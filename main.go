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
		statsd.WithNamespace("home.net."),
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
	pinger, err := ping.NewPinger("www.google.com")
	if err != nil {
		log.Fatal(err)
	}
	pinger.Count = 5
	pinger.Timeout = time.Second * 2
	pinger.Interval = time.Millisecond * 100
	pinger.SetPrivileged(true)

	// pinger.OnRecv = func(pkt *ping.Packet) {
	// 	fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
	// 		pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
	// }

	pinger.OnFinish = func(stats *ping.Statistics) {
		fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
		fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)

		client.Gauge("up", float64(1 - stats.PacketLoss), []string{}, 1)
		client.Gauge("avg", float64(stats.AvgRtt/time.Millisecond), []string{}, 1)
		client.Gauge("min", float64(stats.MinRtt/time.Millisecond), []string{}, 1)
		client.Gauge("max", float64(stats.MaxRtt/time.Millisecond), []string{}, 1)
		client.Gauge("stddev", float64(stats.StdDevRtt/time.Millisecond), []string{}, 1)
	}

	fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())
	pinger.Run()
}
