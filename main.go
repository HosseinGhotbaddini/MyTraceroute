package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

const (
	RECV_ADDR = "0.0.0.0"
	MAX_QUERY = 3
	MAX_DIS   = 30
	TIMEOUT   = 3000 //ms
)

type Req struct {
	SourceAddr net.IP
	Time       time.Duration
}

func main() {
	//find the ip
	targetAddr, err := GetTargetAddr()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(targetAddr)

	timeout, _ := time.ParseDuration(strconv.Itoa(TIMEOUT) + "ms")

	//find hops and print them
	_, err = TraceRoute(targetAddr, MAX_QUERY, MAX_DIS, timeout)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

//get one of the addresses of a host name
func GetTargetAddr() (targetAddr net.IP, err error) {
	var target string
	var ipList []net.IP

	fmt.Print("Enter the target host (ipv4 or name): ")
	fmt.Scanln(&target)

	ipList, err = net.LookupIP(target)
	if err == nil {
		targetAddr = ipList[0]
	}
	return
}

func TraceRoute(targetAddr net.IP, maxQuery int, maxDis int, timeout time.Duration) (hopsAddrList []net.IP, err error) {
	//send packet for ttl from 1 to maxdis
Loop:
	for i := 1; i <= maxDis; i++ {
		var thisReq Req

		//send packet for a few more times if it was not successful
		for query := 0; query < maxQuery; query++ {
			thisReq, err = Request(targetAddr, i, timeout)
			if err != nil {
				return
			}

			if thisReq.SourceAddr != nil {
				break
			}
		}

		//add another hop to list
		hopsAddrList = append(hopsAddrList, thisReq.SourceAddr)

		//print it
		fmt.Print(i)
		fmt.Print(" : ")

		if hopsAddrList[i-1] == nil {
			fmt.Println("*")
		} else {
			fmt.Print(thisReq.Time)
			fmt.Print(" ")
			fmt.Println(hopsAddrList[i-1])
			if targetAddr.Equal(hopsAddrList[i-1]) {
				fmt.Println("FIN")
				break Loop
			}
		}
	}
	return
}

func Request(targetAddr net.IP, ttl int, timeout time.Duration) (thisReq Req, err error) {

	//make connection to send packet
	sndConnection, err := net.Dial("ip4:icmp", targetAddr.String())
	if err != nil {
		return
	}
	defer sndConnection.Close()

	newConnection := ipv4.NewConn(sndConnection)
	err = newConnection.SetTTL(ttl)
	if err != nil {
		return
	}

	//make packet
	echo := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID:   rand.Int(),
			Seq:  1,
			Data: []byte("TABS"),
		}}

	request, err := echo.Marshal(nil)
	if err != nil {
		return
	}

	//make receive connection
	recvConnection, err := icmp.ListenPacket("ip4:icmp", RECV_ADDR)
	if err != nil {
		return
	}
	defer recvConnection.Close()

	sTime := time.Now()

	//send packet
	_, err = sndConnection.Write(request)
	if err != nil {
		return
	}

	err = recvConnection.SetDeadline(time.Now().Add(timeout))
	if err != nil {
		return
	}

	//receive response packet
	recvData := make([]byte, 1500)
	_, addr, errr := recvConnection.ReadFrom(recvData)

	thisReq.Time = time.Since(sTime)
	if errr == nil {
		thisReq.SourceAddr = net.ParseIP(addr.String())
	}

	return thisReq, err
}
