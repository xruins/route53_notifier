package notifier

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/xruins/route53_notifier/address"
)

type Notifier struct {
	FQDN         string
	HostedZoneId string
	IPAddr       *address.IPAddr
	TTL int64
}

var sess *session.Session

func init() {
	sess = session.Must(session.NewSession())
}

func (n *Notifier) Notify() (string, error) {
	r := route53.New(sess)
	recordSets := n.IPAddr.ToResourceRecordSet(n.FQDN, n.TTL)
	changeBatch := generateChangeBatch(recordSets)

	input := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch:  changeBatch,
		HostedZoneId: &n.HostedZoneId,
	}
	output, err := r.ChangeResourceRecordSets(input)
	if err != nil {
		return "", fmt.Errorf("failed to update route53 resource record sets: %s", err)
	}

	msg := output.GoString()
	return msg, nil
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
