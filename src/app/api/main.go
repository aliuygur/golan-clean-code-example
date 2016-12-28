package main

import (
	"app/infra"
	"app/interfaces"
	"app/interfaces/errs"
	"app/interfaces/handlers"
	"app/interfaces/repos/gormdb"
	"app/usecases"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"log"

	bugsnag "github.com/bugsnag/bugsnag-go"
	"github.com/gorilla/mux"
	"github.com/rs/cors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	durl, err := parseDBURL(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open("mysql", durl)
	if err != nil {
		log.Fatalf("cannot connect to db, err:%s", err)
	}

	if os.Getenv("AUTO_MIGRATE") == "yes" {
		db.LogMode(true)
		interfaces.InitDB(db)
	}

	// Dependencies
	mail := infra.NewFakeMail()
	errH := &errs.Handler{Debug: "on"}

	// Repos
	gormRepo := gormdb.NewRepo(db)
	catalogRepo := gormdb.NewCatalog(gormRepo)

	// middlewares
	authReqMid := interfaces.NewAuthRequiredMid(errH)
	setUserMid := interfaces.NewSetUserMid(gormRepo, errH)

	// services
	userSrv := usecases.NewUser(gormRepo, mail)
	catalogSrv := usecases.NewCatalog(catalogRepo)

	// handlers
	authH := handlers.NewAuthHandler(userSrv, errH)
	accountH := handlers.NewAccount(userSrv, errH)
	catalogH := handlers.NewCatalog(catalogSrv, errH)

	r := mux.NewRouter()
	// r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	gores.String(w, 200, "Okey")
	// })

	authH.SetRoutes(r)
	accountH.SetRoutes(r, authReqMid)
	catalogH.SetRoutes(r)

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("clients/web")))

	bugsnag.Configure(bugsnag.Configuration{
		APIKey:              os.Getenv("BUGSNAG_API_KEY"),
		ReleaseStage:        os.Getenv("ENV"),
		NotifyReleaseStages: []string{"production"},
	})

	corsMid := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	h := bugsnag.Handler(corsMid.Handler(setUserMid(r)))

	log.Printf("server starting port: %s", os.Getenv("PORT"))
	if err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), h); err != nil {
		log.Fatal(err)
	}
}

func parseDBURL(s string) (string, error) {
	durl, err := url.Parse(s)
	if err != nil {
		return "", fmt.Errorf("cannot parse database url, err:%s", err)
	}
	user := durl.User.Username()
	password, _ := durl.User.Password()
	host := durl.Host
	dbname := durl.Path // like: /path

	return fmt.Sprintf("%s:%s@tcp(%s)%s?charset=utf8&parseTime=True&loc=Local", user, password, host, dbname), nil
}
