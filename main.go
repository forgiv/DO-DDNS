package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type DomainRecord struct {
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
}

type RecordResponse struct {
	DomainRecord DomainRecord `json:"domain_record"`
}

type RecordsResponse struct {
	DomainRecords []DomainRecord `json:"domain_records"`
	Links         struct {
	} `json:"links"`
	Meta struct {
		Total int `json:"total"`
	} `json:"meta"`
}

type NewRecordRequest struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Data string `json:"data"`
	Ttl  int    `json:"ttl"`
}

type UpdateRecordRequest struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func main() {
	// args = [domain, subdomain, apiKey]
	args := os.Args[1:]

	// Create HTTP client
	client := &http.Client{}

	// First step, find domain A record that matches our domain and subdomain
	recordId := getRecordId(client, args[0], args[1], args[2])

	// Second step, decide whether to create new record or update existing and do it
	var record RecordResponse
	if recordId == 0 {
		record = createNewRecord(client, args[0], args[1], args[2])
	} else {
		record = updateExistingRecord(client, args[0], recordId, args[2])
	}

	fmt.Println(prettyPrint(record))
}

// Return pretty formatted string
func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")

	return string(s)
}

// Get preferred outbound ip of this machine
func getOutboundIP(client *http.Client) string {
	req, _ := http.NewRequest("GET", "https://api.ipify.org", nil)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return string(body)
}

// Return ID of existing record to update or zero if no record
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

	recordId := 0
	for _, record := range records.DomainRecords {
		if record.Name == subdomain {
			recordId = record.ID
		}
	}

	return recordId
}

// Creates a new A record pointing to IP of machine
func createNewRecord(client *http.Client, domain string, subdomain string, apiKey string) RecordResponse {
	data := NewRecordRequest{
		Type: "A",
		Name: subdomain,
		Data: getOutboundIP(client),
		Ttl:  30,
	}

	jsonData, _ := json.Marshal(data)

	req, _ := http.NewRequest("POST", fmt.Sprintf("https://api.digitalocean.com/v2/domains/%s/records", domain), strings.NewReader(string(jsonData)))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var record RecordResponse
	if err = json.Unmarshal(body, &record); err != nil {
		fmt.Println("Unable to unmarshal JSON")
	}

	return record
}

func updateExistingRecord(client *http.Client, domain string, recordID int, apiKey string) RecordResponse {
	data := UpdateRecordRequest{
		Type: "A",
		Data: getOutboundIP(client),
	}

	jsonData, _ := json.Marshal(data)

	req, _ := http.NewRequest("PATCH", fmt.Sprintf("https://api.digitalocean.com/v2/domains/%s/records/%d", domain, recordID), strings.NewReader(string(jsonData)))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var record RecordResponse
	if err = json.Unmarshal(body, &record); err != nil {
		fmt.Println("Unable to unmarshal JSON")
	}

	return record
}
