package subscriberep

import (
	"claime-verifier/lib/functions/lib/common/log"
	"claime-verifier/lib/functions/lib/subscribe"
	"context"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

type (
	// Repository repository
	Repository struct {
		ddb dynamo.DB
	}
)

// New new client
func New() Repository {
	return Repository{
		ddb: *dynamo.New(session.New()),
	}
}

// Subscribe register subscribe
func (r Repository) Subscribe(ctx context.Context, sub subscribe.Subscription) error {
	if err := r.ddb.Table(tableName()).Put(&sub).RunWithContext(ctx); err != nil {
		log.Error("put item failed", err)
		return err
	}
	return nil
}

// Subscribing All Subscriptions
func (r Repository) Subscribing(ctx context.Context) ([]subscribe.Subscription, error) {
	return r.scan(ctx, r.sc())
}

// SubscribingFrom subscribe from
func (r Repository) SubscribingFrom(ctx context.Context, t time.Time) ([]subscribe.Subscription, error) {
	return r.scan(ctx, r.sc().Filter("Timestamp > ?", t.Unix()))
}

func (r Repository) sc() *dynamo.Scan {
	return r.ddb.Table(tableName()).Scan()
}

func (r Repository) scan(ctx context.Context, op *dynamo.Scan) ([]subscribe.Subscription, error) {
	subs := []subscribe.Subscription{}
	err := op.AllWithContext(ctx, &subs)
	if err != nil {
		log.Error("scan failed", err)
	}
	return subs, err
}

// UnSubscribe delete subscribe
func (r Repository) UnSubscribe(ctx context.Context, in subscribe.Input) error {
	return nil
}

func tableName() string {
	return "claime-verifier-main"
}
