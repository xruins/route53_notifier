package notifier

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/session"
)

type Notifier struct {
	Sess    *session.Session
	Context context.Context
}

func (n *Notifier) Notify(ipv4addr, ipv6addr string) error {
	// [TODO] implementation
	return nil
}
