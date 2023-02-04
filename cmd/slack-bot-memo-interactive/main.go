package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// api := slack.New(os.Getenv("SLACK_OAUTH_TOKEN"))

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
		default:
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
