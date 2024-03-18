package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	addressB = "localhost:8080"
)

func main() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)
	isClose := make(chan os.Signal, 1)
	go func() {
		if _, ok := <-ch; ok {
			fmt.Println("Bob is Exiting...")
			close(isClose)
			os.Exit(0)
			return
		}
	}()
	for {
		select {
		case <-isClose:
			return
		default:
			sendTcp(isClose)
		}

	}
}

func sendTcp(isClose chan os.Signal) {
	conn, err := net.Dial("tcp", addressB)
	if err != nil {
		fmt.Printf("Error connecting to Alice: %s\n", err)
		time.Sleep(time.Second)
		return
	}
	defer conn.Close()
	pid := os.Getpid()
	// 向Alice发送一个特殊的探测消息
	fmt.Fprintf(conn, "probe %d \n", pid)

	fmt.Println("Number of ", pid, " BoB Start, connected to Alice")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		windCode := scanner.Text()
		if strings.HasPrefix(windCode, "query") {
			fmt.Fprintf(conn, "%s \n", strings.TrimLeft(windCode, "query"))

		} else if strings.ContainsAny(windCode, "byb") {
			fmt.Fprintf(conn, "byb %d \n", pid)
			isClose <- syscall.SIGINT
			return
		}
		fmt.Fprintf(conn, " \n")
		response, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading response: %s\n", err)
			return
		}
		if strings.ContainsAny(response, "not found") {
			fmt.Print("windCode ", response)
		} else {
			fmt.Print("price =", response)
		}
	}
}
