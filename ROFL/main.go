package main

import (
	"context"
	"encoding/json"
	"log"
	"math/big"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ---- ENV ----
var (
	RPC_URL       = "https://testnet.sapphire.oasis.io"
	CONTRACT_ADDR = common.HexToAddress("0xBb18E81753179d29071772DcEf8f8B2dcd368184")
	IPFS_GATEWAY  = "nftstorage.link"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]string{"message": "Hello World"}
		enc := json.NewEncoder(w)
		enc.Encode(resp)
	})

	http.HandleFunc("/test", testfunc)

	http.HandleFunc("/compute", computeHandler)

	log.Println("API listening on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func testfunc(w http.ResponseWriter, r *http.Request) {
	text, err := readContract()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	log.Printf("IPFS content: %s", text)

	// Send the content as a response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(text))
	log.Printf("IPFS content sent as response")
}

func computeHandler(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var req computeReq
	if err := dec.Decode(&req); err != nil {
		http.Error(w, "bad json", 400)
		return
	}

	log.Printf("Received request: %v", req)
	log.Printf("Order ID: %d", req.OrderId)

	order, err := getStake(req.OrderId, req.DatasetId)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	log.Printf("Order: %v", order)
	log.Printf("Order Dataset ID: %d", order.DatasetId)

	datares, err := getDataHash(order.DatasetId)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	log.Printf("Data: %v", datares)

	// text, err := fetchIPFS(data.IPFSHash)
	text, err := fetchIPFS(datares.IPFSHash) // Test IPFS hash
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	log.Printf("IPFS content: %s", text)

	fixedText := strings.ReplaceAll(text, "{timestamp:", "{\"timestamp\":")
	fixedText = strings.ReplaceAll(fixedText, ", heartRate:", ", \"heartRate\":")
	fixedText = strings.ReplaceAll(fixedText, ", bloodOxygenLevel:", ", \"bloodOxygenLevel\":")

	var dataEntries []DataEntry
	if err := json.Unmarshal([]byte(fixedText), &dataEntries); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// --- 3. process data ----------------------------------------------
	var averageHR int64
	var averageBOL float64
	for i := range dataEntries {
		entry := &dataEntries[i]
		averageHR += entry.HeartRate
		averageBOL += entry.BloodOxygenLevel
	}
	averageHR /= int64(len(dataEntries))
	averageBOL /= float64(len(dataEntries))

	//add to ipfs an object with the average heart rate and blood oxygen level
	averageData := map[string]interface{}{
		"averageHeartRate":        averageHR,
		"averageBloodOxygenLevel": averageBOL,
	}
	averageDataJson, err := json.Marshal(averageData)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	averageDataCID, err := addIPFS(string(averageDataJson))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	log.Printf("Average data CID: %s", averageDataCID)

	err = completeOrder(order.OrderId, order.DatasetId)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Send the average data CID as a response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(averageDataCID))
	log.Printf("Average data CID sent as response")
}

func readContract() (string, error) {
	// read a contract method getDataset and input 0
	log.Printf("Reading contract method getDataset")
	cli, err := ethclient.Dial(RPC_URL)
	if err != nil {
		return "", err
	}
	escAbi, _ := abi.JSON(strings.NewReader(ABI_JSON))
	input, _ := escAbi.Pack("getDataset", big.NewInt(0))
	msg := ethereum.CallMsg{To: &CONTRACT_ADDR, Data: input}
	out, err := cli.CallContract(context.Background(), msg, nil) // eth_call :contentReference[oaicite:4]{index=4}
	if err != nil {
		return "", err
	}
	log.Printf("out: %x", out)

	var data any
	err = escAbi.UnpackIntoInterface(&data, "getDataset", out)
	if err != nil {
		return "", err
	}

	// Convert the data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	// Convert JSON to string
	jsonString := string(jsonData)
	log.Printf("JSON string: %s", jsonString)
	return jsonString, nil

}
