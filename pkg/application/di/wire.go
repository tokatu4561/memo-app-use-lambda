package di

// import (
// 	"github.com/google/wire"

// 	"github.com/tokatu4561/memo-app-use/pkg/application"
// 	"github.com/tokatu4561/memo-app-use/pkg/infrastructure/dynamo"
// 	"github.com/tokatu4561/memo-app-use/pkg/usecases"
// )

// func NewMemoController() *application.MemoController {
// 	wire.Build(application.NewMemoController,
// 		usecases.NewMemoUsecase,
// 		dynamo.NewMemoRepository,
// 		dynamo.NewDynamoDatabaseHandler,
// 	)

// 	return &application.MemoController{}
// }

// func NewSlackMemoController() *application.SlackMemoController {
// 	wire.Build(application.NewSlackMemoController,
// 		usecases.NewMemoUsecase,
// 		dynamo.NewMemoRepository,
// 		dynamo.NewDynamoDatabaseHandler,
// 	)

// 	return &application.SlackMemoController{}
// }