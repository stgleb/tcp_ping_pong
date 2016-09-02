package main

import (
	"log"
	"net"
	"time"
	"sync"
)

var wg sync.WaitGroup

func Server() {
	log.Println("Start Server")
	log.Println("Start Client")
	servAddr := "0.0.0.0:8000"
	tcpAddr, _ := net.ResolveTCPAddr("tcp", servAddr)
	listener, err := net.ListenTCP("tcp", tcpAddr)

	if err != nil {
		log.Printf("Error %s", err.Error())
	}

	wg.Done()
	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Printf("Error while accepting connection %s",
				err.Error())
		}
		log.Println("Accept connection")
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	buffer := make([]byte, 1024)
	log.Println("Handle connection")

	for {
		n, err := conn.Read(buffer)

		if n == 0 {
			continue
		}

		log.Printf("Server read %s", buffer)

		if err != nil {
			log.Printf("Error while reading from connection the server %s",
				err.Error())
		}

		log.Printf("Server write %s", buffer)
		_, err = conn.Write(buffer)

		if err != nil {
			log.Printf("Error while writing to connection the server %s",
				err.Error())
		}
	}
}

func Client() {
	wg.Wait()
	log.Println("Start Client")
	servAddr := "0.0.0.0:8000"
	tcpAddr, _ := net.ResolveTCPAddr("tcp", servAddr)

	for {
		conn, err := net.DialTCP("tcp", nil, tcpAddr)

		if err != nil {
			log.Printf("Error %s", err.Error())
		}

		buffer := []byte("hello")

		log.Printf("Client write %s", buffer)
		_, err = conn.Write(buffer)

		if err != nil {
			log.Printf("error while writing %s", err.Error())
		}

		_, err = conn.Read(buffer)
		log.Printf("Client read %s", buffer)

		if err != nil {
			log.Printf("error while reading %s", err.Error())
		}
	}
}

func main() {
	wg.Add(1)
	go Client()
	go Server()
	time.Sleep(1 * time.Second)
}
