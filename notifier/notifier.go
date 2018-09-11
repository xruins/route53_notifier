package notifier

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/xruins/route53_notifier/address"
)

type Notifier struct {
	Session      *session.Session
	Context      context.Context
	FQDN         string
	HostedZoneId string
	IPAddr       *address.IPAddr
}

func (n *Notifier) Notify() error {
	r53 := route53.New(n.Session)
	recordSets := n.IPAddr.ToResourceRecordSet(n.fqdn)

	var changeBatch *route53.ChangeBatch
	var changes *route53.Change
	for _, recordSet := range recordSets {
		change := initChange(recordSet)
		changeBatch = append(changeBatch, change)
	}

	input := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch:  changeBatch,
		HostedZoneId: n.HostedZoneId,
	}
	output, err := r53.ChangeResourceRecordSets(input)
	if err != nil {
		return fmt.Errorf("failed to update route53 resource record sets: %s", err)
	}
	log.Printf(
		"succeed to update resource record sets of HostedZoneID %s. detail: %s",
		n.HostedZoneId,
		output.GoString(),
	)

	return nil
}

var route53ChangeAction = "CREATE"

func generateChangeBatch(recordSets []*route53.ResourceRecordSet) *route53.ChangeBatch {
	var changes []*route53.Change
	for _, recordSet := range recordSets {
		changes = append(
			changes,
			&route53.Change{
				Action:            &route53ChangeAction,
				ResourceRecordSet: recordSet,
			},
		)
	}
	return &route53.ChangeBatch{Changes: changes}
}
