package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/slack-go/slack"
)

type Actions struct {
	Name string `json:"name"`
	Value string `json:"value"`
	ActionType string `json:"type"`
}

type MessageActionPayload struct {
	Actions []*Actions `json:"actions"`
	CallbackId string `json:"callback_id"`
	Channel struct {
		Id string `json:"id"`
		Name string `json:"name"`
	} `json:"channel"`
	User struct {
		Id string `json:"id"`
		Name string `json:"name"`
	} `json:"User"`
	ActionTs string `json:"action_ts"`
	MessageTs string `json:"message_ts"`
	AttachmentId string `json:"attachment_id"`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	slackApi := slack.New(os.Getenv("SLACK_OAUTH_TOKEN"))

	// JSONを読み取る
	type RequestPayload struct {
		Payload MessageActionPayload `json:"payload"`
	}
	var payload RequestPayload
	// readJson(request, &payload)

	_ = json.Unmarshal([]byte(request.Body), &payload)

	fmt.Printf("Request body. %+v\n", request.Body)
	log.Println(payload)

	// 対象のメモを削除する
	// ctl := di.NewSlackMemoController()
	// ctl.DeleteMemo()

	// 削除したことをslackのチャンネルに通知する
	responseMsg := fmt.Sprintf("%sを追加しました!", "a")
	_, _, err := slackApi.PostMessage(payload.Payload.Channel.Id, slack.MsgOptionText(responseMsg, false))
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "bad request", StatusCode: 400}, err
	}
	
	return events.APIGatewayProxyResponse{
		Body:       "",
		StatusCode: 200,
	}, nil
}

func readJson(req events.APIGatewayProxyRequest, data interface{}) error {
	err := json.Unmarshal([]byte(req.Body), &data)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
