package infrastructure

import (
	"database/sql"
)

type MemoRepositoryGateway struct {
	databaseHandler *sql.DB
}

type DatabaseHandler struct {
	Conn *sql.DB
}

// func NewMemoRepository(db *sql.DB) domain.MemoRepositoryInterface {
// 	return &MemoRepositoryGateway{
// 		databaseHandler: db,
// 	}
// }

// func (t *MemoRepositoryGateway) AddMemo(memo *domain.Memo) (*domain.Memo, error) {
// 	newMemo, err := Insert(t.databaseHandler, memo)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return newMemo, nil
// }

// func (t *MemoRepositoryGateway) GetMemos() ([]*domain.Memo, error) {
// 	memos, err := GetAll(t.databaseHandler)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return memos, nil
// }

// func Insert(db *sql.DB, memo *domain.Memo) (*domain.Memo, error) {
// 	stmt := `insert into memos (user_id, title, created_at, updated_at)
// 		values ($1, $2, $3, $4) returning id`

// 	_, err := db.Exec(stmt,
// 		memo.UserID,
// 		memo.Title,
// 		time.Now(),
// 		time.Now(),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return nil, nil
// }

// // GetALl returns all memos in db
// func GetAll(db *sql.DB) ([]*domain.Memo, error) {
// 	query := `select id, user_id, title, created_at, updated_at from memos`

// 	var memos []*domain.Memo
// 	rows, err := db.Query(query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var memo domain.Memo
// 		err := rows.Scan(
// 			&memo.ID,
// 			&memo.UserID,
// 			&memo.Title,
// 			&memo.CreatedAt,
// 			&memo.UpdatedAt,
// 		)
// 		if err != nil {
// 			log.Println("Error scanning", err)
// 			return nil, err
// 		}

// 		memos = append(memos, &memo)
// 	}

// 	return memos, nil
// }
