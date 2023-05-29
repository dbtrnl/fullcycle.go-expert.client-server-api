package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

type EconomiaAwesomeAPICotacaoUSDBRLResponse struct {
	Usdbrl EconomiaAwesomeAPICotacaoObject `json:"USDBRL"`
}

type EconomiaAwesomeAPICotacaoObject struct {
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

func initDbConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./db/db.sqlite3")
	if err != nil { return nil, err }

	err = pingDbConnection(db)
	if err != nil { fmt.Fprintf(os.Stderr, "Error while pinging DB connection: %v\n", err); panic(err) }
	return db, nil
}

func getCotacaoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("-----\ngetCotacaoHandler initiating...")
	defer fmt.Println("getCotacaoHandler finished.")
	
	var resJson EconomiaAwesomeAPICotacaoUSDBRLResponse
	var url = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

	var getCotacaoTimeout = 200 * time.Millisecond
	httpCtx, cancelHttpCtx := context.WithTimeout(context.Background(), getCotacaoTimeout)
	defer cancelHttpCtx()

	req, err := http.NewRequestWithContext(httpCtx, "GET", url, nil)
	if err != nil { fmt.Fprintf(os.Stderr, "Error while creating NewRequestWithContext(): %v\n", err); panic(err) }
	
	res, err := http.DefaultClient.Do(req)
	if err != nil { fmt.Fprintf(os.Stderr, "Error while doing http.DefaultClient.Do(): %v\n", err); panic(err) }
	defer res.Body.Close()
	fmt.Printf("Get request to %s sucessful.\n", url)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil { fmt.Fprintf(os.Stderr, "Error while doing ioutil.ReadAll(): %v\n", err); panic(err) }

	err = json.Unmarshal(body, &resJson)
	if err != nil { fmt.Fprintf(os.Stderr, "Error while Unmarshalling JSON response: %v\n", err); panic(err) }

	dbConnection, err := initDbConnection()
	if err != nil { fmt.Fprintf(os.Stderr, "Error while initializing database connection: %v\n", err); panic(err) }
	defer dbConnection.Close()

	err = saveCotacao(dbConnection, resJson.Usdbrl)
	if err != nil { fmt.Fprintf(os.Stderr, "Error while saving JSON response to Database: %v\n", err); panic(err) }
	fmt.Println("Data sucessfully written to database.")

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(resJson.Usdbrl.Bid))
}

func pingDbConnection(db *sql.DB) error {
	err := db.Ping()
	if err != nil { return err }
	return nil
}

func saveCotacao(db *sql.DB, c EconomiaAwesomeAPICotacaoObject) error {
	var saveCotacaoTimeout = 10 * time.Millisecond
	ctx, cancelCtx := context.WithTimeout(context.Background(), saveCotacaoTimeout)
	defer cancelCtx()
	
	var query = "INSERT INTO COTACAO_USD_BRL(id, code, codein, name, high, low, var_bid, pct_change, bid, ask, timestamp, create_date) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)"
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil { return err }
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, nil, c.Code, c.Codein, c.Name, c.High, c.Low, c.VarBid, c.PctChange, c.Bid, c.Ask, c.Timestamp, c.Create_date)
	if err != nil { return err }
	return nil
}
