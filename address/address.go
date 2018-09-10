package address

type Addresses struct {
	Ipv4addr, Ipv6addr string
}

func GetIpAddrs(iface string) (*Addresses, error) {
	// [TODO] implementation
	return &Addresses{
		Ipv4addr: "127.0.0.1",
		Ipv6addr: ":::",
	}, nil
}
