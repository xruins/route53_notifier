package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/xruins/route53_notifier/address"
	"github.com/xruins/route53_notifier/notifier"
)

func main() {
	var fqdn, iface, ifacev6, ipv4addr, ipv6addr string

	flag.StringVar(&fqdn, "fqdn", "", "FQDN for the key of A/AAAA records.")
	flag.StringVar(&iface, "iface", "", "Network interface name to get IPv4 addresses.")
	flag.StringVar(&ifacev6, "ifacev6", "", "Network interface name to get IPv6 addresses. If blank, use the one of v4.")
	flag.StringVar(&ipv4addr, "ipv4", "", "IPv4 address to notify. used for override auto detected one.")
	flag.StringVar(&ipv6addr, "ipv6", "", "IPv6 address to notify. used for override auto detected one.")

	if ipv4addr == "" || ipv6addr == "" {
		ipaddrs, err := address.GetIPAddr(iface, ifacev6)
		if err != nil {
			log.Fatalf("failed to get ip addresses: %s\n", err)
		}

		if ipv4addr == "" {
			ipv4addr = ipaddrs.Ipv4addr
		}
		if ipv6addr == "" {
			ipv6addr = ipaddrs.Ipv6addr
		}
	}

	sess := session.Must(session.NewSession())

	ctx := context.Background()
	ntf := &notifier.Notifier{
		Sess:    sess,
		Context: ctx,
	}
	ntfErr := ntf.Notify(ipv4addr, ipv6addr)
	if ntfErr != nil {
		log.Fatalf("an error occured when notify route53: %s\n", ntfErr)
	}
	log.Printf("successfully updated. ipv4: %s, ipv6: %s\n", ipv4addr, ipv6addr)
	os.Exit(0)
}
