package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	// "io"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zde37/pinata-go-sdk/pinata"
)

// ---- ENV ----
var (
	RPC_URL       = "https://testnet.sapphire.oasis.io"
	CONTRACT_ADDR = common.HexToAddress("0xd02E5Fe32468C5e3857E8958ECcCb6616b0F16Fb")
	auth          *pinata.Auth
	client        *pinata.Client
)

func main() {
	cli, _ := ethclient.Dial("wss://testnet.sapphire.oasis.io/ws")

	sig := []byte("OrderCreated(uint256,uint256,address,uint256)")
	topic := crypto.Keccak256Hash(sig)

	log.Printf("Topic: %s", topic.Hex())

	q := ethereum.FilterQuery{
		Addresses: []common.Address{CONTRACT_ADDR},
		Topics:    [][]common.Hash{{topic}},
	}

	logs := make(chan types.Log)
	sub, err := cli.SubscribeFilterLogs(context.Background(), q, logs)
	if err != nil {
		log.Fatal(err)
	}

	auth = pinata.NewAuthWithJWT(os.Getenv("JWT_TOKEN"))
	client = pinata.New(auth)

	log.Println("Listening for events...")

	for {
		select {
		case err := <-sub.Err():
			log.Println("Error:", err)
		case vLog := <-logs:
			go handle(vLog, topic)
		}
	}
}

func handle(vLog types.Log, topic common.Hash) {
	log.Printf("Received log: %v", vLog)

	var ev struct {
		DatasetId  *big.Int
		OrderId    *big.Int
		Researcher common.Address
		Amount     *big.Int
	}

	abiObj, _ := abi.JSON(strings.NewReader(`[{
	  "anonymous":false,
	  "inputs":[
		 {"indexed":true,"name":"datasetId","type":"uint256"},
		 {"indexed":true,"name":"orderId","type":"uint256"},
		 {"indexed":true,"name":"researcher","type":"address"},
		 {"indexed":false,"name":"amount","type":"uint256"}],
	  "name":"OrderCreated",
	  "type":"event"}]`))

	if vLog.Topics[0] == topic {
		// indexed fields come from topics[1..3]
		ev.DatasetId = new(big.Int).SetBytes(vLog.Topics[1].Bytes())
		ev.OrderId = new(big.Int).SetBytes(vLog.Topics[2].Bytes())
		ev.Researcher = common.BytesToAddress(vLog.Topics[3].Bytes())

		// non‑indexed amount sits in Data
		if err := abiObj.UnpackIntoInterface(&ev, "OrderCreated", vLog.Data); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("order %d on dataset %d by %s amount %s\n",
			ev.OrderId, ev.DatasetId, ev.Researcher.Hex(), ev.Amount)

		orderId := ev.OrderId.Uint64()
		datasetId := ev.DatasetId.Uint64()

		computeHandler(orderId, datasetId)
	}
}

func computeHandler(orderId uint64, datasetId uint64) {
	log.Printf("Order ID: %d", orderId)

	order, err := getStake(orderId, datasetId)
	if err != nil {
		log.Printf("Error getting order: %v", err)
		return
	}
	log.Printf("Order: %v", order)
	log.Printf("Order Dataset ID: %d", order.DatasetId)

	datares, err := getDataHash(order.DatasetId)
	if err != nil {
		log.Printf("Error getting data hash: %v", err)
		return
	}
	log.Printf("Data: %v", datares)

	// text, err := fetchIPFS(data.IPFSHash)
	text, err := fetchIPFS(datares.IPFSHash) // Test IPFS hash
	if err != nil {
		log.Printf("Error fetching IPFS content: %v", err)
		return
	}

	log.Printf("IPFS content: %s", text)
	fixedText := strings.ReplaceAll(text, "\"[", "[")
	fixedText = strings.ReplaceAll(fixedText, "]\"", "]")
	fixedText = strings.ReplaceAll(fixedText, "{timestamp:", "{\"timestamp\":")
	fixedText = strings.ReplaceAll(fixedText, ", heartRate:", ", \"heartRate\":")
	fixedText = strings.ReplaceAll(fixedText, ", bloodOxygenLevel:", ", \"bloodOxygenLevel\":")

	log.Printf("Fixed text: %s", fixedText)
	var dataEntries []DataEntry
	if err := json.Unmarshal([]byte(fixedText), &dataEntries); err != nil {
		log.Printf("Error unmarshalling JSON: %v", err)
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
		log.Printf("Error marshalling average data: %v", err)
		return
	}
	averageDataCID, err := addIPFS(string(averageDataJson))
	if err != nil {
		log.Printf("Error adding average data to IPFS: %v", err)
		return
	}

	log.Printf("Average data CID: %s", averageDataCID)

	err = completeOrder(order.OrderId, order.DatasetId)
	if err != nil {
		log.Printf("Error completing order: %v", err)
		return
	}
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
