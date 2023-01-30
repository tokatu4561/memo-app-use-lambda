package mock

import (
	"time"

	"github.com/guregu/dynamo"
	"github.com/tokatu4561/memo-app-use/pkg/domain"
)

type MockDatabase struct {}

type MemoRepositoryGateway struct {
	databaseHandler *MockDatabase
}

type DatabaseHandler struct {
	Conn *dynamo.DB
}

func NewMockDatabaseHandler() *MockDatabase {
	return &MockDatabase{}
}

func NewMockMemoRepository(db *MockDatabase) domain.MemoRepositoryInterface {
	return &MemoRepositoryGateway{
		databaseHandler: db,
	}
}

func (t *MemoRepositoryGateway) AddMemo(memo *domain.Memo) (*domain.Memo, error) {
	newMemo := domain.Memo{
		ID: "test1",
		Title: "テスト",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return &newMemo, nil
}

func (t *MemoRepositoryGateway) UpdateMemo(memo *domain.Memo) (*domain.Memo, error) {
	updatedMemo := domain.Memo{
		ID: "test1",
		Title: "テスト",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return &updatedMemo, nil
}

func (t *MemoRepositoryGateway) DeleteMemo(memo *domain.Memo) error {
	return nil
}

func (t *MemoRepositoryGateway) GetMemo(id string) (*domain.Memo, error) {
	memo := domain.Memo{
		ID: "test1",
		Title: "テスト",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return &memo, nil
}

func (t *MemoRepositoryGateway) GetMemos() ([]*domain.Memo, error) {
	memos := []*domain.Memo{
			{
			ID: "test1",
			Title: "テスト",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			},
			{
			ID: "test1",
			Title: "テスト",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			},
		}

	return memos, nil
}
