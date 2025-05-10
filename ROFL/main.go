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
	CONTRACT_ADDR = common.HexToAddress("0xe041b50CA3ED1c23F8D7139a11Ed107a010937D5")
	auth          *pinata.Auth
	client        *pinata.Client
	privKey       *string
	pubKey        *string
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

	// pk, pu, err := GenerateKeyPair()
	// if err != nil {
	// 	log.Printf("Error generating key pair: %v", err)
	// 	return
	// }

	pk := "0x69b32ac0113ca525d3c354771d2ed687f1edc2b014400afc2e04995ffe27a964"
	pu := "0x04fb9e93abb1862b1d8c06340d340e84058f9ce545e9d84d1a4d29258286a08c800c460a7bca92c155f74029fcc9a8e3ba8ae46ecd37311fa349ff2c75b45f001c"

	privKey = &pk
	pubKey = &pu

	// err = storePubKeyInSC(pu)
	// if err != nil {
	// 	log.Printf("Error storing public key in SC: %v", err)
	// 	return
	// }

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
	encryptedText, err := fetchIPFS(datares.IPFSHash) // Test IPFS hash
	if err != nil {
		log.Printf("Error fetching IPFS content: %v", err)
		return
	}

	text, err := DecryptData([]byte(encryptedText))
	if err != nil {
		log.Printf("Error decrypting IPFS content: %v", err)
		return
	}

	log.Printf("IPFS content: %s", text)
	// Remove outer quotes if needed
	fixedText := strings.ReplaceAll(text, "\"[", "[")
	fixedText = strings.ReplaceAll(fixedText, "]\"", "]")

	// Fix missing quotes around timestamp and ensure it's consistent
	fixedText = strings.ReplaceAll(fixedText, "{timestamp:", "{\"timestamp\":")
	fixedText = strings.ReplaceAll(fixedText, "{ timestamp:", "{\"timestamp\":")

	// Fix other fields
	fixedText = strings.ReplaceAll(fixedText, ", heartRate:", ", \"heartRate\":")
	fixedText = strings.ReplaceAll(fixedText, ", bloodOxygenLevel:", ", \"bloodOxygenLevel\":")

	// Fix missing commas between objects
	fixedText = strings.ReplaceAll(fixedText, "},{", "},{")

	// Use regex to fix any remaining issues - better approach would be to add this
	fixedText = strings.ReplaceAll(fixedText, "}{", " },{")

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
