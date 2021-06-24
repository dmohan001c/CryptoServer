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

type Currencies struct {
	Currency Currency
}

type currencyHandlers struct {
	store map[string]Currency
}

func (h *currencyHandlers) get(w http.ResponseWriter, r *http.Request) {
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
	//w.Header().Add("content-type", "application/json")
	//w.WriteHeader(http.StatusOK)
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
		//w.Header().Add("content-type", "application/json")
		//w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)
	}
}

func newCurrencyHandlers() *currencyHandlers {
	return &currencyHandlers{
		store: map[string]Currency{},
	}
}
func main() {
	currencyHandler := newCurrencyHandlers()
	//symbol := "ETHBTC"
	//if symbol == "ETHBTC" || symbol == "BTCUSD" {
	currencyHandler.getCurrencyValue()
	//} else {
	//fmt.Println("Unsupported Symbols", symbol)
	//}

	http.HandleFunc("/currency/all", currencyHandler.get)
	http.HandleFunc("/currency/", currencyHandler.getCurrency)
	err := http.ListenAndServe(":8102", nil)
	if err != nil {
		panic(err)
	}
}

func (h *currencyHandlers) getCurrencyValue() {
	//if symbol == "ETHBTC" || symbol == "BTCUSD" {
	//response, err := http.Get("https://api.hitbtc.com/api/2/public/ticker/" + symbol)
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
			fmt.Println(currency.Symbol)
			h.store[currency.ID] = currency
		}
	}

	//var responseObject Currency
	//json.Unmarshal(responseData, &responseObject)
	//fmt.Println(string(responseData))
	/*} else {
		fmt.Println("Unsupported Symbols", symbol)
	}*/
}
