package ping_pong_client

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"syscall"
)

var (
	Info *log.Logger
	Err  *log.Logger
)

func init() {
	Info = log.New(os.Stdout,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Err = log.New(os.Stderr,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

type Client interface {
	RunTest()
}

type TCPClient struct {
	LoaderAddr      string
	LoaderPort      int
	LocalAddr       string
	LocalPort       int
	ConnectionCount int
	BlockSize       int
	MinTimeout      int
	MaxTimeout      int
	Runtime         int
	Times           []syscall.Tms
	sync.WaitGroup
}

func NewClient(loaderAddr string, loaderPort, connectionCount, msize, minTimeout,
	       maxTimeout int, localAddr string, localPort int, runtime int) *TCPClient {
	return &TCPClient{
		LoaderAddr: loaderAddr,
		LoaderPort: loaderPort,
		ConnectionCount:      connectionCount,
		BlockSize:  msize,
		LocalAddr:  localAddr,
		LocalPort:  localAddr,
		Runtime:    runtime,
	}
}

func (client TCPClient) RunTest() {
	Info.Printf("Start load test on %s:%d", client.LoaderAddr, client.LoaderPort)
	servAddr := fmt.Sprintf("%s:%d", client.LoaderAddr, client.LoaderPort)
	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)

	if err != nil {
		Err.Printf("Error while resolving tcp addr %s",
			err.Error())
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)

	defer conn.Close()
	result := make([]byte, 1024*64)

	if err != nil {
		Err.Printf("Error %s during establishing connection to %s",
			err.Error(), client.LoaderAddr)
		return
	}

	readyFunc := func() {
		conn.Write(fmt.Sprintf("%s %d %d "+
			"%d %d %d %d", client.LocalAddr, client.LocalPort,
			client.ConnectionCount, client.Runtime, client.MinTimeout,
			client.MaxTimeout, client.BlockSize))
	}

	client.doLoad(readyFunc)
	_, err = conn.Read(result)

	if err != nil {
		Err.Printf("Error %s while gathering results from server",
			err.Error())
		return
	}
}

// Create tcp server on localAddr:localPort and waits for count of
// conection to be established, then starts load on server loaderAddr:loaderPort
func (client TCPClient) doLoad(readyFunc func()) {
	servAddr := fmt.Sprintf("%s:%d", client.LocalAddr, client.LocalPort)
	tcpAddr, _ := net.ResolveTCPAddr("tcp", servAddr)
	masterSocket, err := net.ListenTCP("tcp", tcpAddr)

	if err != nil {
		Err.Printf("Error while opening master socket %s", err.Error())
	}
	// Signal server to open connections to the client
	readyFunc()
	// Prepare all connections
	client.prepare(masterSocket)
	masterSocket.Close()
	// Save time spent on test
	client.getTime()
	// Wait until all workers are finished
	client.Done()
	client.getTime()
}

func (client TCPClient) worker(conn net.TCPConn) {
	defer conn.Close()
	defer client.Done()
	buffer := make([]byte, client.BlockSize)

	for {
		count, err := conn.Read(buffer)

		if err != nil {
			Err.Printf("Error while reading from socker %s",
				err.Error())
			break
		}

		if count == 0 {
			break
		} else if count < client.BlockSize {
			Err.Printf("Partial message")
		}
		conn.Write(buffer)
	}
}

// Wait for count connections to be established
func (client TCPClient) prepare(masterSocket net.TCPListener) {
	for i := 0; i < client.ConnectionCount; i++ {
		conn, err := masterSocket.Accept()

		if err != nil {
			Err.Printf("Error while accepting connection from server")
		}
		client.Add(1)
		go client.worker(conn)
	}
}

// Save information about current time
func (client TCPClient) getTime() {
	var tms syscall.Tms
	syscall.Times(&tms)
	client.Times = append(client.Times, tms)
}
