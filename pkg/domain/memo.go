package domain

import "time"

// FIXEME: domainなのにdynamoの情報を書いてる
type Memo struct {
	ID        string    `dynamo:"id" json:"id"`
	UserID    int       `dynamo:"userId" json:"user_id"`
	Title     string    `dynamo:"title" json:"title"`
	CreatedAt time.Time `dynamo:"createdAt" json:"created_at"`
	UpdatedAt time.Time `dynamo:"updatedAt" json:"updated_at"`
}

type MemoRepositoryInterface interface {
	AddMemo(memo *Memo) (*Memo, error)
	UpdateMemo(memo *Memo) (*Memo, error)
	DeleteMemo(memo *Memo) error
	GetMemo(id string) (*Memo, error)
	GetMemos() ([]*Memo, error)
}

type MemoUseCaseInterface interface {
	AddMemo(t *Memo) (*Memo, error)
	UpdateMemo(t *Memo) (*Memo, error)
	DeleteMemo(t *Memo) error
	GetMemo(id string) (*Memo, error)
	GetMemos() ([]*Memo, error)
}
