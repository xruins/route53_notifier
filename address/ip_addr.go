package address

import "github.com/aws/aws-sdk-go/service/route53"

type IPAddr struct {
	IPv4Addr, IPv6Addr string
}

func (ipaddr *IPAddr) ToResourceRecordSet(fqdn string) []*route53.ResourceRecordSet {
	var recordSets []*route53.ResourceRecordSet
	if ipaddr.IPv4Addr != "" {
		recordSets = append(recordSets, generateResourceRecordSet(fqdn, "A", ipaddr.IPv4Addr))
	}

	if ipaddr.IPv6Addr != "" {
		recordSets = append(recordSets, generateResourceRecordSet(fqdn, "AAAA", ipaddr.IPv6Addr))
	}
	return recordSets
}

func generateResourceRecordSet(fqdn, recordType, value string) *route53.ResourceRecordSet {
	return &route53.ResourceRecordSet{
		Name: &fqdn,
		Type: &recordType,
		ResourceRecords: []*route53.ResourceRecord{
			&route53.ResourceRecord{Value: &value},
		},
	}
}
