package main

import (
	"flag"
	"log"
	"os"

	"github.com/xruins/route53_notifier/address"
	"github.com/xruins/route53_notifier/notifier"
)

func main() {
	var fqdn, hostedZoneId, iface, ifacev6, ipv4addr, ipv6addr string

	flag.StringVar(&fqdn, "fqdn", "", "FQDN for the key of A/AAAA records.")
	flag.StringVar(&hostedZoneId, "hosted_zone_id", "", "FQDN for the key of A/AAAA records.")
	flag.StringVar(&iface, "iface", "", "Network interface name to get IPv4 addresses.")
	flag.StringVar(&ifacev6, "ifacev6", "", "Network interface name to get IPv6 addresses. If blank, use the one of v4.")
	flag.StringVar(&ipv4addr, "ipv4", "", "IPv4 address to notify. used for override auto detected one.")
	flag.StringVar(&ipv6addr, "ipv6", "", "IPv6 address to notify. used for override auto detected one.")
	flag.Parse()


	if ipv4addr == "" {
		ipaddrs, err := address.GetIPAddr(iface)
		if err != nil {
			log.Fatalf("failed to get ipv4 addresses: %s\n", err)
		}
		ipv4addr = ipaddrs.IPv4Addr
	}

	if ipv6addr == "" {
		ipaddrs, err := address.GetIPAddr(ifacev6)
		if err != nil {
			log.Fatalf("failed to get ipv6 addresses: %s\n", err)
		}
		ipv6addr = ipaddrs.IPv6Addr
	}

	if ipv4addr == "" && ipv6addr == "" {
		log.Fatalln("couldn't get neither ipv4 address nor ipv6 one")
	}

	ntf := &notifier.Notifier{
		FQDN:         fqdn,
		HostedZoneId: hostedZoneId,
		IPAddr:       &address.IPAddr{IPv4Addr: ipv4addr, IPv6Addr: ipv6addr},
	}
	ntfMsg, ntfErr := ntf.Notify()
	if ntfErr != nil {
		log.Fatalf("an error occured when notify route53: %s\n", ntfErr)
	}
	log.Printf("successfully updated. ipv4: %s, ipv6: %s, message: %s\n", ipv4addr, ipv6addr, ntfMsg)
	os.Exit(0)
}
