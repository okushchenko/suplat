package worker

import (
	"log"
	"time"

	"github.com/alexgear/suplat/common"
	"github.com/alexgear/suplat/datastore"
	"github.com/alexgear/suplat/network"
)

var err error

func producer(c chan common.Point) {
	point, err := network.Ping()
	if err != nil {
		log.Fatalf("Failed to ping: %s", err.Error())
	}
	c <- point
}

func consumer(c chan common.Point) {
	for {
		datastore.Write(<-c)
	}
}

func InitWorker() {
	go func() {
		err = datastore.Flush()
		if err != nil {
			log.Fatalf("Failed to flush: %s", err.Error())
		}
	}()
	c := make(chan common.Point, 100000)
	go consumer(c)
	ticker := time.NewTicker(200 * time.Millisecond)
	go func() {
		for _ = range ticker.C {
			go producer(c)
		}
	}()
}
