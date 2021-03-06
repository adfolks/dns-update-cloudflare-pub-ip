package main

import (
		"bytes"
		"encoding/json"
		"fmt"
		"io/ioutil"
		"net/http"
		"github.com/joho/godotenv"
		"os"
		"time"
		"strings"
)

//Payload defines the json request object
type Payload struct {
		Type string `json:"type"`
		Name string `json:"name"`
		Content string `json:"content"`
		TTL int `json:"ttl"`
		Proxied bool `json:"proxied"`
}

//IpAPI defines the json response object for Origin IP api
type IpAPI struct {
		Origin string `json:"origin"`
}

type AutoGenerated struct {
		Result []Result `json:"result"`
		Success bool `json:"success"`
		Errors []interface{} `json:"errors"`
		Messages []interface{} `json:"messages"`
		ResultInfo ResultInfo `json:"result_info"`
}
type Meta struct {
		AutoAdded bool `json:"auto_added"`
		ManagedByApps bool `json:"managed_by_apps"`
		ManagedByArgoTunnel bool `json:"managed_by_argo_tunnel"`
		Source string `json:"source"`
}
type Result struct {
		ID string `json:"id"`
		ZoneID string `json:"zone_id"`
		ZoneName string `json:"zone_name"`
		Name string `json:"name"`
		Type string `json:"type"`
		Content string `json:"content"`
		Proxiable bool `json:"proxiable"`
		Proxied bool `json:"proxied"`
		TTL int `json:"ttl"`
		Locked bool `json:"locked"`
		Meta Meta `json:"meta"`
		CreatedOn time.Time `json:"created_on"`
		ModifiedOn time.Time `json:"modified_on"`
}
type ResultInfo struct {
		Page int `json:"page"`
		PerPage int `json:"per_page"`
		Count int `json:"count"`
		TotalCount int `json:"total_count"`
		TotalPages int `json:"total_pages"` 
}

var Dns_names []string

var Record_ids = make([]string,0) 
var Content_ip string 


func main() {

		fmt.Println("Starting CloudFlare ip update service")

		var err error
		err = godotenv.Load() // Loads environment values from .env file
			if err != nil {
					fmt.Printf("Error getting env, not comming through %v\n", err)
			} else {
					fmt.Print("Getting the env values\n")
			}

		Dns_names = getDnsnames() // Gets the list all the dns names from the env file .

		Content_ip = getMyIP() // Gets ip of the current system. 

		Record_ids= getRecordId(Dns_names) // Gets record_id of all the dns . Array , holding dns names is passed . 

		updateMultipleDns() // Multiple dns update function is called. 

}

// The function updates all the dns .
// Calls update cloudflare function passing payload and record_id.
// Repeats the function call for each of the dns_name. 

func getDnsnames() []string{

		str := os.Getenv("DNS_NAMES")
		res1 := strings.Split(str, ",") 
		return res1
}

// The function calls updateCloudflareRecord() for each dns_names 
// Passing payload and record_id for each dns names.
func updateMultipleDns(){

		for a := range Dns_names{
				payLoad := Payload{
									Type: "A",
									Name: Dns_names[a], // Dns names from the array.
									Content: Content_ip,
									TTL: 120,
									Proxied: false,
				}
				record_id := Record_ids[a] // Record id from the array. 

		updateCloudflareRecord(payLoad, record_id) // payload and record_id are passed. 
	}
}




// Funtion returns current ip address of the server.( This system)
func getMyIP () string {

		response, _ := http.Get("https://httpbin.org/ip")
		var responseIPJson IpAPI
		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
				print("Error parsing response from service")
			} else {
				json.Unmarshal(body, &responseIPJson)
			}
		content_ip := responseIPJson.Origin 
		return content_ip
}


// The function gets all the record_id and unmarshlled to 'value' variable.
// Value is parsed each time comparing each of the dns_names. The corresponding record_id is fetched and appended to array.
// An array containing all the record_ids is returned.
func getRecordId( DNS_nameslist []string ) []string {

		ListDns_names:= DNS_nameslist // Holds the dns_names array .

		req,err := http.NewRequest("GET","https://api.cloudflare.com/client/v4/zones/"+os.Getenv("ZONE_ID")+"/dns_records",nil)

		req.Header.Add("Authorization", "Bearer "+os.Getenv("API_KEY"))
		res, err := http.DefaultClient.Do(req)

		if err != nil {
			fmt.Print("something wrong while sending req " + err.Error())
			} else {
				res, _ := ioutil.ReadAll(res.Body)

				var value AutoGenerated // ' value' will hold the unmarshalled response json.

			if err == nil && res != nil {
						err = json.Unmarshal([]byte(res), &value)
					}

		for b := range ListDns_names{
			for l := range value.Result {
					if (ListDns_names[b] == value.Result[l].Name ){
					//fmt.Printf( value.Result[l].Name)
					//fmt.Println()
					//Record_ID = value.Result[l].ID
					//fmt.Println()
					r := value.Result[l].ID
					Record_ids = append(Record_ids , r)
					}
			}
		}

	}
	return Record_ids
}

// Function updates the cloudflare
func updateCloudflareRecord(payload Payload, record_id string) {

		url := "https://api.cloudflare.com/client/v4/zones/"+os.Getenv("ZONE_ID")+"/dns_records/"+record_id
		jsonPayload, _ := json.Marshal(payload)
		apiKey := os.Getenv("API_KEY")
		req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonPayload))
		req.Header.Add("Authorization", "Bearer "+apiKey)
		req.Header.Add("Content-Type", "application/json")
		res, err := http.DefaultClient.Do(req)
		if err != nil {
		fmt.Print("something wrong while sending req " + err.Error())
		} else {
		res, _ := ioutil.ReadAll(res.Body)
		print(string(res))
		fmt.Println("")
		}

}
