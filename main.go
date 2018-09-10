package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/xruins/route53_notifier/address"
	"github.com/xruins/route53_notifier/notifier"
)

func main() {
	var credential, host, iface, ipv4addr, ipv6addr string
	var daemon bool
	var interval time.Duration

	flag.StringVar(&credential, "c", "", "Path to Credential.")
	flag.StringVar(&host, "h", "", "Host for the key of A/AAAA records.")
	flag.StringVar(&iface, "i", "", "Network interface name to get IP addresses.")
	flag.StringVar(&ipv4addr, "4", "", "IPv4 address to notify. used for override auto detected one.")
	flag.StringVar(&ipv6addr, "6", "", "IPv6 address to notify. used for override auto detected one.")
	flag.BoolVar(&daemon, "d", false, "If true, this program persists and continue to notify IP addresses.")
	flag.DurationVar(&interval, "i", 600, "Seconds of notification interval. Works only for daemon mode.")

	if ipv4addr == "" || ipv6addr == "" {
		addr, err := address.GetIpAddrs(iface)
		if err != nil {
			log.Fatalf("an error occured when get IP addresses from interface: %v", err)
		}
		if ipv4addr == "" {
			ipv4addr = addr.Ipv4addr
		}
		if ipv6addr == "" {
			ipv6addr = addr.Ipv6addr
		}
	}

	sess := session.Must(session.NewSession())

	ctx := context.Background()
	ntf := &notifier.Notifier{
		Sess:    sess,
		Context: ctx,
	}
	err := ntf.Notify(ipv4addr, ipv6addr)
	if err != nil {
		log.Fatalf("an error occured when notify route53: %v", err)
	}
	os.Exit(0)
}
