package main

import (
	"claime-verifier/lib/functions/lib"
	"claime-verifier/lib/functions/lib/common/log"
	"claime-verifier/lib/functions/lib/infrastructure/eth"
	"claime-verifier/lib/functions/lib/infrastructure/ssm"
	"claime-verifier/lib/functions/lib/subscribe"
	repository "claime-verifier/lib/functions/lib/subscribe/persistence"
	"context"
	"os"
)

func main() {
	ctx := context.Background()
	app, err := newApp(ctx)
	if err != nil {
		os.Exit(1)
	}
	if err := app.StartSubscription(ctx); err != nil {
		log.Error("an error occured", err)
	}

	os.Exit(1)
}

func newApp(ctx context.Context) (subscribe.Subscriber, error) {
	parameterstore := ssm.New()
	sender, err := lib.Slackcli(ctx, parameterstore)
	if err != nil {
		log.Error("initialize sender failed", err)
		return subscribe.Subscriber{}, err
	}
	return subscribe.NewSubscriber(ctx, parameterstore, sender, repository.New(), eth.Subscriber{})
}
