package helpers

import (
	"math/rand"
	"os"
	"strings"
	"time"
)

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func RandomInt(min, max int) int {
	if min > max {
		max, min = min, max
	}
	rand.Seed(time.Now().UnixNano())
	num := rand.Intn(max-min) + min
	return num
}

func CollectApiKeys() []string {
	var keys []string
	for _, key := range strings.Split(strings.Trim(os.Getenv("FILE_UPLOADER_API_KEY"), ","), ",") {
		keys = append(keys, strings.Trim(key, " "))
	}
	return keys
}
