package network

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/alexgear/suplat/common"
	"github.com/alexgear/suplat/config"
)

var ief string

func getInterface(addr string) (string, error) {
	for tag, network := range config.C.Networks {
		for _, cidr := range network.CIDRs {
			_, network, err := net.ParseCIDR(cidr)
			if err != nil {
				return "", fmt.Errorf("Failed to parse CIDR: %s", err.Error())
			}
			if network.Contains(net.ParseIP(addr)) {
				return tag, nil
			}
		}
	}
	return "", fmt.Errorf("No networks specified in config.")
}

func Ping() (common.Point, error) {
	start := time.Now().UTC()
	point := common.Point{Time: start}
	d := net.Dialer{Timeout: 5 * time.Second}
	conn, err := d.Dial("tcp", "8.8.8.8:53")
	if err != nil {
		point.Latency = 5 * time.Second
		point.Error = 1
		point.Interface = ief
		return point, nil
	}
	point.Latency = time.Since(start)
	point.Error = 0
	ief, err = getInterface(strings.Split(conn.LocalAddr().String(), ":")[0])
	if err != nil {
		return point, fmt.Errorf("Failed to get interface: %s", err.Error())
	}
	point.Interface = ief
	conn.Close()
	return point, nil
}
