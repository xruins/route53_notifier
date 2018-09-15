package address

import "github.com/aws/aws-sdk-go/service/route53"

type IPAddr struct {
	IPv4Addr, IPv6Addr string
}

func (ipaddr *IPAddr) ToResourceRecordSet(fqdn string, ttl int64) []*route53.ResourceRecordSet {
	var recordSets []*route53.ResourceRecordSet
	if ipaddr.IPv4Addr != "" {
		recordSets = append(recordSets, generateResourceRecordSet(fqdn, "A", ipaddr.IPv4Addr, ttl))
	}
	if ipaddr.IPv6Addr != "" {
		recordSets = append(recordSets, generateResourceRecordSet(fqdn, "AAAA", ipaddr.IPv6Addr, ttl))
	}
	return recordSets
}

func generateResourceRecordSet(fqdn, recordType, value string, ttl int64) *route53.ResourceRecordSet {
	return &route53.ResourceRecordSet{
		Name: &fqdn,
		Type: &recordType,
		TTL: &ttl,
		ResourceRecords: []*route53.ResourceRecord{
			&route53.ResourceRecord{Value: &value},
		},
	}
}
