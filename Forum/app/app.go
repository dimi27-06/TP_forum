package app

import (
	"database/sql"

	"exemple_api/config"
	"exemple_api/controllers"
	"exemple_api/repositories"
	"exemple_api/routers"
	"exemple_api/services"

	"github.com/gorilla/mux"
)

type App struct {
	Db     *sql.DB
	Router *mux.Router
}

func InitApp() *App {
	config.LoadEnv()

	db := config.InitDB()

	userRepository := repositories.InitUserRepository(db)
	forumRepository := repositories.InitForumRepository(db)

	authService := services.InitAuthService(userRepository)
	forumService := services.InitForumService(forumRepository)

	authController := controllers.AuthProductController(authService)
	forumController := controllers.InitForumControllers(forumService)

	router := mux.NewRouter()
	registerWebRoutes(router)

	apiRouter := router.PathPrefix("/api").Subrouter()
	routers.AuthProductRoutes(apiRouter, authController)
	routers.RegisterForumRoutes(apiRouter, forumController)

	return &App{
		Db:     db,
		Router: router,
	}
}

func (a *App) Close() {
	if a.Db != nil {
		a.Db.Close()
	}
}
