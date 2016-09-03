package main

import (
	"../../ping-pong-client/client"
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	fileName        string
	loaderAddr      string
	loaderPort      int
	localAddr       string
	localPort       int
	connectionCount int
	blockSize       int
	runtime         int
	minTimeout      int
	maxTimeout      int
)

func init() {
	flag.StringVar(&loaderAddr, "loaderAddr", "0.0.0.0", "loader ip address")
	flag.IntVar(&loaderPort, "loaderPort", 33331, "loader port")
	flag.StringVar(&localAddr, "localAddr", "0.0.0.0", "local address")
	flag.IntVar(&localPort, "localPort", 33332, "local port to receive connections")
	flag.IntVar(&connectionCount, "connectionCount", 1000, "count of connections")
	flag.IntVar(&blockSize, "blockSize", 1024, "size of block in bytes")
	flag.IntVar(&runtime, "runtime", 30, "time for load test")
	flag.IntVar(&minTimeout, "minTimeout", 10, "minimal timeout for server")
	flag.IntVar(&maxTimeout, "maxTimeout", 10, "maximal timeout for server")
	flag.StringVar(&fileName, "fileName", "output.txt", "filename to store results")
}

func main() {
	client := ping_pong_client.NewClient(loaderAddr, loaderPort, localAddr, localPort,
		connectionCount, blockSize, runtime, minTimeout, maxTimeout)
	stats := client.RunTest()
	percentiles := []float64{0.5, 0.75, 0.95}
	perc50, perc75, perc95 := ping_pong_client.GetLatencies(stats.LatencyDistribution,
		stats.LatencyBase,
		percentiles)

	f, err := os.Create(fileName)
	defer f.Close()

	if err != nil {
		log.Printf("Error while opening the file %s", err.Error())
	}

	// Print result data to file, to be proceed later
	f.WriteString(fmt.Sprintf("%d %d %d %s %s %s %e %e %d", stats.Utime, stats.Stime, stats.Ctime,
		ping_pong_client.FormatTime(perc50), ping_pong_client.FormatTime(perc75),
		ping_pong_client.FormatTime(perc95), stats.Percentiles[0],
		stats.Percentiles[len(stats.Percentiles)-1], stats.MessageCount))
}
