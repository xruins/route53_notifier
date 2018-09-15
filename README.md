# Description
A golang implementation to notify your IP address to AWS Route 53 as A/AAAA redords.

# Usage

```
Usage of ./route53_notifier:
  -fqdn string
    	FQDN for the key of A/AAAA records. (required)
  -hosted_zone_id string
    	FQDN for the key of A/AAAA records. (required)
  -iface string
    	Network interface name to get IPv4 addresses.
  -ifacev6 string
    	Network interface name to get IPv6 addresses. If blank, use the one of v4.
  -ipv4 string
    	IPv4 address to notify. used for override auto detected one.
  -ipv6 string
    	IPv6 address to notify. used for override auto detected one.
```

route53_notifier requires at least one of `-iface`, `-ifacev6`, `-ipv4` and `-ipv6`.

# Build

`GO111MODULE=on go build`
