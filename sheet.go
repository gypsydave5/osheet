package main

import (
	"encoding/json"
	"fmt"
	"google.golang.org/api/sheets/v4"
	"net/http"
	"strings"
)

type Sheet struct {
	googleSheet *sheets.ValueRange
}

func (s Sheet) writeJSONMap(w http.ResponseWriter) error {
	keys := headers(s.googleSheet)

	var resp struct {
		Rows []map[string]interface{} `json:"rows"`
	}

	for _, row := range s.googleSheet.Values[1:] {
		r := make(map[string]interface{})
		for j := range row {
			if len(keys[j]) > 0 {
				r[keys[j]] = row[j]
			}
		}
		resp.Rows = append(resp.Rows, r)
	}

	return json.NewEncoder(w).Encode(resp)
}

func (s Sheet) writeJSONArray(w http.ResponseWriter) error {
	return json.NewEncoder(w).Encode(s.googleSheet.Values)
}

func headers(sheet *sheets.ValueRange) []string {
	var result []string
	for i := range sheet.Values[0] {
		result = append(result, cleanHeader(fmt.Sprintf("%s", sheet.Values[0][i])))
	}
	return result
}

func cleanHeader(s string) string {
	return strings.Replace(strings.ToLower(s), " ", "-", -1)
}
