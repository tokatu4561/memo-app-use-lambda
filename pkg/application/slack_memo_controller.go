package application

import (
	"time"

	"github.com/google/uuid"
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

type MemoPayload struct {
	id string
	title string
}

func (t *SlackMemoController) GetMemos() ([]*domain.Memo, error) {
	memos, err := t.memoUsecase.GetMemos()
	if err != nil {
		return nil, err
	}

	return memos, nil
}

func (t *SlackMemoController) GetMemo(id string) (*domain.Memo, error) {
	memo := &domain.Memo{
		ID: "test1",
		Title: "テスト",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return memo, nil
}

func (t *SlackMemoController) CreateMemo(text string) (*domain.Memo, error) {
	msg := text

	newId := uuid.New()

	memo := domain.Memo{
		ID:        newId.String(),
		Title:     msg,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	newMemo, err := t.memoUsecase.AddMemo(&memo)
	if err != nil {
		return nil, err
	}

	return newMemo, nil
}

func (t *SlackMemoController) UpdateMemo(memo MemoPayload) (*domain.Memo, error) {
	updatedMemo := &domain.Memo{
		ID: "test1",
		Title: "テスト",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return updatedMemo, nil
}

func (t *SlackMemoController) DeleteMemo(id string) error {
	// msg := strings.Split(event.Text, " ")

	memo := domain.Memo{
		ID: "test1",
		Title: "テスト",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := t.memoUsecase.DeleteMemo(&memo)
	if err != nil {
		return err
	}

	return nil
}
