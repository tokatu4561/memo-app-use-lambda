package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/slack-go/slack"
	"github.com/tokatu4561/memo-app-use/pkg/application/di"
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

	// body を json にエンコード
	q, _ := url.ParseQuery(request.Body)
	qPayload := q.Get("payload")

	// JSONを struct に読み取る
	var payload MessageActionPayload
	err := json.Unmarshal([]byte(qPayload), &payload)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "bad request", StatusCode: 400}, err
	}

	// continue 継続であれば削除処理 せずに レスポンス返す
	if (payload.Actions[0].Name == "continue") {
		return events.APIGatewayProxyResponse{Body: "", StatusCode: 200 }, nil
	}
	
	// 対象のメモを削除する
	ctl := di.NewSlackMemoController()
	ctl.DeleteMemo(payload.Actions[0].Value)

	// 削除したことをslackのチャンネルに通知する
	responseMsg := fmt.Sprintf("%sをリストから削除しました!", "a")
	_, _, err = slackApi.PostMessage(payload.Channel.Id, slack.MsgOptionText(responseMsg, false))
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "bad request", StatusCode: 400}, err
	}
	
	return events.APIGatewayProxyResponse{
		Body:       "",
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
