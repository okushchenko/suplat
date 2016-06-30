package datastore

import (
	"fmt"
	"time"

	"github.com/alexgear/suplat/common"
	"github.com/alexgear/suplat/config"
	"github.com/cenkalti/backoff"
	client "github.com/influxdb/influxdb/client/v2"
)

var err error
var clnt client.Client

const MyDB = "suplat"

func queryDB(cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: MyDB,
	}
	if response, err := clnt.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	}
	return res, nil
}

func InitDB() (client.Client, error) {
	clnt, err = client.NewHTTPClient(client.HTTPConfig{
		Addr: config.C.InfluxDBUrl,
	})
	if err != nil {
		return clnt, fmt.Errorf("Failed to connect to InfluxDB: %s", err)
	}
	_, err := queryDB(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", MyDB))
	if err != nil {
		return clnt, fmt.Errorf("Failed to create database: %s", err)
	}
	return clnt, nil
}

var cPoints = make(chan common.Point, 1000)

func Write(r common.Point) {
	cPoints <- r
}

func Flush() error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database: MyDB,
	})
	if err != nil {
		return err
	}
	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = 5 * time.Minute
	for {
		bp, err = client.NewBatchPoints(client.BatchPointsConfig{
			Database: MyDB,
		})
		if err != nil {
			return err
		}
		for i := 0; i < 10; i++ {
			point := <-cPoints
			tags := map[string]string{"interface": point.Interface}
			fields := map[string]interface{}{
				"error":   point.Error,
				"latency": point.Latency.Seconds(),
			}
			pt, err := client.NewPoint("network", tags, fields, point.Time)
			if err != nil {
				return err
			}
			bp.AddPoint(pt)
		}
		err = backoff.Retry(func() error {
			return clnt.Write(bp)
		}, bo)
		if err != nil {
			bo.Reset()
			return err
		}
	}
	return nil
}
