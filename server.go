package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

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
	/* 
		db, err := sql.Open("sqlite", "root:root@tcp(localhost:9000)/root/db/db.sqlite3")
		Using docker host always results in this error: "unable to open database file: out of memory (14)"
	*/
	db, err := sql.Open("sqlite", "./db/db.sqlite3")
	if err != nil { return nil, err }
	return db, nil
}

func getCotacaoHandler(w http.ResponseWriter, r *http.Request) {
	var resJson EconomiaAwesomeAPICotacaoUSDBRLResponse
	var url = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

	req, err := http.Get(url)
	if err != nil { fmt.Fprintf(os.Stderr, "Error while doing http.Get(): %v\n", err); panic(err) }
	defer req.Body.Close()

	res, err := ioutil.ReadAll(req.Body)
	if err != nil { fmt.Fprintf(os.Stderr, "Error while doing ioutil.ReadAll(): %v\n", err); panic(err) }

	dbConnection, err := initDbConnection()
	if err != nil { fmt.Fprintf(os.Stderr, "Error while initializing database connection: %v\n", err); panic(err) }
	defer dbConnection.Close()

	err = pingDbConnection(dbConnection)
	if err != nil { fmt.Fprintf(os.Stderr, "Error while pinging DB connection: %v\n", err); panic(err) }
	
	err = json.Unmarshal(res, &resJson)
	if err != nil { fmt.Fprintf(os.Stderr, "Error while Unmarshalling JSON response: %v\n", err); panic(err) }

	err = saveCotacao(dbConnection, resJson.Usdbrl)
	if err != nil { fmt.Fprintf(os.Stderr, "Error while saving JSON response to Database: %v\n", err); panic(err) }

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(res))
}

func pingDbConnection(db *sql.DB) error {
	err := db.Ping()
	if err != nil { return err }
	return nil
}

func saveCotacao(db *sql.DB, c EconomiaAwesomeAPICotacaoObject) error {
	var query = "INSERT INTO COTACAO_USD_BRL(id, code, codein, name, high, low, var_bid, pct_change, bid, ask, timestamp, create_date) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)"
	stmt, err := db.Prepare(query)
	if err != nil { return err }
	defer stmt.Close()
	
	res, err := stmt.Exec(nil, c.Code, c.Codein, c.Name, c.High, c.Low, c.VarBid, c.PctChange, c.Bid, c.Ask, c.Timestamp, c.Create_date)
	if err != nil { return err }
	fmt.Println(res.LastInsertId())
	return nil
}

/* Not used since test table was removed
func testDbConnection(db *sql.DB) error {
	query, err := db.Query("SELECT * FROM test;")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while querying database: %v\n", err)
		return err
	}
	println(query)
	return nil
}
*/