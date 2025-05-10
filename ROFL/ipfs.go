package main

import (
	"io"
	"log"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
)

func fetchIPFS(ipfsHash string) (string, error) {
	sh := shell.NewShell("127.0.0.1:5001") // IPFS shell :contentReference[oaicite:8]{index=8}

	// Get content from IPFS
	rc, err := sh.Cat(ipfsHash)
	if err != nil {
		log.Printf("Failed to retrieve IPFS content: %v", err)
		return "", err
	}

	bytes, err := io.ReadAll(rc)
	if err != nil {
		log.Printf("Failed to read IPFS content: %v", err)
		return "", err
	}

	text := string(bytes)
	log.Printf("IPFS content: %s", text)

	return text, nil
}

func addIPFS(content string) (string, error) {
	sh := shell.NewShell("127.0.0.1:5001") // IPFS shell :contentReference[oaicite:8]{index=8}
	// Add content to IPFS
	cid, err := sh.Add(strings.NewReader(content))
	if err != nil {
		log.Printf("Failed to add content to IPFS: %v", err)
		return "", err
	}
	log.Printf("Added content to IPFS with CID: %s", cid)
	return cid, nil
}
