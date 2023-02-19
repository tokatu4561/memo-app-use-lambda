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

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	slackApi := slack.New(os.Getenv("SLACK_OAUTH_TOKEN"))

	// slack からのリクエストかどうかを検証する
	err := Verify(ConvertHeaders(request.Headers), []byte(request.Body))
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "slack conection error", StatusCode: 400}, err
	}
	
	body := request.Body
	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if err != nil {
		return events.APIGatewayProxyResponse{Body: "slack conection error", StatusCode: 500}, err
	}

	// 受け取ったslack のイベントに応じて処理
	switch eventsAPIEvent.Type {
		case slackevents.URLVerification:
			// slack に登録する の コールバック url(こちら側で処理するアプリ側のURL) が有効かどうか確かめるため
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
				case "memo":
					// 新しくメモを追加する
					memo, err := ctl.CreateMemo(msg[2])
					if err != nil {
						return events.APIGatewayProxyResponse{Body: "bad request", StatusCode: 400}, err
					}

					responseMsg := fmt.Sprintf("%sを追加しました!", memo.Title)
					_, _, err = slackApi.PostMessage(event.Channel, slack.MsgOptionText(responseMsg, false))
					if err != nil {
						return events.APIGatewayProxyResponse{Body: "bad request", StatusCode: 400}, err
					}
				case "list":
					// 保存されているメモのリストを通知する

					// 全てのメモを取得
					memos, err := ctl.GetMemos()
					if err != nil {
						return events.APIGatewayProxyResponse{Body: "bad request", StatusCode: 400}, err
					}

					if (len(memos) == 0) {
						_, _, err = slackApi.PostMessage(event.Channel, slack.MsgOptionText("メモがありません。", false))
					}

					// slack へメモのリストを 1つづつ通知
					for _, memo := range memos {
						attachment := slack.Attachment{
							Pretext:    memo.Title,
							Fallback:   "We don't currently support your client",
							CallbackID: "continue_or_complete",
							Color:      "#3AA3E3",
							Actions: []slack.AttachmentAction{
								{
									Name:  "continue",
									Text:  "継続",
									Type:  "button",
									Value: memo.ID,
								},
								{
									Name:  "complete",
									Text:  "完了",
									Type:  "button",
									Value: memo.ID,
									Style: "primary",
								},
							},
						}
						message := slack.MsgOptionAttachments(attachment)
						_, _, err = slackApi.PostMessage(event.Channel, slack.MsgOptionText("", false), message)
					}
					
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
