package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// TODO: env管理する
const AWS_REGION = "ap-northeast-1"
const DYNAMO_ENDPOINT = "http://dynamodb:8000"

var secret = os.Getenv("SLACK_SIGNING_SECRET")
var oAuthToken = os.Getenv("SLACK_OAUTH_TOKEN")

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	api := slack.New(os.Getenv("SLACK_OAUTH_TOKEN"))

	// slack からのリクエストかを検証 外部からのリクエストを受け付けないように
	// ヘッダー、body、Signing Secretで検証
	verifier, err := slack.NewSecretsVerifier(ConvertHeaders(request.Headers), os.Getenv("SLACK_SIGNING_SECRET"))
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "slack conection error", StatusCode: 500}, err
	}
	verifier.Write([]byte(request.Body))
	if err := verifier.Ensure(); err != nil {
		return events.APIGatewayProxyResponse{Body: "slack conection error", StatusCode: 400}, err
	}
	
	body := request.Body
	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "slack conection error", StatusCode: 500}, err
	}
	
	switch eventsAPIEvent.Type {
		case slackevents.URLVerification:
			var res *slackevents.ChallengeResponse
			if err := json.Unmarshal([]byte(body), &res); err != nil {
				log.Println(err)
				if err != nil {
					return events.APIGatewayProxyResponse{Body: "slack conection error", StatusCode: 500}, err
				}
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
				log.Println(fmt.Sprintf("%s\n", msg))
				cmd := msg[1]
				switch cmd {
				case "ping": 
					// MsgOptionText() の第二引数に true を設定すると特殊文字をエスケープする
					_, _, err := api.PostMessage(event.Channel, slack.MsgOptionText("pong", false))
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

func main() {
	lambda.Start(handler)
}
