package app

import (
	"context"
	"log"
	"net/http"

	"github.com/Basu008/Better-ESPN/config"
	"github.com/Basu008/Better-ESPN/database"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct {
	Router   *mux.Router
	Database *mongo.Database
}

func SetUpApp(c *config.Config) {
	app := new(App)
	app.Initialize(c)
	app.StartServer(c.GetServerHost())
}

func (app *App) Initialize(c *config.Config) {
	app.Database = database.ConnectToDatabase(c)
	app.Router = mux.NewRouter()
	app.setRoutes()
}

type RequestHandlerFunction func(db *mongo.Database, w http.ResponseWriter, r *http.Request)

func (app *App) handleRequest(handlerFunction RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlerFunction(app.Database, w, r)
	}
}

func (app *App) StartServer(host string) {
	log.Printf("Server is listening on port%s", host)
	log.Fatal(http.ListenAndServe(host, app.Router))
	app.Database.Client().Disconnect(context.Background())
}

func (app *App) Get(path string, endpoint http.HandlerFunc, queries ...string) {
	app.Router.HandleFunc(path, endpoint).Methods("GET").Queries(queries...)
}

func (app *App) Post(path string, endpoint http.HandlerFunc, queries ...string) {
	app.Router.HandleFunc(path, endpoint).Methods("POST").Queries(queries...)
}
func (app *App) Put(path string, endpoint http.HandlerFunc, queries ...string) {
	app.Router.HandleFunc(path, endpoint).Methods("PUT").Queries(queries...)
}
func (app *App) Delete(path string, endpoint http.HandlerFunc, queries ...string) {
	app.Router.HandleFunc(path, endpoint).Methods("DELETE").Queries(queries...)
}
