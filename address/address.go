package address

import (
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

func GetIPAddr(iface_name string) (*IPAddr, error) {
	var ipv4, ipv6 string
	iface, err := getIPAddrFromInterface(iface_name)
	if err != nil {
		return nil, fmt.Errorf("failed to get IPv4 address: %s", err)
	}
	if iface.IPv4Addr != "" {
		ipv4 = iface.IPv4Addr
	}
	if iface.IPv6Addr != "" {
		ipv6 = iface.IPv6Addr
	}

	return &IPAddr{IPv4Addr: ipv4, IPv6Addr: ipv6}, nil
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
			if isIPv4Addr(&ip) {
				ret.IPv4Addr = ip_str
			} else if isIPv6Addr(&ip) {
				ret.IPv6Addr = ip_str
			} else {
				return nil, fmt.Errorf("ambigious ip address detected: %s", ip_str)
			}
		}
	}

	return ret, nil
}

func isIPv4Addr(ip *net.IP) bool {
	if len(*ip) == net.IPv4len {
		return true
	} else {
		return false
	}
}

func isIPv6Addr(ip *net.IP) bool {
	if len(*ip) == net.IPv6len {
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
