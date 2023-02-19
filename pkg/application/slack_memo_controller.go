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

func (t *SlackMemoController) GetMemos() ([]*domain.Memo, error) {
	memos, err := t.memoUsecase.GetMemos()
	if err != nil {
		return nil, err
	}

	return memos, nil
}

func (t *SlackMemoController) CreateMemo(text string) (*domain.Memo, error) {
	newId := uuid.New()

	memo := domain.Memo{
		ID:        newId.String(),
		Title:     text,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	newMemo, err := t.memoUsecase.AddMemo(&memo)
	if err != nil {
		return nil, err
	}

	return newMemo, nil
}

func (t *SlackMemoController) DeleteMemo(id string) error {
	memo := domain.Memo{
		ID: id,
	}

	err := t.memoUsecase.DeleteMemo(&memo)
	if err != nil {
		return err
	}

	return nil
}
