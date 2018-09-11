package address

import (
	"errors"
	"fmt"
	"net"
)

var privateIPAddrs = [...]string{
	"127.0.0.0/8",    // IPv4 loopback
	"10.0.0.0/8",     // RFC1918
	"172.16.0.0/12",  // RFC1918
	"192.168.0.0/16", // RFC1918
	"::1/128",        // IPv6 loopback
	"fe80::/10",      // IPv6 link-local
}

var privateIPBlocks []*net.IPNet

func init() {
	for _, cidr := range privateIPAddrs {
		_, block, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(err)
		}
		privateIPBlocks = append(privateIPBlocks, block)
	}
}

type IPAddr struct {
	Ipv4Addr, Ipv6Addr string
}

func GetIPAddr(iface, ifacev6 string) (*IPAddr, error) {
	var ipv4, ipv6 string
	ifv4, err := getIPAddrFromInterface(iface)
	if err != nil {
		return nil, err
	}
	if ifv4.Ipv4Addr != "" {
		ipv4 = ifv4.Ipv4Addr
	}
	// if ifacev6 was specified, get v6 address from it
	if ifacev6 != "" {
		ifv6, err := getIPAddrFromInterface(iface)
		if err != nil {
			return nil, err
		}
		if ifv6.Ipv6Addr != "" {
			ipv6 = ifv6.Ipv6Addr
		}
	} else { // if ifacev6 was not specified, use the v6 address of ifacev4
		if ifv4.Ipv6Addr != "" {
			ipv4 = ifv4.Ipv6Addr
		}
	}

	if ipv4 == "" && ipv6 == "" {
		return nil, errors.New("couldn't get neither ipv4 addr not ipv6 one")
	}

	return &IPAddr{Ipv4Addr: ipv4, Ipv6Addr: ipv6}, nil
}

func getIPAddrFromInterface(iface_name string) (*IPAddr, error) {
	iface, err := net.InterfaceByName(iface_name)
	if err != nil {
		return nil, err
	}
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}

	var ret *IPAddr
	for _, addr := range addrs {
		addr_str := addr.String()
		ip, _, err := net.ParseCIDR(addr_str)
		if err != nil {
			return nil, fmt.Errorf("ambigious ip address detected: %s", addr_str)
		}

		if isGlobalIP(&ip) {
			ip_str := ip.String()
			if isIpv4Addr(ip_str) {
				ret.Ipv4Addr = ip_str
			} else if isIpv6Addr(ip_str) {
				ret.Ipv6Addr = ip_str
			} else {
				return nil, fmt.Errorf("ambigious ip address detected: %s", ip_str)
			}
		}
	}

	return ret, nil
}

// func GetIPAddrs(iface string) (*Addresses, error) {

// enp5s0, err := net.InterfaceByName("enp5s0")
// addrs, err := enp5s0.Addrs()
// var ips []net.IP
// var ipaddrs []string
// for _, addr := range addrs {
//     ipaddrs = append(ipaddrs, addr.Network())
// }
// for _, addr := range addrs {
//     ipaddrs = append(ipaddrs, addr.String())
// }
// ip, ipnet, error := net.ParseCIDR("172.16.0.2/16")
// ip.String()
// ip.IsGlobalUnicast()
// ip.IsGlobalUnicast()
// net.IP{8, 8, 8, 8}.IsGlobalUnicast()
// net.IP{172, 16, 0, 1}.IsGlobalUnicast()

// [TODO] implementation
//	return &Addresses{
//		Ipv4Addr: "127.0.0.1",
//		Ipv6Addr: ":::",
//	}, nil

func isIPv4Addr(addr string) bool {
	if len(addr) == net.IPv4len {
		return true
	} else {
		return false
	}
}

func isIpv6Addr(addr string) bool {
	if len(addr) == net.IPv6len {
		return true
	} else {
		return false
	}
}

func isGlobalIP(ip *net.IP) bool {
	for _, block := range privateIPBlocks {
		if block.Contains(*ip) {
			return false
		}
	}
	return true
}
