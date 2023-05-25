package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type EconomiaAwesomeAPICotacaoUSDBRLResponse struct {
	Usdbrl EconomiaAwesomeAPICotacaoResponse `json:"USDBRL"`
}

type EconomiaAwesomeAPICotacaoResponse struct {
	Code        string `json:"code"`
	Codein      string `json:"codein"`
	Name        string `json:"name"`
	High        string `json:"high"`
	Low         string `json:"low"`
	VarBid      string `json:"varBid"`
	PctChange   string `json:"pctChange"`
	Bid         string `json:"bid"`
	Ask         string `json:"ask"`
	Timestamp   string `json:"timestamp"`
	Create_date string `json:"create_date"`
}

func main() {
	http.HandleFunc("/cotacao", getCotacaoHandler)
	http.ListenAndServe(":8080", nil)
}

func getCotacaoHandler(w http.ResponseWriter, r *http.Request) {
	url := "https://economia.awesomeapi.com.br/json/last/USD-BRL"

	req, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while doing http.Get(): %v\n", err)
		panic(err)
	}
	defer req.Body.Close()

	res, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while doing ioutil.ReadAll(): %v\n", err)
		panic(err)
	}
	fmt.Println("string(res):", string(res))

	// Create struct and save on SQLite
	// var resJson EconomiaAwesomeAPICotacaoUSDBRLResponse
	// err = json.Unmarshal(res, &resJson)
	// if err != nil {
	// 	fmt.Println("Erro no Unmarshall do JSON do response:", err)
	// }

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(res))
}
