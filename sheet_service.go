package main

import (
	"context"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"log"
	"net/url"
)

type SheetsService struct {
	sheetId string
	svc     *sheets.Service
}

type Link struct {
	Path string `json:"path"`
	Rel  string `json:"rel"`
}

func newSheetsService(serviceAccountJSON []byte, spreadsheetID string) (*SheetsService, error) {
	// Notes
	// 1. Create a service account on the google app whatsface thing mhttps://console.developers.google.com/apis/credentials
	// 2. download the json for the service account and call it service_account.json or whatever

	// 3. give the email address of the service account access to the sheet
	// 4. find the id of the sheet

	// this creates a JWT config thing from the service account JSON
	config, err := google.JWTConfigFromJSON(serviceAccountJSON, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	// config can make an HTTP client with some context lol
	client := config.Client(context.TODO())

	// make a sheets client
	srv, err := sheets.NewService(context.TODO(), option.WithHTTPClient(client))
	if err != nil {
		return &SheetsService{}, fmt.Errorf("Unable to retrieve Sheets client: %v", err)
	}

	s := &SheetsService{
		sheetId: spreadsheetID,
		svc:     srv,
	}
	return s, err
}

func (s SheetsService) getSheet(sheetName string) (*Sheet, error) {
	sheet, err := s.svc.Spreadsheets.Values.Get(s.sheetId, sheetName).Do()
	return &Sheet{sheet}, err
}

func sheetLink(title string) *Link {
	return &Link{
		Rel:  title,
		Path: sheetPath(title),
	}
}

func sheetPath(title string) string {
	return "/" + url.PathEscape(title)
}
