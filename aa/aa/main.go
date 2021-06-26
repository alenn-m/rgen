package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	// [services]
	"github.com/aa/aa/api/user"

	"github.com/aa/aa/api/post"

	"github.com/aa/aa/api/auth"
	"github.com/jmoiron/sqlx"
	"github.com/patrickmn/go-cache"

	mdw "github.com/aa/aa/middleware"
	authService "github.com/aa/aa/util/auth"
	"github.com/aa/aa/util/cache/memory"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var seed bool

func main() {
	flag.BoolVar(&seed, "seed", false, "Seed the database")

	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_TABLE")),
	)
	if err != nil {
		panic(err)
	}

	db.MapperFunc(func(s string) string {
		return s
	})

	// memory cache can be replaced with any other type of cache
	c := cache.New(time.Minute*60*24, 10*time.Minute)
	caching := memory.NewMemoryCache(c)

	authSvc := authService.NewAuthService(caching)
	mdw.AuthSvc = authSvc

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/api", func(r chi.Router) {
		// authentication
		auth.New(r, db, authSvc)

		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(authService.TokenAuth), mdw.AuthMiddleware)

			// [protected routes]

			// users
			user.New(r, db, authSvc)

			// posts
			post.New(r, db, authSvc)

		})

		// [public routes]
	})

	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("HTTP_PORT")), r)
}
