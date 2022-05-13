package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"time"
)

// create stock data structure
type stock struct {
	Time   string  `json:"time"`
	Symbol string  `json:"symbol"`
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume int     `json:"volume"`
}

// get 10 random data for stock
func getStockData() []stock {
	var stocks []stock

	var stockSymbols = []string{"AAPL", "MSFT", "AMZN", "GOOG", "FB", "TWTR", "INTC", "CSCO", "NVDA", "ORCL"}

	for i := 0; i < 10; i++ {
		stocks = append(stocks, stock{
			Time:   "2020-01-01T00:00:00.000Z", // don't care about this value, will be updated to now
			Symbol: stockSymbols[rand.Intn(10)],
			Open:   100.00,
			High:   100.00,
			Low:    100.00,
			Close:  100.00,
			Volume: 10000,
		})
	}

	return stocks
}

func updateValues(stock *stock) {
	rand.Seed(time.Now().UnixNano())

	stock.Time = time.Now().UTC().Format(time.RFC3339)

	var percentage = rand.Float64()
	var operators = []string{"+", "-"}
	var operator = operators[rand.Intn(2)]

	switch operator {
	case "+":
		stock.Close = math.Round(stock.Close+(stock.Close*percentage)*100) / 100
	case "-":
		stock.Close = math.Round(stock.Close-(stock.Close*percentage)*100) / 100
	}

	if stock.Close > stock.High {
		stock.High = stock.Close
	}

	if stock.Close < stock.Low {
		stock.Low = stock.Close
	}

	stock.Volume = stock.Volume + rand.Intn(1000)
}

func main() {
	var address string
	var network string

	flag.StringVar(&address, "address", ":9000", "address to listen on")
	flag.StringVar(&network, "network", "tcp", "network to listen on")
	flag.Parse()

	//validate network
	switch network {
	case "tcp":
		fmt.Printf("Listening on tcp network at prot 9000\n")
	default:
		log.Fatalln("Invalid network", network)
	}

	//create listener
	listener, err := net.Listen(network, address)
	if err != nil {
		log.Fatal("failed to create listener", err)
	}

	defer listener.Close()

	// connection loop
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("failed to accept connection", err)
			continue
		}
		log.Println("accepted connection", conn.RemoteAddr())

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	if _, err := conn.Write([]byte("Connected... Stock data incoming\n")); err != nil {
		log.Println("failed to write to connection", err)
		return
	}

	// loop to stay connected
	for {
		if err := sendStockData(conn); err != nil {
			log.Println("failed to send stock data", err)
			return
		}

		// wait for 100 miliseconds
		time.Sleep(time.Millisecond * 100)
	}
}

func sendStockData(conn net.Conn) error {

	// send stock data
	stocks := getStockData()

	stock := stocks[rand.Intn(10)]
	updateValues(&stock)

	// write stock json to conn
	jsonData, _ := json.MarshalIndent(stock, "", "  ")

	if _, err := conn.Write([]byte(fmt.Sprintf(
		"%s,\n",
		string(jsonData),
	))); err != nil {
		return err
	}

	return nil
}
