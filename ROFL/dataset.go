package main

import (
	"context"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
)

func getDataHash(id uint64) (DataResponse, error) {
	cli, err := ethclient.Dial(RPC_URL) // JSON‑RPC call :contentReference[oaicite:3]{index=3}
	if err != nil {
		return DataResponse{}, err
	}

	escAbi, _ := abi.JSON(strings.NewReader(ABI_JSON))
	input, _ := escAbi.Pack("getDatasetHash", big.NewInt(int64(id)))

	// Note: Sapphire requires ECIES envelope, but JSON‑RPC GET is fine for eth_call.
	msg := ethereum.CallMsg{To: &CONTRACT_ADDR, Data: input}
	out, err := cli.CallContract(context.Background(), msg, nil) // eth_call :contentReference[oaicite:4]{index=4}
	if err != nil {
		return DataResponse{}, err
	}

	var data DataResponse
	err = escAbi.UnpackIntoInterface(&data, "getDatasetHash", out)
	if err != nil {
		return DataResponse{}, err
	}
	return data, nil
}
