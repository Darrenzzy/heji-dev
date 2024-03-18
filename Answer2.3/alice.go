// Alice.go
package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	address      = "localhost:8080"
	logFileAlice = "alice.log"
)

var (
	latestPrices sync.Map // 使用sync.Map来存储最新的成交价格

)

func main() {
	// 设置日志文件
	setupLoggingAlice()
	// 从CSV文件中加载数据
	loadData("transaction.1min.csv")

	// 监听TCP连接
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Printf("Error listening on %s: %s\n", address, err)
		return
	}
	defer listener.Close()
	log.Printf("Waiting for Bob \n")

	// 接受连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %s\n", err)
			continue
		}
		go handleConnection(conn)
	}
}
func setupLoggingAlice() {
	file, err := os.OpenFile(logFileAlice, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %s", err)
	}
	log.SetOutput(file)
}

// 从CSV文件中读取数据并更新最新的成交价格
func loadData(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening file: %s\n", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	// 第一行header
	reader.Read()
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error reading row: %s\n", err)
			continue
		}

		windCode := record[0]
		price, err := strconv.Atoi(record[4])
		if err != nil {
			log.Printf("Error parsing price for %s: %s\n", windCode, err)
			continue
		}
		realPrice := float64(price) / 10000
		latestPrices.Store(windCode, realPrice)
	}
}

// 处理每个Bob的TCP连接
func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		windCode, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading from connection: %s\n", err)
			}
			return
		}

		if strings.HasPrefix(windCode, "byb") {
			log.Println("Byebye Bob", strings.Split(windCode, " ")[1])
			log.Println("Waiting for Bob")
			continue
		}
		if strings.HasPrefix(windCode, "probe") {
			log.Println("Bob", strings.Split(windCode, " ")[1], "is coming")
			continue
		}
		windCode = strings.TrimSpace(windCode)
		log.Println(windCode)
		if price, ok := latestPrices.Load(windCode); ok {
			fmt.Fprintf(conn, "%s Price: %.2f\n", windCode, price)
		} else {
			fmt.Fprintf(conn, "%s not found\n", windCode)
		}
	}
}
