package main

import (
	"io"
	"log"
	"os"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
)

func getIPFSHost() string {
	host := os.Getenv("IPFS_HOST")
	if host == "" {
		return "127.0.0.1:5001" // Default fallback
	}
	return host
}

func fetchIPFS(ipfsHash string) (string, error) {
	ipfsHost := getIPFSHost()
	sh := shell.NewShell(ipfsHost)
	log.Printf("Connecting to IPFS at: %s", ipfsHost)

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
	ipfsHost := getIPFSHost()
	sh := shell.NewShell(ipfsHost)
	log.Printf("Connecting to IPFS at: %s", ipfsHost)

	// Add content to IPFS
	cid, err := sh.Add(strings.NewReader(content))
	if err != nil {
		log.Printf("Failed to add content to IPFS: %v", err)
		return "", err
	}
	log.Printf("Added content to IPFS with CID: %s", cid)
	return cid, nil
}
