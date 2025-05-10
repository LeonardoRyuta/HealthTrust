package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

func fetchIPFS(cid string) (string, error) {
	//make a noraml http get request to https://ipfs.io/ipfs/<cid> and print the response
	resp, err := http.Get("https://ipfs.io/ipfs/" + cid)
	if err != nil {
		log.Printf("Failed to fetch IPFS content: %v", err)
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read IPFS response: %v", err)
		return "", err
	}

	log.Printf("Fetched content from IPFS: %s", string(body))
	return string(body), nil

}

func addIPFS(content string) (string, error) {
	// Upload the content to IPFS using Pinata
	pin, err := client.PinJSON(content, nil)
	if err != nil {
		log.Printf("Failed to upload to IPFS: %v", err)
		return "", err
	}

	log.Printf("Uploaded content to IPFS: %s", pin.IpfsHash)
	return pin.IpfsHash, nil
}
