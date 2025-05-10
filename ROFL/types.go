package main

type computeReq struct {
	OrderId   uint64 `json:"orderId"`
	RAddress  string `json:"raAddress"`
	DatasetId uint64 `json:"datasetId"`
	Amount    int64  `json:"amount"`
}

type DataEntry struct {
	Timestamp        int64   `json:"timestamp"`
	HeartRate        int64   `json:"heartRate"`
	BloodOxygenLevel float64 `json:"bloodOxygenLevel"`
}

type Data struct {
	Owner   string      `json:"owner"`
	Entries []DataEntry `json:"entries"`
}

type DataResponse struct {
	IPFSHash string `json:"ipfsHash"`
}

type Order struct {
	OrderId      uint64 `json:"orderId"`
	DatasetId    uint64 `json:"datasetId"`
	Researcher   string `json:"researcher"` // the payer
	Patient      string `json:"patient"`    // dataset owner / payee
	Amount       uint64 `json:"amount"`
	TokenAddress string `json:"tokenAddress"`
	Timestamp    uint64 `json:"timestamp"` // creation time
	Completed    bool   `json:"completed"`
}
