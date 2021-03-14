package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	// [services]
	"{{Root}}/api/auth"
	authDB "{{Root}}/api/auth/repositories/mysql"
	"{{Root}}/api/user"
	userDB "{{Root}}/api/user/repositories/mysql"
	"{{Root}}/database/seeds"
	"github.com/jmoiron/sqlx"
	"github.com/patrickmn/go-cache"

	mdw "{{Root}}/middleware"
	authService "{{Root}}/util/auth"
	"{{Root}}/util/cache/memory"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	_ "github.com/jinzhu/gorm/dialects/mysql"
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

	if seed {
		fmt.Println("Seeding database")

		s := seeds.NewDatabaseSeeder(db)
		err = s.Run()
		if err != nil {
			panic(err)
		}

		fmt.Println("Done seeding database")

		return
	}

	// memory cache can be replaced with any other type of cache
	c := cache.New(time.Minute*60*24, 10*time.Minute)
	caching := memory.NewMemoryCache(c)

	authSvc := authService.NewAuthService(caching)
	mdw.AuthSvc = authSvc

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/api", func(r chi.Router) {
		// authentication
		auth.New(r, auth.NewController(
			authDB.NewAuthDB(db), authSvc),
		)

		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(authService.TokenAuth), mdw.AuthMiddleware)
			// users
			user.New(r, user.NewController(
				userDB.NewUserDB(db), authSvc),
			)

			// [protected routes]
		})

		// [public routes]
	})

	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("HTTP_PORT")), r)
}
