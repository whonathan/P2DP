package main

import (
	"net/url"
	"strings"
)

func buildPayload(pairs map[string]string) string {
	var result strings.Builder
	first := true
	for k, v := range pairs {
		if !first {
			result.WriteString("&")
		}
		result.WriteString(url.QueryEscape(k))
		result.WriteString("=")
		result.WriteString(url.QueryEscape(v))
		first = false
	}
	return result.String()
}

func buildMeterDataPayload(data MeterData) string {
	pairs := map[string]string{
		"tglbaca":          data.TglBaca,
		"sisakwh":          data.SisaKWH,
		"kwhkomulatif":     data.KWHKomulatif,
		"latitude":         data.Latitude,
		"jumlah_terminal":  data.JumlahTerminal,
		"indikatordisplay": data.IndikatorDisplay,
		"keypad":           data.Keypad,
		"cosphi":           data.CosPhi,
		"lcd":              data.LCD,
		"kdbaca":           data.KDBaca,
		"kdbaca2":          data.KDBaca2,
		"kdbaca3":          data.KDBaca3,
		"namafoto":         data.NamaFoto,
		"blth":             data.BLTH,
		"tegangan":         data.Tegangan,
		"tutup_meter":      data.TutupMeter,
		"longitude":        data.Longitude,
		"arus":             data.Arus,
		"tarifindex":       data.TarifIndex,
		"kondisi_segel":    data.KondisiSegel,
		"relay":            data.Relay,
		"indikator_temper": data.IndikatorTemper,
		"akurasi":          data.Akurasi,
		"idpel":            data.IDPEL,
		"powerlimit":       data.PowerLimit,
		"transaksiby":      data.TransaksiBy,
		"status_temper":    data.StatusTemper,
	}
	return buildPayload(pairs)
}

func buildPhotoPayload(data PhotoData) string {
	pairs := map[string]string{
		"idpel":       data.IDPEL,
		"blth":        data.BLTH,
		"unitup":      data.UnitUP,
		"namafile":    data.NamaFile,
		"filefoto":    data.PhotoContent,
		"transaksiby": data.TransaksiBy,
	}
	return buildPayload(pairs)
}

func getDefaultHeaders(userAgent string) map[string]string {
	return map[string]string{
		"User-Agent":      userAgent,
		"Content-Type":    "application/x-www-form-urlencoded; charset=UTF-8",
		"Connection":      "Keep-Alive",
		"Accept-Encoding": "gzip",
	}
}
