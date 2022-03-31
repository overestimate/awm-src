package main

import (
	"crypto/tls"
	"time"
)

// GetLatencyMc() returns the latency to api.minecraftservices.com:443, in milliseconds, as a float64.
func GetLatencyMc() (float64, error) {
	conn, err := tls.Dial("tcp", "api.minecraftservices.com:443", nil)
	if err != nil {
		return -1, err
	}

	tmpVar := make([]byte, 4096)
	payload := []byte("HEAD /minecraft/profile HTTP/1.1\r\nUser-Agent: httping.go/0.1\r\n\r\n")
	t1 := time.Now()
	conn.Write(payload)
	_, err = conn.Read(tmpVar)
	if err != nil {
		return -1, err
	}
	t2 := time.Now()
	return float64(t2.Sub(t1).Milliseconds()), nil
}
