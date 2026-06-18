package app

import (
	"database/sql"
	"forum/config"
	"forum/controllers"
	"forum/middleware"
	"forum/repositories"
	"forum/router"
	"forum/services"
	"forum/templates"

	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	db     *sql.DB
}

func InitApp() *App {
	config.LoadEnv()
	db := config.InitDB()

	userRepo := repositories.InitUserRepository(db)
	tagRepo := repositories.InitTagRepository(db)
	threadRepo := repositories.InitThreadRepository(db)
	messageRepo := repositories.InitMessageRepository(db)
	reactionRepo := repositories.InitReactionRepository(db)

	authService := services.InitAuthService(userRepo)
	threadService := services.InitThreadService(threadRepo, tagRepo)
	messageService := services.InitMessageService(messageRepo, threadRepo)
	reactionService := services.InitReactionService(reactionRepo)
	adminService := services.InitAdminService(userRepo, threadRepo)

	tmpl := templates.NewManager()

	authCtrl := controllers.InitAuthController(authService, tmpl)
	threadCtrl := controllers.InitThreadController(threadService, messageService, tmpl)
	messageCtrl := controllers.InitMessageController(messageService, tmpl)
	reactionCtrl := controllers.InitReactionController(reactionService)
	adminCtrl := controllers.InitAdminController(adminService, threadService, messageService, tmpl)

	r := mux.NewRouter()
	r.Use(middleware.LoadUser)

	router.RegisterAssetRoutes(r)
	router.RegisterAuthRoutes(r, authCtrl)
	router.RegisterThreadRoutes(r, threadCtrl)
	router.RegisterMessageRoutes(r, messageCtrl, reactionCtrl)
	router.RegisterAdminRoutes(r, adminCtrl)

	return &App{Router: r, db: db}
}

func (a *App) Close() {
	if a.db != nil {
		a.db.Close()
	}
}
