package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

func fillTheBook(symbol string, depth int, data map[string]interface{}, conn *sql.DB) {
	bids := data["bids"].([]interface{})
	asks := data["asks"].([]interface{})

	for _, bid := range bids {
		price := bid.([]interface{})[0].(string)
		quantity := bid.([]interface{})[1].(string)

		saveOrderBookEntry("bid", price, quantity, symbol, conn)
	}

	for _, ask := range asks {
		price := ask.([]interface{})[0].(string)
		quantity := ask.([]interface{})[1].(string)

		saveOrderBookEntry("ask", price, quantity, symbol, conn)
	}
}

func saveOrderBookEntry(side string, price string, quantity string, symbol string, conn *sql.DB) {
	var priceDB, quantityDB string
	err := conn.QueryRow("SELECT price, quantity FROM "+side+" WHERE s_price = $1 AND symbol = $2", price, symbol).Scan(&priceDB, &quantityDB)
	if err == sql.ErrNoRows {
		_, err := conn.Exec("INSERT INTO "+side+" (s_price, price, quantity, symbol) VALUES ($1, $2, $3, $4)",
			price, price, quantity, symbol)
		if err != nil {
			log.Println("Error inserting into database:", err)
		}
	} else if err == nil {
		_, err := conn.Exec("UPDATE "+side+" SET quantity = $1 WHERE s_price = $2 AND symbol = $3",
			quantity, price, symbol)
		if err != nil {
			log.Println("Error updating database:", err)
		}
	} else {
		log.Println("Error querying database:", err)
	}
}

func getData(url string, symbol string, depth int, resultChannel chan<- string, wg *sync.WaitGroup, conn *sql.DB) {
	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil {
		resultChannel <- fmt.Sprintf("Error en la solicitud GET: %s", err.Error())
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		resultChannel <- fmt.Sprintf("Error al leer la respuesta: %s", err.Error())
		return
	}

	var parsedData map[string]interface{}
	err = json.Unmarshal(data, &parsedData)
	if err != nil {
		resultChannel <- fmt.Sprintf("Error al analizar la respuesta JSON: %s", err.Error())
		return
	}

	fillTheBook(symbol, depth, parsedData, conn)

	resultChannel <- "Datos procesados correctamente..."
}

func BinanceData() error {
	depth := 1000

	coins := os.Getenv("COIN_LIST")
	coinList := strings.Split(coins, ",")

	resultChannel := make(chan string)

	db := GetConnection()
	defer db.Close()

	var wg sync.WaitGroup

	for _, coin := range coinList {
		wg.Add(1)
		url := fmt.Sprintf("https://api.binance.com/api/v3/depth?symbol=%s&limit=%d", coin, depth)
		go getData(url, coin, depth, resultChannel, &wg, db)
	}

	go func() {
		wg.Wait()
		close(resultChannel)
	}()

	for result := range resultChannel {
		fmt.Println(result)
	}

	return nil
}
