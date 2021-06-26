package test

import (
	"context"

	"github.com/test/testApp/models"
	authService "github.com/test/testApp/util/auth"
)

type DBRepository interface {
	Insert(*models.Test) (models.TestID, error)
	Update(*models.Test) error
}

type TestsController struct {
	db   DBRepository
	auth *authService.AuthService
}

type Repository interface {
	Store(context.Context, *StoreReq) (models.TestID, error)
	Update(context.Context, *UpdateReq, models.TestID) error
}

func NewController(db DBRepository, auth *authService.AuthService) *TestsController {
	return &TestsController{
		db:   db,
		auth: auth,
	}
}
