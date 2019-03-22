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

const LISTEN_PORT = ":3000"
const CONTROL_LISTEN_ADDR = "localhost:2999"

var MY_IP string = os.Getenv("MY_POD_IP")
var peers map[string]bool = make(map[string]bool)

func logTime() {
	for {
		delay := time.Duration(rand.Intn(30) + 45)
		<-time.After(delay * time.Second)
		timeStr := time.Now().Format("3:04 PM")
		fmt.Printf("[%s] It's %s and all is well.\n", MY_IP, timeStr)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	addr := conn.RemoteAddr()
	host, _, _ := net.SplitHostPort(addr.String())
	peers[host] = true
	fmt.Printf("Accepted connection %s <== %s\n", MY_IP, host)

	for {
		timeout := time.Now().Add(100 * time.Millisecond)
		buffer := make([]byte, 1)
		conn.SetReadDeadline(timeout)
		_, err := conn.Read(buffer)
		if err == io.EOF {
			delete(peers, host)
			fmt.Println("Incoming connection from", host, "closed.")
			return
		}
	}
}

func listen() {
	listenAddr := MY_IP + LISTEN_PORT
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		panic(fmt.Sprintf("Error listening: %v", err.Error()))
	}

	fmt.Println("Listening on " + listenAddr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(fmt.Sprintf("Error accepting: %v", err.Error()))
		}
		go handleConnection(conn)
	}
}

func listenForControl() {
	ln, err := net.Listen("tcp", CONTROL_LISTEN_ADDR)
	if err != nil {
		panic(fmt.Sprintf("Error listening for control connection: %v", err.Error()))
	}

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
	host, _, _ := net.SplitHostPort(address)

	// do not connect to peers that already connected to us
	if _, exists := peers[host]; exists {
		return
	}

	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Error dialing connection to", address, ":", err.Error())
		return
	}
	defer conn.Close()

	peers[host] = true
	fmt.Printf("Established connection %s ==> %s\n", MY_IP, address)
	for {
		timeout := time.Now().Add(100 * time.Millisecond)
		buffer := make([]byte, 1)
		conn.SetReadDeadline(timeout)
		_, err = conn.Read(buffer)
		if err == io.EOF {
			fmt.Println("Outgoing connection to", address, "closed.")
			delete(peers, host)
			return
		}
	}
}

func main() {
	go logTime()
	go listenForControl()
	listen()
}
