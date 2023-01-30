// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"github.com/tokatu4561/memo-app-use/pkg/application"
	"github.com/tokatu4561/memo-app-use/pkg/infrastructure/dynamo"
	"github.com/tokatu4561/memo-app-use/pkg/infrastructure/mock"
	"github.com/tokatu4561/memo-app-use/pkg/usecases"
)

// Injectors from wire.go:

func NewMemoController() *application.MemoController {
	db := dynamo.NewDynamoDatabaseHandler()
	memoRepositoryInterface := dynamo.NewMemoRepository(db)
	memoUseCaseInterface := usecases.NewMemoUsecase(memoRepositoryInterface)
	memoController := application.NewMemoController(memoUseCaseInterface)
	return memoController
}

// Injectors from wire.go:

func NewSlackMemoController() *application.SlackMemoController {
	// db := dynamo.NewDynamoDatabaseHandler()
	// memoRepositoryInterface := dynamo.NewMemoRepository(db)
	db := mock.NewMockDatabaseHandler()
	memoRepositoryInterface := mock.NewMockMemoRepository(db)
	memoUseCaseInterface := usecases.NewMemoUsecase(memoRepositoryInterface)
	memoController := application.NewSlackMemoController(memoUseCaseInterface)
	return memoController
}
