package main

import (
	"reflect"
	"testing"
)

func TestRandNum(t *testing.T) {
	var x string
	check := randNum()
	if reflect.TypeOf(check) != reflect.TypeOf(x) {
		t.Error("fuck")
	}
}

func TestUrlReturn(t *testing.T) {
	var (
		engine      = "currency"
		market      = "index"
		typeOfCheck = "marketdata"
		answer      = "http://iss.moex.com/iss/engines/currency/markets/index/securities.json?iss.only=marketdata&sort_order=desc&first=1&marketdata.columns=UPDATETIME&sort_column=UPDATETIME&securities=MICEXINDEXCF,RTSI,MICEXBMI,RTSSTD,RVI"
	)
	checkAnswer := urlReturn(engine, market, typeOfCheck)
	if checkAnswer != answer {
		t.Error("Expected 1.5, got ", checkAnswer)
	}
}
