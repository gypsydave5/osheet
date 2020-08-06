package main

import (
	"github.com/patrickmn/go-cache"
	"time"
)
import "gopkg.in/matryer/try.v1"
import "github.com/jpillora/backoff"

type InMemCachingSheetService struct {
	ss *SheetsService
	c  *cache.Cache
}

func (s InMemCachingSheetService) getSheet(sheetName string) (sheet *Sheet, err error) {
	cacheSheet, found := s.c.Get(sheetName)
	if found {
		aSheet := cacheSheet.(Sheet)
		return &aSheet, nil
	}

	b := backoff.Backoff{}
	err = try.Do(func(attempt int) (retry bool, err error) {
		sheet, err = s.ss.getSheet(sheetName)
		time.Sleep(b.Duration())
		return true, err
	})

	return sheet, err
}
