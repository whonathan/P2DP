package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/valyala/fasthttp"
)

const dbPath = "prepaidData.db"

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS prepaidData (
			idpel TEXT PRIMARY KEY,
			nama TEXT,
			tarif TEXT,
			daya TEXT,
			kdrbm TEXT,
			blth TEXT,
			merk_meter TEXT,
			nomor_meter TEXT
		)
	`)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func (c *PLNClient) FetchAndStorePrepaidData(username, unitup string) error {
	req := c.reqPool.Get().(*fasthttp.Request)
	resp := c.respPool.Get().(*fasthttp.Response)
	defer c.reqPool.Put(req)
	defer c.respPool.Put(resp)

	payload := fmt.Sprintf("username=%s&tgllogin=%s&unitup=%s",
		username,
		time.Now().Format("01/02/2006"),
		unitup,
	)

	req.SetRequestURI(baseURL + "/getAllStanPrabayar")
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Apache-HttpClient/UNAVAILABLE (java 1.4)")
	req.Header.Set("Connection", "Keep-Alive")
	req.SetBodyString(payload)

	if err := c.client.Do(req, resp); err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	var response Response
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		return fmt.Errorf("json unmarshal failed: %v", err)
	}

	if response.Success != 1 {
		return fmt.Errorf("API returned unsuccessful response")
	}

	if len(response.Stan) == 0 {
		return fmt.Errorf("no data returned from API")
	}

	return c.storeMeterData(response.Stan)
}

func (c *PLNClient) storeMeterData(data []DBMeterData) error {
	tx, err := c.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM prepaidData")
	if err != nil {
		return fmt.Errorf("failed to clear existing data: %v", err)
	}

	stmt, err := tx.Prepare(`
		INSERT INTO prepaidData (
			idpel, nama, tarif, daya, kdrbm, blth, merk_meter, nomor_meter
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	for _, meter := range data {
		_, err = stmt.Exec(
			meter.IDPEL,
			meter.Nama,
			meter.Tarif,
			meter.Daya,
			meter.KDRBM,
			meter.BLTH,
			meter.MerkMeter,
			meter.NomorMeter,
		)
		if err != nil {
			return fmt.Errorf("insert failed for IDPEL %s: %v", meter.IDPEL, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

func (c *PLNClient) GetKDRBMData() ([]KDRBMData, error) {
	query := `SELECT kdrbm, COUNT(*) as count 
			 FROM prepaidData 
			 GROUP BY kdrbm 
			 ORDER BY kdrbm`

	rows, err := c.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []KDRBMData
	for rows.Next() {
		var data KDRBMData
		if err := rows.Scan(&data.KDRBM, &data.Count); err != nil {
			return nil, err
		}
		result = append(result, data)
	}
	return result, nil
}

func (c *PLNClient) GetReadingsByKDRBM(kdrbm string) ([]DBMeterData, error) {
	query := `SELECT idpel, nama, tarif, daya, kdrbm, blth, merk_meter, nomor_meter 
			 FROM prepaidData 
			 WHERE kdrbm = ?`

	rows, err := c.db.Query(query, kdrbm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var readings []DBMeterData
	for rows.Next() {
		var data DBMeterData
		if err := rows.Scan(&data.IDPEL, &data.Nama, &data.Tarif, &data.Daya,
			&data.KDRBM, &data.BLTH, &data.MerkMeter, &data.NomorMeter); err != nil {
			return nil, err
		}
		readings = append(readings, data)
	}
	return readings, nil
}
