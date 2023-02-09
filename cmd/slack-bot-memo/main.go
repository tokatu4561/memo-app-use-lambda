package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/tokatu4561/memo-app-use/pkg/application/di"
)

// TODO: env管理する
const AWS_REGION = "ap-northeast-1"
const DYNAMO_ENDPOINT = "http://dynamodb:8000"

var secret = os.Getenv("SLACK_SIGNING_SECRET")
var oAuthToken = os.Getenv("SLACK_OAUTH_TOKEN")

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	api := slack.New(os.Getenv("SLACK_OAUTH_TOKEN"))

	err := Verify(ConvertHeaders(request.Headers), []byte(request.Body))
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "slack conection error", StatusCode: 400}, err
	}
	
	body := request.Body
	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "slack conection error", StatusCode: 500}, err
	}

	switch eventsAPIEvent.Type {
		case slackevents.URLVerification:
			res, err := HandleURLVerification(body)
			if err != nil {
				return events.APIGatewayProxyResponse{Body: "slack conection error", StatusCode: 500}, err
			}
			return events.APIGatewayProxyResponse{
				Body:       res.Challenge,
				StatusCode: 200,
			}, nil
		case slackevents.CallbackEvent:
			innerEvent := eventsAPIEvent.InnerEvent
			switch event := innerEvent.Data.(type) {
			case *slackevents.AppMentionEvent:
				msg := strings.Split(event.Text, " ")
				cmd := msg[1]
				ctl := di.NewSlackMemoController()

				switch cmd {
				case "ping": 
					// MsgOptionText() の第二引数に true を設定すると特殊文字をエスケープする
					_, _, err := api.PostMessage(event.Channel, slack.MsgOptionText("pong", false))
					if err != nil {
						return events.APIGatewayProxyResponse{Body: "bad request", StatusCode: 400}, err
					}
				case "memo":
					memo, err := ctl.CreateMemo(msg[2])
					if err != nil {
						return events.APIGatewayProxyResponse{Body: "bad request", StatusCode: 400}, err
					}

					responseMsg := fmt.Sprintf("%sを追加しました!", memo.Title)
					_, _, err = api.PostMessage(event.Channel, slack.MsgOptionText(responseMsg, false))
					if err != nil {
						return events.APIGatewayProxyResponse{Body: "bad request", StatusCode: 400}, err
					}
				case "list":
					_, err := ctl.GetMemos()
					if err != nil {
						return events.APIGatewayProxyResponse{Body: "bad request", StatusCode: 400}, err
					}

					attachment := slack.Attachment{
						Pretext:    "pretext",
						Fallback:   "We don't currently support your client",
						CallbackID: "accept_or_reject",
						Color:      "#3AA3E3",
						Actions: []slack.AttachmentAction{
							{
								Name:  "accept",
								Text:  "Accept",
								Type:  "button",
								Value: "accept",
							},
							{
								Name:  "reject",
								Text:  "Reject",
								Type:  "button",
								Value: "reject",
								Style: "danger",
							},
						},
					}

					message := slack.MsgOptionAttachments(attachment)
					_, _, err = api.PostMessage(event.Channel, slack.MsgOptionText("", false), message)
					if err != nil {
						return events.APIGatewayProxyResponse{Body: "bad request", StatusCode: 400}, err
					}
				}
			}
	}
	
	return events.APIGatewayProxyResponse{
		Body:       "",
		StatusCode: 200,
	}, nil
}

func ConvertHeaders(headers map[string]string) http.Header {
    h := http.Header{}
    for key, value := range headers {
        h.Set(key, value)
    }
    return h
}

func HandleURLVerification(body string) (*slackevents.ChallengeResponse ,error) {
	var res *slackevents.ChallengeResponse
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		return nil, err
	}

	return res, nil
}

// slack からのリクエストかを検証 外部からのリクエストを受け付けないように
// ヘッダー、body、Signing Secretで検証
func Verify (header http.Header, body []byte) error {
	verifier, err := slack.NewSecretsVerifier(header, os.Getenv("SLACK_SIGNING_SECRET"))
	if err != nil {
		return err
	}

	verifier.Write(body)
	if err := verifier.Ensure(); err != nil {
		return err
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
