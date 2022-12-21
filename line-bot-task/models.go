package main

import (
	"time"

	"github.com/guregu/dynamo"
)

type Memo struct {
	Id    string `dynamo:"Id,hash"`
	Text      string `dynamo:"Text"`
	CreatedAt time.Time `dynamo:"CreatedAt"`
	UpdatedAt time.Time `dynamo:"UpdatedAt"`
}

func (m *Memo) Insert(db *dynamo.DB, text string) (*Memo, error) {
	table := db.Table("Memo")

	newMemo := &Memo{
		Id: "1234", 
		Text: text, 
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := table.Put(&newMemo).Run()
	if err != nil {
		return nil, err
	}

	return newMemo, nil
}

func (m *Memo) Update(db *dynamo.DB, memo *Memo) (*Memo, error) {
	table := db.Table("Task")

	var updatedMemo *Memo

	err := table.Update("id", memo.Id).Set("text", memo.Text).Set("updatedAt", memo.UpdatedAt).Value(&updatedMemo)
	if err != nil {
		return nil, err
	}

	return updatedMemo, nil
}

func Delete(db *dynamo.DB, memo *Memo) error {
	table := db.Table("Task")

	err := table.Delete("id", memo.Id).Run()
	if err != nil {
		return err
	}

	return nil
}

func Get(db *dynamo.DB, id string) (*Memo, error) {
	table := db.Table("Task")

	var memo *Memo

	err := table.Get("id", id).One(&memo)
	if err != nil {
		return nil, err
	}

	return memo, nil
}

func GetAll(db *dynamo.DB) ([]*Memo, error) {
	table := db.Table("Task")

	var memos []*Memo

	err := table.Scan().All(&memos)
	if err != nil {
		return nil, err
	}

	return memos, nil
}