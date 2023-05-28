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

// func initDbConnection() (*sql.DB, error) {
// 	// db, err := sql.Open("sqlite", "root:root@tcp(localhost:9000)/root/db/db.sqlite3") //Qualquer operação resulta em unable to open database file: out of memory (14) 
// 	db, err := sql.Open("sqlite", "./db/db.sqlite3")
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Error while connecting to SQLite database: %v\n", err)
// 		return nil, err
// 	}
// 	query, err := db.Query("SELECT * FROM test;")
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Error while querying database: %v\n", err)
// 		return nil, err
// 	}
// 	println(query)
// 	return db, nil
// }

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

	// db, err := sql.Open("sqlite", "root:root@tcp(localhost:9000)/root/db/db.sqlite3") //Qualquer operação resulta em unable to open database file: out of memory (14) 
	db, err := sql.Open("sqlite", "./db/db.sqlite3")
	// if err != nil {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while connecting to SQLite database: %v\n", err)
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while pinging DB connection: %v\n", err)
		panic(err)
	}

	/* Test
	queryRes, err := db.Exec("INSERT INTO test VALUES(NULL, 'testeGolang', 69);")
	if err != nil {
		fmt.Fprintf(os.Stderr, "-----\nError while querying database: %v\n-----\n", err)
		panic(err)
	}
	println(queryRes.RowsAffected())
 	*/
	
	var resJson EconomiaAwesomeAPICotacaoUSDBRLResponse
	err = json.Unmarshal(res, &resJson)
	if err != nil {
		fmt.Println("Erro no Unmarshall do JSON do response:", err)
	}
	err = saveCotacao(db, resJson.Usdbrl)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(res))
}

func saveCotacao(db *sql.DB, c EconomiaAwesomeAPICotacaoObject) error {
	stmt, err := db.Prepare("INSERT INTO COTACAO_USD_BRL(id, code, codein, name, high, low, var_bid, pct_change, bid, ask, timestamp, create_date) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)")
	if err != nil { return err }
	defer stmt.Close()
	res, err := stmt.Exec(nil, c.Code, c.Codein, c.Name, c.High, c.Low, c.VarBid, c.PctChange, c.Bid, c.Ask, c.Timestamp, c.Create_date)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while writing to database: %v\n", err)
		return err
	}
	fmt.Println(res.LastInsertId())
	return nil
}