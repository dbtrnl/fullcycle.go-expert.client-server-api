package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	fmt.Println("-----\nInitiating USD_BRL cotacao request...")
	defer fmt.Println("Client request finished.")
	var url = "http://localhost:8080/cotacao"
	var bid float32
	var filename = "cotacao.txt"
	var requestTimeout = 300 * time.Millisecond
	ctx, cancelCtx := context.WithTimeout(context.Background(), requestTimeout)
	defer cancelCtx()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil { fmt.Fprintf(os.Stderr, "Error while creating NewRequestWithContext(): %v\n", err); }

	res, err := http.DefaultClient.Do(req)
	if err != nil { fmt.Fprintf(os.Stderr, "Error while doing HTTP request: %v\n", err); }
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil { fmt.Fprintf(os.Stderr, "Error while reading HTTP response: %v\n", err); }

	err = json.Unmarshal(body, &bid)
	if err != nil { fmt.Fprintf(os.Stderr, "Error while unmarshalling HTTP response: %v\n", err); }
	bidToStr := strconv.FormatFloat(float64(bid), 'f', 4, 32)
	fmt.Printf("Got '%s' as result...\n", bidToStr)

	str := "Dólar: " + bidToStr + "\n"

	f, err := os.Create(filename)
	if err != nil { fmt.Fprintf(os.Stderr, "Error while trying to create %s: %v\n", filename, err); }
	defer f.Close()
	fmt.Printf("File %s created successfully!\n", filename)

	f.Write([]byte(str))
}