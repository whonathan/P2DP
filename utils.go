package main

import (
	"fmt"
	"math/rand"
)

func getPowerLimit(daya string) string {
	powerVA := 0
	fmt.Sscanf(daya, "%d", &powerVA)
	powerLimit := float64(powerVA) * 0.0017
	return fmt.Sprintf("%.2f", powerLimit)
}

func getTarifIndex(daya string) string {
	powerVA := 0
	fmt.Sscanf(daya, "%d", &powerVA)
	switch powerVA {
	case 450:
		return "01"
	case 900:
		return "02"
	case 1300:
		return "03"
	case 2200:
		return "04"
	case 3500:
		return "05"
	case 4400:
		return "06"
	case 5500:
		return "07"
	case 7700:
		return "08"
	case 11000:
		return "09"
	default:
		return "01"
	}
}

func randomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func (d *DBMeterData) ToSubmissionData(timeGen *TimeGenerator, coordGen *CoordinateGenerator) *MeterData {
	lat, long := coordGen.GenerateCoordinates()

	return &MeterData{
		IDPEL:            d.IDPEL,
		BLTH:             d.BLTH,
		TglBaca:          timeGen.NextTime(),
		SisaKWH:          fmt.Sprintf("%.2f", randomFloat(4.00, 110.00)),
		KWHKomulatif:     fmt.Sprintf("%d", rand.Intn(2501)+500),
		Latitude:         lat,
		Longitude:        long,
		JumlahTerminal:   "5",
		IndikatorDisplay: "0",
		LCD:              "Normal",
		Keypad:           "Normal",
		CosPhi:           fmt.Sprintf("%.2f", randomFloat(0.50, 1.00)),
		KDBaca:           "NORMAL",
		KDBaca2:          "",
		KDBaca3:          "",
		Tegangan:         fmt.Sprintf("%d", rand.Intn(31)+195),
		TutupMeter:       "0",
		Arus:             fmt.Sprintf("%d", rand.Intn(2)),
		TarifIndex:       getTarifIndex(d.Daya),
		KondisiSegel:     "Ada",
		Relay:            "Tutup",
		IndikatorTemper:  "Tidak Nyala",
		Akurasi:          fmt.Sprintf("%d", rand.Intn(16)+4),
		PowerLimit:       getPowerLimit(d.Daya),
		TransaksiBy:      "",
		StatusTemper:     "0",
	}
}
