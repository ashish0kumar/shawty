package utils

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"time"
)

func GetShortCode() string {
	fmt.Println("Shortening URL")
	ts := time.Now().UnixNano()
	fmt.Println("Timestamp: ", ts)

	// convert the timestamp to byte slice and encode it to base64
	ts_bytes := strconv.AppendInt(nil, ts, 10)
	key := base64.StdEncoding.EncodeToString(ts_bytes)
	fmt.Println("Key: ", key)

	// remove the last 2 chars since they are always '=='
	key = key[:len(key)-2]

	// return the last chars after 16 chars
	return key[16:]
}
