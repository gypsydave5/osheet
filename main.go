package main

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	// Notes
	// 1. Create a service account on the google app whatsface thing mhttps://console.developers.google.com/apis/credentials
	// 2. download the json for the service account and call it service_account.json or whatever
	b, err := ioutil.ReadFile("service_account.json")

	// 3. give the email address of the service account access to the sheet
	// 4. find the id of the sheet
	spreadsheetId := "1SGb-MWvYvw9Jiw3_Y5rhYtWIEOdAk3fqnIvGzPmQz90"

	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// this creates a JWT config thing from the service account JSON
	config, err := google.JWTConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	// config can make an HTTP client with some context lol
	client := config.Client(context.TODO())

	// make a sheets client
	srv, err := sheets.NewService(context.TODO(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	s := Sheet{
		sheetId: spreadsheetId,
		svc:     srv,
	}

	log.Fatal(http.ListenAndServe(":8080", s))
}

type Sheet struct  {
	sheetId string
	svc *sheets.Service
}

func (s Sheet) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		s.sheetsHandler(w, r)
		return
	default:
		s.sheetHandler(w, r)
		return
	}
}

type Link struct {
	Path string `json:"path"`
	Rel string `json:"rel"`
}

func (s *Sheet) sheetsHandler(w http.ResponseWriter, r *http.Request) {
	// profit I guess
	sheet, err := s.svc.Spreadsheets.Get(s.sheetId).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	type sheetResponse struct {
		Names []string `json:"names"`
		Links []*Link `json:"links"`
	}

	resp := &sheetResponse{}

	for i := range sheet.Sheets {
		resp.Names = append(resp.Names, sheet.Sheets[i].Properties.Title)
		resp.Links = append(resp.Links, sheetLink(sheet.Sheets[i].Properties.Title))
	}

	json.NewEncoder(w).Encode(resp)
}

func (s Sheet) sheetHandler(w http.ResponseWriter, r *http.Request) {
	split := strings.Split(r.URL.Path, "/")
	sheetName, _ := url.PathUnescape(split[1])
	if len(split) == 3 {
		 s.rowHandler(w, r)
		 return
	}

	sheet, _ := s.svc.Spreadsheets.Values.Get(s.sheetId, sheetName).Do()
	keys := sheet.Values[0]
	var resp struct {
		Rows []map[string]interface{} `json:"rows"`
	}

	for _, row := range sheet.Values[1:] {
		r := make(map[string]interface{})
		for j := range row {
			r[keys[j].(string)] = row[j]
		}
		resp.Rows = append(resp.Rows, r)
	}

	fmt.Println(resp)
	json.NewEncoder(w).Encode(resp)
}

func (s Sheet) rowHandler(w http.ResponseWriter, r *http.Request) {
	split := strings.Split(r.URL.Path, "/")
	sheetName, _ := url.PathUnescape(split[1])
	rowNum, _ := url.PathUnescape(split[2])

	sheet, _ := s.svc.Spreadsheets.Values.Get(s.sheetId, sheetName+"!1:1").Do()
	keys := sheet.Values[0]

	sheet, _ = s.svc.Spreadsheets.Values.Get(s.sheetId, sheetName+"!"+rowNum+":"+rowNum).Do()
	resp := make(map[string]interface{})
	for i := range sheet.Values[0] {
		resp[keys[i].(string)]  = sheet.Values[0][i]
	}

	json.NewEncoder(w).Encode(resp)
}

func sheetLink(title string) *Link {
	return &Link{
		Rel: title,
		Path: sheetPath(title),
	}
}

func sheetPath(title string) string {
	return "/" + url.PathEscape(title)
}
