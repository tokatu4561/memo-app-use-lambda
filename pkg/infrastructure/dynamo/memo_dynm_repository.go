package dynamo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/tokatu4561/memo-app-use/pkg/domain"
)

//FIXME:環境変数の管理
const AWS_REGION = "ap-northeast-1"
const DYNAMO_ENDPOINT = "http://dynamodb:8000"

type MemoRepositoryGateway struct {
	databaseHandler *dynamo.DB
}

type DatabaseHandler struct {
	Conn *dynamo.DB
}

func NewDynamoDatabaseHandler() *dynamo.DB {
	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String(AWS_REGION),
		// Endpoint:    aws.String(os.Getenv("DYNAMODB_ENDPOINT")), //FIXME: ローカル接続のためには必要？？s
		// Credentials: credentials.NewStaticCredentials("dummy", "dummy", "dummy"),
	})

	return dynamo.New(sess)
}

func NewMemoRepository(db *dynamo.DB) domain.MemoRepositoryInterface {
	return &MemoRepositoryGateway{
		databaseHandler: db,
	}
}

func (t *MemoRepositoryGateway) AddMemo(memo *domain.Memo) (*domain.Memo, error) {
	newMemo, err := Insert(t.databaseHandler, memo)
	if err != nil {
		return nil, err
	}

	return newMemo, nil
}

func (t *MemoRepositoryGateway) UpdateMemo(memo *domain.Memo) (*domain.Memo, error) {
	updatedMemo, err := Update(t.databaseHandler, memo)

	if err != nil {
		return nil, err
	}

	return updatedMemo, nil
}

func (t *MemoRepositoryGateway) DeleteMemo(memo *domain.Memo) error {
	err := Delete(t.databaseHandler, memo)
	if err != nil {
		return err
	}

	return nil
}

func (t *MemoRepositoryGateway) GetMemo(id string) (*domain.Memo, error) {
	memo, err := Get(t.databaseHandler, id)
	if err != nil {
		return nil, err
	}

	return memo, nil
}

func (t *MemoRepositoryGateway) GetMemos() ([]*domain.Memo, error) {
	memos, err := GetAll(t.databaseHandler)

	if err != nil {
		return nil, err
	}

	return memos, nil
}

func Insert(db *dynamo.DB, memo *domain.Memo) (*domain.Memo, error) {
	table := db.Table("Memo")

	newMemo := memo

	err := table.Put(&newMemo).Run()
	if err != nil {
		return nil, err
	}

	return newMemo, nil
}

func Update(db *dynamo.DB, memo *domain.Memo) (*domain.Memo, error) {
	table := db.Table("Memo")

	var updatedMemo *domain.Memo

	err := table.Update("id", memo.ID).Set("title", memo.Title).Set("updatedAt", memo.UpdatedAt).Value(&updatedMemo)
	if err != nil {
		return nil, err
	}

	return updatedMemo, nil
}

func Delete(db *dynamo.DB, memo *domain.Memo) error {
	table := db.Table("Memo")

	err := table.Delete("id", memo.ID).Run()
	if err != nil {
		return err
	}

	return nil
}

func Get(db *dynamo.DB, id string) (*domain.Memo, error) {
	table := db.Table("Memo")

	var memo *domain.Memo

	err := table.Get("id", id).One(&memo)
	if err != nil {
		return nil, err
	}

	return memo, nil
}

func GetAll(db *dynamo.DB) ([]*domain.Memo, error) {
	table := db.Table("Memo")

	var memos []*domain.Memo

	err := table.Scan().All(&memos)
	if err != nil {
		return nil, err
	}

	return memos, nil
}
