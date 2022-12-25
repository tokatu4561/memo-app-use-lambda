package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// TODO: env管理する
// const AWS_REGION = "ap-northeast-1"
// const DYNAMO_ENDPOINT = "http://dynamodb:8000"


var secret = os.Getenv("SLACK_SIGNING_SECRET")
var oAuthToken = os.Getenv("SLACK_OAUTH_TOKEN")

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	api := slack.New("xoxb-4588063634176-4577008852193-IXeR4FDhiNU4GK9tIvqBiCrc")
	
	body := request.Body
	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "slack conection error", StatusCode: 500}, err
	}

	var res *slackevents.ChallengeResponse

	switch eventsAPIEvent.Type {
		case slackevents.URLVerification:
			if err := json.Unmarshal([]byte(body), &res); err != nil {
				log.Println(err)
				if err != nil {
					return events.APIGatewayProxyResponse{Body: "slack conection error", StatusCode: 500}, err
				}
			}
		case slackevents.CallbackEvent:
			innerEvent := eventsAPIEvent.InnerEvent
			switch event := innerEvent.Data.(type) {
			case *slackevents.AppMentionEvent:
				msg := strings.Split(event.Text, " ")
				// if len(msg) < 2 {
                //     return
                // }
				log.Println(fmt.Sprintf("%s\n", msg))
				cmd := msg[1]
				switch cmd {
				case "ping": 
					log.Println("通過1")
					_, _, err := api.PostMessage(event.Channel, slack.MsgOptionText("pong", false))
					if err != nil {
						return events.APIGatewayProxyResponse{Body: "bad request", StatusCode: 400}, err
					}
				}
			}
	}
	

	return events.APIGatewayProxyResponse{
		Body:       res.Challenge,
		StatusCode: 200,
	}, nil
}


// func setUpLineClient() (*Line, error) {
// 	line := &Line{
// 		ChannelSecret: os.Getenv("LINE_BOT_CHANNEL_SECRET"),
// 		ChannelToken:  os.Getenv("LINE_BOT_CHANNEL_TOKEN"),
// 	}

// 	bot, err := linebot.New(
// 		line.ChannelSecret,
// 		line.ChannelToken,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	line.Client = bot

// 	return line, nil
// }

// func setUpDB() (*dynamo.DB, error) {
// 	sess, err := session.NewSession(&aws.Config{
// 		Region:      aws.String(AWS_REGION),
// 		Endpoint:    aws.String(DYNAMO_ENDPOINT),
// 		Credentials: credentials.NewStaticCredentials("dummy", "dummy", "dummy"),
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	db := dynamo.New(sess)

// 	return db, nil
// }

func main() {
	lambda.Start(handler)
}
