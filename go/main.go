package main

import (
	"encoding/hex"
	"fmt"
	"time"

	avl "example.com/teltonica/teltonica"
)

func main() {
	hexString := "000000000000003608010000016B40D8EA30010000000000000000000000000000000105021503010101425E0F01F10000601A014E0000000000000000010000C7CF"

	// Decode the hex string into a byte slice
	byteSlice, err := hex.DecodeString(hexString)
	if err != nil {
		fmt.Println("Error decoding hex string:", err)
		return
	}

	// Measure the time taken for deserialization
	startTime := time.Now()
	iterations := 0
	seconds := 10.0 

	for {
		_, err := avl.NewAvlPacketTcp(byteSlice, true)
		if err != nil {
			fmt.Println("Error during deserialization:", err)
			return
		}

		iterations++

		// Check if one second has passed
		if time.Since(startTime).Seconds() >= seconds {
			break
		}
	}

	fmt.Printf("Number of deserializations in one second: %.2f\n", float64(iterations) / seconds)
}
