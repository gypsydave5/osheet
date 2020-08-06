package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type HTTPSheetsService struct {
	*SheetsService
}

func newHTTPSheetsService(s *SheetsService) *HTTPSheetsService {
	return &HTTPSheetsService{s}
}

func (s HTTPSheetsService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.URL.Path {
	case "/":
		s.sheetsHandler(w, r)
		return
	default:
		s.sheetHandler(w, r)
		return
	}
}

func (s *HTTPSheetsService) sheetsHandler(w http.ResponseWriter, r *http.Request) {
	// profit I guess
	sheet, err := s.svc.Spreadsheets.Get(s.sheetId).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	type sheetResponse struct {
		Names []string `json:"names"`
		Links []*Link  `json:"links"`
	}

	resp := &sheetResponse{}

	for i := range sheet.Sheets {
		resp.Names = append(resp.Names, sheet.Sheets[i].Properties.Title)
		resp.Links = append(resp.Links, sheetLink(sheet.Sheets[i].Properties.Title))
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Printf("Unable to JSON encode response : %s\n", resp)
	}
}

func (s HTTPSheetsService) sheetHandler(w http.ResponseWriter, r *http.Request) {
	split := strings.Split(r.URL.Path, "/")

	sheetName, _ := url.PathUnescape(split[1])
	sheet, err := s.getSheet(sheetName)

	if err != nil {
		s.handleError(w, sheetName, err)
		return
	}

	if len(split) > 2 && split[2] == "array" {
		err = sheet.writeJSONArray(w)
	} else {
		err = sheet.writeJSONMap(w)
	}

	if err != nil {
		log.Printf("Unable to JSON encode sheet '%s' : %s\n", err)
	}
}

func (s HTTPSheetsService) handleError(w http.ResponseWriter, sheetName string, err error) {
	w.WriteHeader(500)
	w.Header().Add("Content-Type", "application/json")
	errorMessage := fmt.Sprintf("Unable to retrieve data from sheet '%s': %v\n", sheetName, err)
	log.Printf(errorMessage)
	fmt.Fprint(w, errorMessage)
}
