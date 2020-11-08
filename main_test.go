// main_test.go
package main

import (
	"log"
	"net"
	"testing"
	"time"
)

func TestLocal(test *testing.T) {
	ip := net.ParseIP("127.0.0.1")
	_, err := TraceRoute(ip, MAX_QUERY, MAX_DIS, 10*time.Second)
	if err != nil {
		test.Error(err)
		log.Print(err)
		return
	}
}

func TestBroadcast(test *testing.T) {
	ip := net.ParseIP("255.255.255.255")
	_, err := TraceRoute(ip, MAX_QUERY, MAX_DIS, 10*time.Second)
	if err == nil {
		test.Error(err)
		log.Print(err)
		return
	}
}

func TestDNS(test *testing.T) {
	ip := net.ParseIP("8.8.8.8")
	_, err := TraceRoute(ip, MAX_QUERY, MAX_DIS, 10*time.Second)
	if err != nil {
		test.Error(err)
		log.Print(err)
		return
	}
}
