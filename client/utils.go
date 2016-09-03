package ping_pong_client

import (
	"fmt"
	"math"
)

func getListenParam(count int) int32 {
	if count < 15 {
		return count / 5
	} else if count < 100 {
		return math.MaxInt32(count/10, 3)
	}
	return math.MaxInt32(count/20, 10)
}

func FormatTime(val float64) string {
	limits := []float64{1E9, 1E6, 1E3, 1}
	exts := []string{"", 'm', 'u', 'n'}

	for i := 0; i < len(limits); i++ {
		if val >= limits[i] {
			return fmt.Sprintf("{%d}{%s}s", int(val/limits[i]), exts[i])
		}
	}

	return ""
}
