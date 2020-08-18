func update(client *http.Client) {
	fmt.Println("Starting program successfully... Output will be displayed only when records are updated.")
	//Infinite loop to update records over time
	for {
		//GETS current record information
		url := "https://api.cloudflare.com/client/v4/zones/"+os.Getenv("zoneId")+"/dns_records/"+os.Getenv("recordId")
		body := httpRequest(client, "GET", url, nil,  zone)
		jsonData := unjsonify(body)

		numOfRecords = len(jsonData.Result)
		publicIP := getIP()
		for i := 0; i < numOfRecords; i++ {
			recordType := jsonData.Result[i].Type
			recordIP := jsonData.Result[i].Content
			recordIdentifier := jsonData.Result[i].Identifier
			recordName := jsonData.Result[i].Name
			recordProxied := jsonData.Result[i].Proxied

			//Proceeds if is an A Record, AND current IP differs from recorded one
			if recordType == "A" && recordIP != publicIP {
				jsonData := jsonify(recordType, recordName, publicIP, recordProxied) //Creates JSON payload
				//PUTS new record information
				recordURL := url + "/" + recordIdentifier
				httpRequest(client, "PUT", recordURL, jsonData, recordId, zone)

				//Prints after successful update
				setTime[i] = time.Now().Format("2020-18-08 12:0:0 PM")

				tableData[i] = table{
					Name:  recordName,
					IP:    publicIP,
					Proxy: recordProxied,
					Time:  setTime[i],
				}
				fmt.Println("Current Time: " + setTime[i] + "\nUpdated Record: " + recordName + "\nUpdated IP: " + publicIP + "\n")
			} else {
				tableData[i] = table{
					Name:  recordName,
					IP:    publicIP,
					Proxy: recordProxied,
					Time:  setTime[i],
				}
			}
		}
		time.Sleep(interval * time.Second) //Sleeping for n seconds
	}
}

