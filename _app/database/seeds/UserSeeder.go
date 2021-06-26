package seeds

import (
    "{{Root}}/models"
    authService "{{Root}}/util/auth"
    "github.com/jmoiron/sqlx"
    "golang.org/x/crypto/bcrypt"
)

type UserSeeder struct {}

func (u *UserSeeder) Execute(db *sqlx.DB) error {
    password, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    authSvc := new(authService.AuthService)

    user := models.User{
        Email:     "admin@example.com",
        FirstName: "Admin",
        LastName:  "Admin",
        Password:  string(password),
    }

    token, err := authSvc.GenerateToken(&user)
    if err != nil {
        return err
    }

    user.ApiToken = token.Raw

    _, err = db.NamedExec("INSERT INTO Users (ApiToken, Email, FirstName, LastName, Password) VALUES (:ApiToken, :Email, :FirstName, :LastName, :Password)", user)

    return err
}
