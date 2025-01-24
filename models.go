package main

type Response struct {
	Stan    []DBMeterData `json:"stan"`
	Success int           `json:"success"`
}

type DBMeterData struct {
	IDPEL      string `json:"idpel"`
	Nama       string `json:"nama"`
	Tarif      string `json:"tarif"`
	Daya       string `json:"daya"`
	KDRBM      string `json:"kdrbm"`
	BLTH       string `json:"blth"`
	MerkMeter  string `json:"merk_meter"`
	NomorMeter string `json:"nomor_meter"`
}

type MeterData struct {
	IDPEL            string
	BLTH             string
	TglBaca          string
	SisaKWH          string
	KWHKomulatif     string
	Latitude         string
	JumlahTerminal   string
	IndikatorDisplay string
	Keypad           string
	CosPhi           string
	LCD              string
	KDBaca           string
	KDBaca2          string
	KDBaca3          string
	NamaFoto         string
	Tegangan         string
	TutupMeter       string
	Longitude        string
	Arus             string
	TarifIndex       string
	KondisiSegel     string
	Relay            string
	IndikatorTemper  string
	Akurasi          string
	PowerLimit       string
	TransaksiBy      string
	StatusTemper     string
}

type PhotoData struct {
	IDPEL        string
	BLTH         string
	UnitUP       string
	NamaFile     string
	PhotoContent string
	TransaksiBy  string
}

type SubmissionResult struct {
	MeterDataResponse []byte
	MeterDataError    error
	PhotoResults      []PhotoUploadResult
}

type PhotoUploadResult struct {
	Endpoint string
	Response []byte
	Error    error
}

type KDRBMData struct {
	KDRBM    string
	Count    int
	Readings []DBMeterData
}
