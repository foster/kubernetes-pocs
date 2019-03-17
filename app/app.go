package main

import (
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

const LISTEN_ADDR = ":3000"
const CONTROL_LISTEN_ADDR = "127.0.0.1:2999"

// hello
func logTime() {
	for {
		delay := time.Duration(rand.Intn(30) + 45)
		<-time.After(delay * time.Second)
		timeStr := time.Now().Format("3:04 PM")
		fmt.Println("It's", timeStr, "and all is well.")
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	addr := conn.RemoteAddr()
	host, _, _ := net.SplitHostPort(addr.String())
	fmt.Println("Accepted connection from", host)

	for {
		timeout := time.Now().Add(100 * time.Millisecond)
		buffer := make([]byte, 1)
		conn.SetReadDeadline(timeout)
		_, err := conn.Read(buffer)
		if err == io.EOF {
			fmt.Println("Incoming connection from", host, "closed.")
			return
		}
	}
}

func listen() {
	ln, err := net.Listen("tcp", LISTEN_ADDR)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	fmt.Println("Listening on " + LISTEN_ADDR)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn)
	}
}

func listenForControl() {
	ln, err := net.Listen("tcp", CONTROL_LISTEN_ADDR)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	fmt.Println("Listening on " + CONTROL_LISTEN_ADDR)
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting control connection:", err.Error())
			conn.Close()
			continue
		}

		timeout := time.Now().Add(1 * time.Second)
		buffer := make([]byte, 100)

		conn.SetReadDeadline(timeout)
		n, err := conn.Read(buffer)
		conn.Close()
		if err != nil {
			fmt.Println("Error reading from control connection:", err.Error())
			continue
		}
		remoteAddr := strings.TrimSpace(string(buffer[:n]))
		go dial(remoteAddr)
	}
}

func dial(address string) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Error dialing connection to", address, ":", err.Error())
		return
	}
	defer conn.Close()

	fmt.Println("Established outbound connection to", address)
	for {
		timeout := time.Now().Add(100 * time.Millisecond)
		buffer := make([]byte, 1)
		conn.SetReadDeadline(timeout)
		_, err = conn.Read(buffer)
		if err == io.EOF {
			fmt.Println("Outgoing connection to", address, "closed.")
			return
		}
	}
}

func main() {
	go logTime()
	go listenForControl()
	listen()
}
