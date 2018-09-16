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
	var ttl int64

	flag.StringVar(&fqdn, "fqdn", "", "FQDN for the key of A/AAAA records. (required)")
	flag.StringVar(&hostedZoneId, "hosted_zone_id", "", "FQDN for the key of A/AAAA records. (required)")
	flag.StringVar(&iface, "iface", "", "Network interface name to get IPv4 addresses.")
	flag.StringVar(&ifacev6, "ifacev6", "", "Network interface name to get IPv6 addresses. If blank, use the one of v4.")
	flag.StringVar(&ipv4addr, "ipv4", "", "IPv4 address to notify. used for override auto detected one.")
	flag.StringVar(&ipv6addr, "ipv6", "", "IPv6 address to notify. used for override auto detected one.")
	flag.Int64Var(&ttl, "ttl", 3600, "seconds to TTL of DNS record.")
	flag.Parse()

	// check commandline args
	if fqdn == "" || hostedZoneId == "" {
		log.Fatalln("both of -fqdn and -hosted_zone_id are required.")
	}
	if ipv4addr == "" && ipv6addr == "" && iface == "" && ifacev6 == "" {
		log.Fatalln("specify at least one of -iface, -ifacev6, -ipv4addr and -ipv6addr.")
	}

	// get ipv4/v6 addresses from interfaces
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

	// notify Route53
	ntf := &notifier.Notifier{
		FQDN:         fqdn,
		HostedZoneId: hostedZoneId,
		IPAddr:       &address.IPAddr{IPv4Addr: ipv4addr, IPv6Addr: ipv6addr},
		TTL:          ttl,
	}
	ntfMsg, ntfErr := ntf.Notify()
	if ntfErr != nil {
		log.Fatalf("an error occured when notify route53: %s\n", ntfErr)
	}
	log.Printf("successfully updated. ipv4: %s, ipv6: %s, message: %s\n", ipv4addr, ipv6addr, ntfMsg)
	os.Exit(0)
}
