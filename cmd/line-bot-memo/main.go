package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/line/line-bot-sdk-go/linebot"

	"github.com/tokatu4561/memo-app-use/line"
	"github.com/tokatu4561/memo-app-use/models"
)

// TODO: env管理する
const AWS_REGION = "ap-northeast-1"
const DYNAMO_ENDPOINT = "http://dynamodb:8000"

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	line, err := setUpLineClient()
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "LINE接続エラー", StatusCode: 500}, err
	}

	lineEvents, err := line.ParseRequest(request)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "LINE接続エラー", StatusCode: 500}, err
	}

	db, _ := setUpDB()

	for _, event := range lineEvents {
		// イベントがメッセージの受信だった場合
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {

			case *linebot.TextMessage:
				replyMessage := message.Text
				memo := models.Memo{}
				memo.Insert(db, replyMessage)
				_, err = line.Client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do()
				if err != nil {
					return events.APIGatewayProxyResponse{
						Body:       err.Error(),
						StatusCode: 500,
					}, nil
				}
			default:
			}
		}
	}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Hello, %v", string("hello")),
		StatusCode: 200,
	}, nil
}


func setUpLineClient() (*line.Line, error) {
	line := &line.Line{
		ChannelSecret: os.Getenv("LINE_BOT_CHANNEL_SECRET"),
		ChannelToken:  os.Getenv("LINE_BOT_CHANNEL_TOKEN"),
	}

	bot, err := linebot.New(
		line.ChannelSecret,
		line.ChannelToken,
	)
	if err != nil {
		return nil, err
	}

	line.Client = bot

	return line, nil
}

func setUpDB() (*dynamo.DB, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(AWS_REGION),
		Endpoint:    aws.String(DYNAMO_ENDPOINT),
		Credentials: credentials.NewStaticCredentials("dummy", "dummy", "dummy"),
	})
	if err != nil {
		return nil, err
	}

	db := dynamo.New(sess)

	return db, nil
}

func main() {
	lambda.Start(handler)
}
