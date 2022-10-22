package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/line/line-bot-sdk-go/linebot"
)

// 初期作成されるテンプレート 参考のため残しておく
// var (
// 	// DefaultHTTPGetAddress Default Address
// 	DefaultHTTPGetAddress = "https://checkip.amazonaws.com"

// 	// ErrNoIP No IP found in response
// 	ErrNoIP = errors.New("No IP in HTTP response")

// 	// ErrNon200Response non 200 status code in response
// 	ErrNon200Response = errors.New("Non 200 Response found")
// )

// resp, err := http.Get(DefaultHTTPGetAddress)
// if err != nil {
// 	return events.APIGatewayProxyResponse{}, err
// }

// if resp.StatusCode != 200 {
// 	return events.APIGatewayProxyResponse{}, ErrNon200Response
// }

// ip, err := ioutil.ReadAll(resp.Body)
// if err != nil {
// 	return events.APIGatewayProxyResponse{}, err
// }

// if len(ip) == 0 {
// 	return events.APIGatewayProxyResponse{}, ErrNoIP
// }

// TODO: env管理する
const AWS_REGION = "ap-northeast-1"
const DYNAMO_ENDPOINT = "http://dynamodb:8000"

type User struct {
	UserID string `dynamo:"UserID,hash"`
	Name   string `dynamo:"Name,range"`
	Age    int    `dynamo:"Age"`
	Text   string `dynamo:"Text"`
}

type Line struct {
	ChannelSecret string
	ChannelToken  string
	Client        *linebot.Client
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	line, err := setUpLineClient()
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "LINE接続エラー", StatusCode: 500}, err
	}

	lineEvents, err := line.ParseRequest(request)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "LINE接続エラー", StatusCode: 500}, err
	}

	for _, event := range lineEvents {
		// イベントがメッセージの受信だった場合
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {

			case *linebot.TextMessage:
				replyMessage := message.Text
				_, err = line.Client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do()
				if err != nil {
					return events.APIGatewayProxyResponse{}, err
				}
			default:
			}
		}
	}

	// sess, err := session.NewSession(&aws.Config{
	// 	Region:      aws.String(AWS_REGION),
	// 	Endpoint:    aws.String(DYNAMO_ENDPOINT),
	// 	Credentials: credentials.NewStaticCredentials("dummy", "dummy", "dummy"),
	// })
	// if err != nil {
	// 	return events.APIGatewayProxyResponse{
	// 		Body:       err.Error(),
	// 		StatusCode: 500,
	// 	}, nil
	// }

	// db := dynamo.New(sess)

	// db.Table("UserTable").DeleteTable().Run()

	// err = db.CreateTable("UserTable", User{}).Run()
	// if err != nil {
	// 	return events.APIGatewayProxyResponse{
	// 		Body:       err.Error(),
	// 		StatusCode: 500,
	// 	}, nil
	// }
	// table := db.Table("UserTable")

	// var user User

	// err = table.Put(&User{UserID: "1234", Name: "太郎", Age: 20}).Run()
	// if err != nil {
	// 	return events.APIGatewayProxyResponse{
	// 		Body:       err.Error(),
	// 		StatusCode: 500,
	// 	}, nil
	// }

	// err = table.Get("UserID", "1234").Range("Name", dynamo.Equal, "太郎").One(&user)
	// if err != nil {
	// 	return events.APIGatewayProxyResponse{
	// 		Body:       err.Error(),
	// 		StatusCode: 500,
	// 	}, nil
	// }
	// fmt.Printf("GetDB%+v\n", user)

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Hello, %v", string("hello")),
		StatusCode: 200,
	}, nil
}

func (l *Line) ParseRequest(r events.APIGatewayProxyRequest) ([]*linebot.Event, error) {
	req := &struct {
		Events []*linebot.Event `json:"events"`
	}{}
	if err := json.Unmarshal([]byte(r.Body), req); err != nil {
		return nil, err
	}

	return req.Events, nil
}

func setUpLineClient() (*Line, error) {
	line := &Line{
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

func main() {
	lambda.Start(handler)
}
