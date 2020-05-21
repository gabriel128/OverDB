package dialers

import (
	"net/rpc"
	"log"
	"time"
	"strconv"
)

func DialTCP(port int) (*rpc.Client, error) {
	var client *rpc.Client
	var err error

	for {
		client, err = rpc.Dial("tcp", "localhost:" + strconv.Itoa(port))

		if err != nil {
			log.Println("Dial for port", strconv.Itoa(port), "failed", err)
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}
	return client, err
}

func DialHttp(port int) (*rpc.Client, error) {
	var client *rpc.Client
	var err error

	for {
		client, err = rpc.DialHTTP("tcp", "localhost:" + strconv.Itoa(port))

		if err != nil {
			// log.Println("Dial for port", strconv.Itoa(port), "failed", err)
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}
	return client, err
}
