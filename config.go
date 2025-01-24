package main

import "time"

const (
	BaseURL              = "http://portalapp.iconpln.co.id:8000/api-v2-acmt-prod/mobile"
	MeterDataEndpoint    = "/setStanPrabayar"
	MaxRetries           = 3
	RequestTimeout       = 10 * time.Second
	MaxConcurrentUploads = 3
	DBPath               = "prepaidData.db"
)
