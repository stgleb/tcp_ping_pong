package ping_pong_client

import "math"

func getListenParam(count int) int {
	if count < 15 {
		return count / 5
	} else if count < 100 {
		return math.MaxInt32(count/10, 3)
	}
	return math.MaxInt32(count/20, 10)
}
