package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Currency struct {
	ID          string `json:"id"`
	FullName    string `json:"fullName"`
	Ask         string `json:"ask"`
	Bid         string `json:"bid"`
	Last        string `json:"last"`
	Open        string `json:"open"`
	Low         string `json:"low"`
	High        string `json:"high"`
	FeeCurrency string `json:"feeCurrency"`
	Symbol      string `json:"symbol"`
}

type currencyHandlers struct {
	store map[string]Currency
}

func main() {
	currencyHandler := newCurrencyHandlers()
	currencyHandler.getCurrencyValue()

	http.HandleFunc("/currency/all", currencyHandler.getAllValues)
	http.HandleFunc("/currency/", currencyHandler.getCurrency)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func (h *currencyHandlers) getCurrencyValue() {
	response, err := http.Get("https://api.hitbtc.com/api/2/public/ticker")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	currencies := make([]Currency, 0)
	err = json.Unmarshal(responseData, &currencies)
	for _, currency := range currencies {
		if currency.Symbol == "BTCUSD" || currency.Symbol == "ETHBTC" {
			currency.ID = currency.Symbol
			h.store[currency.ID] = currency
		}
	}
}

func (h *currencyHandlers) getAllValues(w http.ResponseWriter, r *http.Request) {
	currencies := make([]Currency, len(h.store))

	i := 0
	for _, currency := range h.store {
		currencies[i] = currency
		i++
	}

	jsonBytes, err := json.Marshal(currencies)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.Write(jsonBytes)
}

func (h *currencyHandlers) getCurrency(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	currency, ok := h.store[parts[2]]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if parts[2] == "ETHBTC" || parts[2] == "BTCUSD" {
		jsonBytes, err := json.Marshal(currency)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		w.Write(jsonBytes)
	}
}

func newCurrencyHandlers() *currencyHandlers {
	return &currencyHandlers{
		store: map[string]Currency{},
	}
}
