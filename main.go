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
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://api.digitalocean.com/v2/domains/%s/records?type=A", args[0]), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", args[2]))
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
		if record.Name == args[1] {
			recordId = record.ID
		}
	}

	fmt.Println(recordId)
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}
