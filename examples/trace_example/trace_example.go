package main

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateTraceID creates a new 16-byte trace ID as a hex string
func GenerateTraceID() string {
	rand.Seed(time.Now().UnixNano())
	bytes := make([]byte, 16)
	rand.Read(bytes)

	// Convert to hex string
	hexStr := fmt.Sprintf("%x", bytes)
	return hexStr
}

// GenerateSpanID creates a new 8-byte span ID as a hex string
func GenerateSpanID() string {
	rand.Seed(time.Now().UnixNano())
	bytes := make([]byte, 8)
	rand.Read(bytes)

	// Convert to hex string
	hexStr := fmt.Sprintf("%x", bytes)
	return hexStr
}

func main() {
	fmt.Println("TraceID and SpanID Examples:")
	fmt.Println("=============================")

	for i := 0; i < 5; i++ {
		traceID := GenerateTraceID()
		spanID := GenerateSpanID()

		fmt.Printf("Example %d:\n", i+1)
		fmt.Printf("  TraceID: %s\n", traceID)
		fmt.Printf("  SpanID:  %s\n", spanID)
		fmt.Println()
	}

	// Show the typical format requirements
	fmt.Println("Format Specifications:")
	fmt.Println("======================")
	fmt.Println("- TraceID: 32-character lowercase hex string (16 bytes)")
	fmt.Println("- SpanID:  16-character lowercase hex string (8 bytes)")
}