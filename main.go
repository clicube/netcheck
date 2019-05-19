package main

import (
	"fmt"
	"log"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	ping "github.com/sparrc/go-ping"
)

func main() {

	client, err := statsd.New(
		"127.0.0.1:8125",
		statsd.WithNamespace("home.net.google."),
	)
	if err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(time.Second * 15)
	defer ticker.Stop()

	for {
		doCheck(client)
		select {
		case <-ticker.C:
		}
	}

}

func doCheck(client *statsd.Client) {

	host := "www.google.com"
	pinger, err := ping.NewPinger(host)
	if err != nil {
		log.Fatal(err)
	}
	pinger.Count = 1
	pinger.Timeout = time.Second * 2
	pinger.Interval = time.Millisecond * 100
	pinger.SetPrivileged(true)

	pinger.OnFinish = func(stats *ping.Statistics) {

		var up float64
		if stats.PacketLoss == 0 {
			up = 1
		} else {
			up = 0
		}
		client.Gauge("up", up, []string{}, 1)
		client.Gauge("time", float64(stats.AvgRtt/time.Millisecond), []string{}, 1)
	}

	fmt.Printf("Ping: %s (%s)", pinger.Addr(), pinger.IPAddr())
	pinger.Run()
}
