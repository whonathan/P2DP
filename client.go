package main

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

const (
	baseURL              = BaseURL              // Use from config.go
	meterDataEndpoint    = MeterDataEndpoint    // Use from config.go
	maxRetries           = MaxRetries           // Use from config.go
	requestTimeout       = RequestTimeout       // Use from config.go
	maxConcurrentUploads = MaxConcurrentUploads // Use from config.go
)

var photoEndpoints = []string{
	"/uploadFoto52",
	"/uploadFoto2",
	"/uploadFoto3",
}

type PLNClient struct {
	client   *fasthttp.Client
	maxConns int
	reqPool  sync.Pool
	respPool sync.Pool
	coordGen *CoordinateGenerator
	timeGen  *TimeGenerator
	db       *sql.DB
}

type PLNClientInterface interface {
	SubmitReading(ctx context.Context, meterData MeterData, photoData PhotoData) SubmissionResult
	GetReadingsByKDRBM(kdrbm string) ([]DBMeterData, error)
	Close() error
}

func NewPLNClient(maxConns int, username string) (*PLNClient, error) {
	db, err := initDB()
	if err != nil {
		return nil, fmt.Errorf("database initialization failed: %v", err)
	}

	seed := time.Now().UnixNano()
	return &PLNClient{
		client: &fasthttp.Client{
			MaxConnsPerHost:          maxConns,
			ReadTimeout:              requestTimeout,
			WriteTimeout:             requestTimeout,
			MaxIdleConnDuration:      5 * time.Minute,
			MaxConnDuration:          10 * time.Minute,
			NoDefaultUserAgentHeader: true,
		},
		maxConns: maxConns,
		reqPool: sync.Pool{
			New: func() interface{} {
				return fasthttp.AcquireRequest()
			},
		},
		respPool: sync.Pool{
			New: func() interface{} {
				return fasthttp.AcquireResponse()
			},
		},
		coordGen: NewCoordinateGenerator(seed, username),
		timeGen:  NewTimeGenerator(seed),
		db:       db,
	}, nil
}

func (c *PLNClient) Close() error {
	if c.client != nil {
		c.client.CloseIdleConnections()
	}
	c.reqPool.New = nil
	c.respPool.New = nil
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

func (c *PLNClient) sendRequest(ctx context.Context, method, url string, payload string, headers map[string]string) ([]byte, error) {
	reqCtx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	select {
	case <-reqCtx.Done():
		return nil, fmt.Errorf("request timeout after %v: %w", requestTimeout, reqCtx.Err())
	default:
		req := c.reqPool.Get().(*fasthttp.Request)
		resp := c.respPool.Get().(*fasthttp.Response)
		defer c.reqPool.Put(req)
		defer c.respPool.Put(resp)

		req.Reset()
		resp.Reset()

		if c.client == nil {
			return nil, fmt.Errorf("client not initialized")
		}

		req.SetRequestURI(url)
		req.Header.SetMethod(method)
		req.SetBodyString(payload)

		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Accept-Encoding", "gzip")

		for key, value := range headers {
			req.Header.Set(key, value)
		}

		var err error
		for i := 0; i < maxRetries; i++ {
			if err = c.client.Do(req, resp); err == nil {
				if resp.StatusCode() != fasthttp.StatusOK {
					return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
				}
				responseBody := make([]byte, len(resp.Body()))
				copy(responseBody, resp.Body())
				return responseBody, nil
			}
			if i < maxRetries-1 {
				backoff := time.Duration(1<<uint(i)) * 100 * time.Millisecond
				time.Sleep(backoff)
			}
		}
		return nil, fmt.Errorf("request failed after %d retries: %v", maxRetries, err)
	}
}

func (c *PLNClient) SubmitReading(ctx context.Context, meterData MeterData, photoData PhotoData) SubmissionResult {
	var result SubmissionResult
	var resultMutex sync.Mutex
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			resultMutex.Lock()
			result.MeterDataError = ctx.Err()
			resultMutex.Unlock()
			return
		default:
			payload := buildMeterDataPayload(meterData)
			headers := getDefaultHeaders("Dalvik/2.1.0")

			response, err := c.sendRequest(ctx, "POST", baseURL+meterDataEndpoint, payload, headers)
			resultMutex.Lock()
			result.MeterDataResponse = response
			if err != nil {
				result.MeterDataError = fmt.Errorf("meter data submission failed: %w", err)
			}
			resultMutex.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		photoResults := make([]PhotoUploadResult, len(photoEndpoints))
		var photoResultsMutex sync.Mutex
		var photoWg sync.WaitGroup
		photoWg.Add(len(photoEndpoints))

		for i, endpoint := range photoEndpoints {
			go func(index int, endp string) {
				defer photoWg.Done()

				select {
				case <-ctx.Done():
					photoResults[index] = PhotoUploadResult{
						Endpoint: endp,
						Error:    ctx.Err(),
					}
					return
				default:
					payload := buildPhotoPayload(photoData)
					headers := getDefaultHeaders("Apache-HttpClient/UNAVAILABLE")

					response, err := c.sendRequest(ctx, "POST", baseURL+endp, payload, headers)
					photoResultsMutex.Lock()
					result := PhotoUploadResult{
						Endpoint: endp,
						Response: response,
						Error:    err,
					}
					if err != nil {
						result.Error = fmt.Errorf("photo upload failed: %w", err)
					}
					photoResults[index] = result
					photoResultsMutex.Unlock()
				}
			}(i, endpoint)
		}

		photoWg.Wait()
		result.PhotoResults = photoResults
	}()

	wg.Wait()
	return result
}
