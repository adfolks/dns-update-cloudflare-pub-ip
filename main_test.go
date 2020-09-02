package main


import (
	"testing"
	"github.com/joho/godotenv"
	"os"
	"fmt"
)

var getRecordMock func(DNS_name string) string

func TestMain(t *testing.T){
	var err error
	err = godotenv.Load()  // Loads environment values from .env file
	if err != nil {
				fmt.Printf("Error getting env, not comming through %v\n", err)
			} else {
					fmt.Print("Getting the env values\n")
			}

}


func TestGetMyIP(t *testing.T){


	my_ip := os.Getenv("Test_MyIP")
    response := getMyIP()
    if response != my_ip  {
	   t.Errorf("Somthing  was incorrect, got: %s, want: %s.", response, my_ip )
	}
}

func TestGetRecordId(t *testing.T) {

	
}