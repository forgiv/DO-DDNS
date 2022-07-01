package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type RecordsResponse struct {
	DomainRecords []struct {
		ID       int         `json:"id"`
		Type     string      `json:"type"`
		Name     string      `json:"name"`
		Data     string      `json:"data"`
		Priority interface{} `json:"priority"`
		Port     interface{} `json:"port"`
		TTL      int         `json:"ttl"`
		Weight   interface{} `json:"weight"`
		Flags    interface{} `json:"flags"`
		Tag      interface{} `json:"tag"`
	} `json:"domain_records"`
	Links struct {
	} `json:"links"`
	Meta struct {
		Total int `json:"total"`
	} `json:"meta"`
}

func main() {
	// args = [domain, subdomain, apiKey]
	args := os.Args[1:]

	// Create HTTP client
	client := &http.Client{}

	// First step, find domain A record that matches our domain and subdomain
	recordId := getRecordId(client, args[0], args[1], args[2])

	fmt.Printf("Record ID: %d\n", recordId)
}

func getRecordId(client *http.Client, domain string, subdomain string, apiKey string) int {
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://api.digitalocean.com/v2/domains/%s/records?type=A", domain), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var records RecordsResponse
	if err = json.Unmarshal(body, &records); err != nil {
		fmt.Println("Unable to unmarshal JSON")
	}

	var recordId int
	for _, record := range records.DomainRecords {
		if record.Name == subdomain {
			recordId = record.ID
		}
	}

	return recordId
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
