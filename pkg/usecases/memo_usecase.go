package usecases

import "github.com/tokatu4561/memo-app-use/pkg/domain"

type MemoUsecase struct {
	Repository domain.MemoRepositoryInterface
}

func NewMemoUsecase(memoRepo domain.MemoRepositoryInterface) domain.MemoUseCaseInterface {
	return &MemoUsecase{
		Repository: memoRepo,
	}
}

func (t *MemoUsecase) GetMemos() ([]*domain.Memo, error) {
	memos, err := t.Repository.GetMemos()
	if err != nil {
		return nil, err
	}

	return memos, err
}

func (t *MemoUsecase) GetMemo(id string) (*domain.Memo, error) {
	memo, err := t.Repository.GetMemo(id)

	if err != nil {
		return nil, err
	}

	return memo, nil
}

func (t *MemoUsecase) AddMemo(memo *domain.Memo) (*domain.Memo, error) {
	newMemo, err := t.Repository.AddMemo(memo)

	if err != nil {
		return nil, err
	}

	return newMemo, nil
}

func (t *MemoUsecase) UpdateMemo(memo *domain.Memo) (*domain.Memo, error) {
	updatedMemo, err := t.Repository.UpdateMemo(memo)

	if err != nil {
		return nil, err
	}

	return updatedMemo, nil
}

func (t *MemoUsecase) DeleteMemo(memo *domain.Memo) error {
	err := t.Repository.DeleteMemo(memo)

	if err != nil {
		return err
	}

	return nil
}
