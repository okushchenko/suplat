package main

import (
	"flag"
	"log"
	"strconv"
	"time"

	old_datastore "github.com/alexgear/checker/datastore"
	"github.com/alexgear/suplat/common"
	"github.com/alexgear/suplat/config"
	"github.com/alexgear/suplat/datastore"
	"github.com/alexgear/suplat/worker"
)

var err error

func Migrate() {
	iefs := []string{"lan", "wifi"}
	go func() {
		err = datastore.Flush()
		if err != nil {
			log.Fatalf("Failed to flush: %s", err.Error())
		}
	}()
	for _, ief := range iefs {
		log.Printf("Migrating %s...\n", ief)
		data, err := old_datastore.Read(ief, 10000*time.Hour)
		if err != nil {
			log.Fatalf("Failed to read: %s", err.Error())
		}
		log.Println("Done reading")
		var ops int
		ticker := time.NewTicker(5 * time.Second)
		go func() {
			for _ = range ticker.C {
				log.Printf("Writing %d points per second\n", ops/5)
				ops = 0
			}
		}()
		for t, status := range data {
			ops += 1
			point := common.Point{
				Time:      t,
				Interface: ief,
			}
			p99 := strconv.FormatFloat(status.Percentile99, 'f', -1, 64)
			point.Latency, err = time.ParseDuration(p99 + "s")
			if err != nil {
				log.Fatalf("Failed to parse duration:", err.Error())
			}

			if status.Uptime > 99.0 {
				point.Error = 0
			} else {
				point.Error = 1
			}
			datastore.Write(point)
		}
		ticker.Stop()
	}
}

func main() {
	log.Printf("Suplat version %s, build %s", config.Version, config.BuildTime)
	log.Println("Load flags...")
	var migrate bool
	flag.BoolVar(&migrate, "migrate", false, "help")
	flag.Parse()
	log.Println("Load config...")
	err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}
	if migrate {
		log.Println("Init old DB...")
		db, err := old_datastore.InitDB()
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()
		log.Println("Init new DB...")
		_, err = datastore.InitDB()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Migrating...")
		Migrate()
		log.Println("Waiting for migration to finish...")
		time.Sleep(15 * time.Second)
	} else {
		log.Println("Init DB...")
		_, err = datastore.InitDB()
		if err != nil {
			log.Fatal(err)
		}
		worker.InitWorker()
		for {
			time.Sleep(60 * time.Second)
		}
	}
}
