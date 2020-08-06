package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	b, err := ioutil.ReadFile("service_account.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	sheetID, _ := os.LookupEnv("SHEET_ID")
	accountJSON, _ := os.LookupEnv("SERVICE_ACCOUNT_JSON")
	port, _ := os.LookupEnv("PORT")

	fmt.Println(accountJSON)

	b = []byte(accountJSON)

	s, err := newSheetsService(b, sheetID)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), newHTTPSheetsService(s)))
}
