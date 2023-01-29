package application

import (
	"time"

	"github.com/google/uuid"
	"github.com/slack-go/slack/slackevents"
	"github.com/tokatu4561/memo-app-use/pkg/domain"
)

type SlackMemoController struct {
	memoUsecase MemoUsecase
}

func NewSlackMemoController(memoUsecase domain.MemoUseCaseInterface) *SlackMemoController {
	return &SlackMemoController{
		memoUsecase: memoUsecase,
	}
}

func (t *SlackMemoController) GetMemos(event slackevents.AppMentionEvent) ([]*domain.Memo, error) {
	memos, err := t.memoUsecase.GetMemos()
	if err != nil {
		return nil, err
	}

	return memos, nil
}

func (t *SlackMemoController) GetMemo(event *slackevents.AppMentionEvent) (*domain.Memo, error) {
	id := request.PathParameters["id"]

	memo, err := t.memoUsecase.GetMemo(id)
	if err != nil {
		return nil, err
	}

	return memo, nil
}

func (t *SlackMemoController) CreateMemo(event *slackevents.AppMentionEvent) (*domain.Memo, error) {
	type RequestPayload struct {
		Memo Memo `json:"memo"`
	}
	var requestPayload RequestPayload
	t.readJson(request, &requestPayload)

	event.InnerEvent.Data.Text

	newId := uuid.New()

	memo := domain.Memo{
		ID:        newId.String(),
		UserID:    requestPayload.Memo.UserID,
		Title:     requestPayload.Memo.Title,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	newMemo, err := t.memoUsecase.AddMemo(&memo)
	if err != nil {
		return nil, err
	}

	return newMemo, nil
}

func (t *SlackMemoController) UpdateMemo(event *slackevents.AppMentionEvent) (*domain.Memo, error) {
	type RequestPayload struct {
		Memo Memo
	}
	var requestPayload RequestPayload
	t.readJson(request, &requestPayload)

	memo := domain.Memo{
		ID:        requestPayload.Memo.ID,
		UserID:    requestPayload.Memo.UserID,
		Title:     requestPayload.Memo.Title,
		UpdatedAt: time.Now(),
	}

	updatedMemo, err := t.memoUsecase.UpdateMemo(&memo)
	if err != nil {
		return nil, err
	}

	return updatedMemo, nil
}

func (t *SlackMemoController) DeleteMemo(event *slackevents.AppMentionEvent) error {
	type RequestPayload struct {
		Memo Memo `json:"memo"`
	}
	var requestPayload RequestPayload
	t.readJson(request, &requestPayload)

	memo := domain.Memo{
		ID:     requestPayload.Memo.ID,
		UserID: requestPayload.Memo.UserID,
		Title:  requestPayload.Memo.Title,
	}

	err := t.memoUsecase.DeleteMemo(&memo)
	if err != nil {
		return err
	}

	return nil
}

// func (c *SlackMemoController) readJson(req events.APIGatewayProxyRequest, data interface{}) error {
// 	err := json.Unmarshal([]byte(req.Body), &data)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (c *MemoController) readJson(w http.ResponseWriter, r *http.Request, data interface{}) error {
// 	maxBytes := 1048576

// 	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

// 	dec := json.NewDecoder(r.Body)
// 	err := dec.Decode(data)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (c *MemoController) writeJson(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
// 	out, err := json.MarshalIndent(data, "", "\t")
// 	if err != nil {
// 		return err
// 	}

// 	if len(headers) > 0 {
// 		for k, v := range headers[0] {
// 			w.Header()[k] = v
// 		}
// 	}

// 	w.Header().Set("Content-type", "application/json")
// 	w.WriteHeader(status)
// 	w.Write(out)

// 	return nil
// }

// func (t *MemoController) badRequest(w http.ResponseWriter, err error) error {
// 	var payload struct {
// 		Error   bool   `json:"error"`
// 		Message string `json:"message"`
// 	}

// 	payload.Error = true
// 	payload.Message = err.Error()

// 	out, err := json.MarshalIndent(payload, "", "\t")
// 	if err != nil {
// 		return err
// 	}

// 	w.Header().Set("Content-type", "application/json")
// 	w.Write(out)

// 	return nil
// }
