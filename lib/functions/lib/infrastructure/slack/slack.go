package slackclient

import (
	"bytes"
	"claime-verifier/lib/functions/lib/common/log"
	"claime-verifier/lib/functions/lib/subscribe"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/kyokomi/goslash/goslash"
	"github.com/slack-go/slack"
)

type (
	// Client client
	Client struct {
		svc           *slack.Client
		signingSecret string
	}
	// CommandResponse response
	CommandResponse struct {
		goslash.SlashCommandMessage
	}
)

// Parse parse response
func (r CommandResponse) Parse() ([]byte, error) {
	return json.Marshal(&r)
}

// New new client
func New(token, signingSecret string) Client {
	return Client{
		svc:           slack.New(token, slack.OptionDebug(true)),
		signingSecret: signingSecret,
	}
}

// Send send event notification
func (c Client) Send(ctx context.Context, event subscribe.Event, subscription subscribe.Subscription) error {
	attachment := []slack.Attachment{
		{
			Fields: attatchmentFields(event),
		},
	}
	_, _, err := c.svc.PostMessageContext(ctx, subscription.ChannelID, slack.MsgOptionAttachments(attachment[:]...))
	return err
}

func attatchmentFields(event subscribe.Event) []slack.AttachmentField {
	res := []slack.AttachmentField{}
	res = append(res, slack.AttachmentField{
		Title: "ContractAddress",
		Value: event.ContractAddress,
	})
	res = append(res, slack.AttachmentField{
		Title: "Event",
		Value: event.EventName,
	})
	res = append(res, slack.AttachmentField{
		Title: "BlockNumber",
		Value: strconv.FormatUint(event.BlockNumber, 10),
	})
	res = append(res, slack.AttachmentField{
		Title: "Hash",
		Value: event.TxHash,
	})
	res = append(res, slack.AttachmentField{
		Title: "Link",
		Value: event.Link,
	})
	for k, v := range event.Parameters {
		res = append(res, slack.AttachmentField{
			Title: k,
			Value: fmt.Sprintf("%v", v),
		})
	}
	return res
}

func parseRequest(body string) (slack.SlashCommand, error) {
	req, _ := http.NewRequest("POST", "", bytes.NewBufferString(body))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	s, err := slack.SlashCommandParse(req)
	if err != nil {
		log.Error("unexpected error occured", err)
		return slack.SlashCommand{}, err
	}
	return s, nil
}

// Verify verify request
func (c Client) Verify(request events.APIGatewayProxyRequest) error {
	body := request.Body
	headers := convertHeaders(request.Headers)
	sv, e := slack.NewSecretsVerifier(headers, c.signingSecret)
	if e != nil {
		return e
	}
	sv.Write([]byte(body))
	return sv.Ensure()
}

// ErrorOutput error to output
func (c Client) ErrorOutput(err error) subscribe.Output {
	out := CommandResponse{c.newCommandMessage()}
	out.Attachments = []goslash.Attachment{{
		Title:   "Invalid Command! :crying_cat_face: ",
		Pretext: "usage example: `/claime-verifier subscribe rinkeby 0x4beb7299221807Cd47C2fa118c597C51Cc2fEC99 Refund https://raw.githubusercontent.com/bridges-inc/kaleido-core/develop/deployments/rinkeby/AdManager.json`",
		Text:    err.Error(),
	}}
	out.ResponseType = "ephemeral"
	return out
}

// ToOutput to output
func (c Client) ToOutput(in subscribe.Input) subscribe.Output {
	out := CommandResponse{c.newCommandMessage()}
	out.Attachments = []goslash.Attachment{{
		Title: "Subscription started.",
		Fields: []goslash.AttachmentField{
			{
				Title: "Address",
				Value: in.Address,
			},
			{
				Title: "Network",
				Value: in.Network,
			},
			{
				Title: "Event",
				Value: in.Event,
			}, {
				Title: "Abi",
				Value: in.AbiURL,
			},
		},
	}}
	out.ResponseType = "in_channel"
	return out
}

func (c Client) newCommandMessage() goslash.SlashCommandMessage {
	return goslash.SlashCommandMessage{
		Attachments: []goslash.Attachment{},
	}
}

func convertHeaders(headers map[string]string) http.Header {
	h := http.Header{}
	for key, value := range headers {
		h.Set(key, value)
	}
	return h
}

// ToInputFromBody to subscribe input from request body.
func (c Client) ToInputFromBody(body string) (subscribe.Input, error) {
	cm, err := parseRequest(body)
	if err != nil {
		return subscribe.Input{}, err
	}
	options := strings.Split(cm.Text, " ")
	in := &subscribe.Input{}
	parseFields(in, options)
	in.ChannelID = cm.ChannelID
	in.WebhookURL = cm.ResponseURL
	return *in, nil
}

func parseFields(input *subscribe.Input, options []string) {
	if len(options) == 0 {
		return
	}
	input.Type = options[0]
	if len(options) == 1 {
		return
	}
	input.Network = options[1]
	if len(options) == 2 {
		return
	}
	input.Address = options[2]
	if len(options) == 3 {
		return
	}
	input.Event = options[3]
	if len(options) == 4 {
		return
	}
	input.AbiURL = options[4]
}
