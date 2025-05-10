package main

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func waitForReceipt(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	for {
		receipt, err := client.TransactionReceipt(ctx, txHash)
		if err == nil {
			return receipt, nil
		}

		if err != ethereum.NotFound {
			return nil, err
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(2 * time.Second):
			// Wait and retry
		}
	}
}
