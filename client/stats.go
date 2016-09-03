package ping_pong_client

import (
	"strings"
	"syscall"
)

type Stats struct {
	Utime               int64
	Stime               int64
	Ctime               int64
	MessageCount        int64
	LatencyBase         float64
	LatencyDistribution map[float64]float64
	Percentiles         []float64
}

func ProcessStats(data []byte, tms []syscall.Tms) Stats {
	utime := tms[1].Utime - tms[0].Utime
	stime := tms[1].Stime - tms[0].Stime
	ctime := tms[1].Cstime - tms[0].Cstime

	result := string(data[:])
	splittedResult := strings.Split(result, " ")
	msgProcessed := splittedResult[:1]
	latencyBase := splittedResult[1:2]
	latencyDistAndPercentiles := splittedResult[2:]
	percentiles, latencies := processLatenciesAndPercentiles(latencyDistAndPercentiles)

	return Stats{
		Utime:               utime,
		Stime:               stime,
		Ctime:               ctime,
		MessageCount:        msgProcessed,
		LatencyBase:         latencyBase,
		LatencyDistribution: latencies,
		Percentiles:         percentiles,
	}
}

func processLatenciesAndPercentiles(latAndPerc []string) ([]int, map[int]int) {
	lNp := make([]int, 0)

	for i := 0; i < len(latAndPerc); i++ {
		lNp = append(lNp, int(latAndPerc[i]))
	}
	latSize := lNp[0]
	latDistributionRaw := lNp[1 : 1+latSize*2]
	percetilesRaw := lNp[1+latSize*2:]
	percSize := percetilesRaw[0]

	if percSize+1 != len(percetilesRaw) {
		Err.Printf("Percentile size expected %d actual %d",
			percSize+1,
			len(percetilesRaw))
	}

	percentiles := percetilesRaw[1:]
	latencies := getLatencyMap(latDistributionRaw)

	return percentiles, latencies
}

func getLatencyMap(latency []int) map[int]int {
	result := make(map[int]int)

	for i := 0; i < len(latency); i += 2 {
		key, value := latency[i], latency[i+1]
		result[key] = value
	}

	return result
}
