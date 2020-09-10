package main

import (
		
		"github.com/joho/godotenv"
		"github.com/adfolks/dns-update-cloudflare-pub-ip/src/utils"
		"github.com/adfolks/dns-update-cloudflare-pub-ip/src/models"

		log "github.com/sirupsen/logrus" 
	
)

var Dns_names []string

var Record_ids = make([]string,0) 
var Content_ip string


func main() {

	log.Println("Starting CloudFlare ip update service")

	var err error
	err = godotenv.Load() // Loads environment values from .env file
		if err != nil {
				log.Printf("Error getting env, not comming through %v\n", err)
		} else {
				log.Print("Getting the env values\n")
		}

	Dns_names = utils.GetDnsnames() // Gets the list all the dns names from the env file .

	Content_ip = utils.GetMyIP() // Gets ip of the current system. 

	Record_ids= utils.GetRecordId(Dns_names) // Gets record_id of all the dns . Array , holding dns names is passed . 

	updateMultipleDns() // Multiple dns update function is called. 

}


// The function calls updateCloudflareRecord() for each dns_names 
// Passing payload and record_id for each dns names.
func updateMultipleDns(){

	for a := range Dns_names{
			payLoad := models.Payload{
								Type: "A",
								Name: Dns_names[a], // Dns names from the array.
								Content: Content_ip,
								TTL: 120,
								Proxied: false,
			}
			record_id := Record_ids[a] // Record id from the array. 

			utils.UpdateCloudflareRecord(payLoad, record_id) // payload and record_id are passed. 
}
}