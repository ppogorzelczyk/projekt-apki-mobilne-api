package config

import (
	"buymeagiftapi/internal/auth"
	"buymeagiftapi/internal/database"
	listitems "buymeagiftapi/internal/listItems"
	"buymeagiftapi/internal/lists"
	"buymeagiftapi/internal/sharing"
	"buymeagiftapi/internal/users"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func InitGin() *gin.Engine {
	router := gin.Default()
	db := configDatabase()

	// Middleware
	router.Use(database.DatabaseMiddleware(*db))
	// Set up routing
	configRoutes(*db, router)
	// cors?

	return router
}

func configDatabase() *database.Database {
	db := database.NewDatabaseConnection()

	db.Ping()

	return db
}

func configRoutes(db database.Database, r *gin.Engine) {
	configV1Group(db, r)

	slog.Info("Routes configured")
}

func configV1Group(db database.Database, r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		configPublicEndpoints(v1)
		configPrivateEndpoints(db, v1)
	}
}

func configPublicEndpoints(r *gin.RouterGroup) {
	userController := users.UsersController{}

	// User routes
	r.POST("/users/register", userController.Register)
	r.POST("/users/login", userController.Login)
}

func configPrivateEndpoints(db database.Database, r *gin.RouterGroup) {
	listController := lists.ListsController{}
	listItemController := listitems.ListItemsController{}
	shareController := sharing.SharingController{}
	userController := users.UsersController{}

	r.Use(auth.AuthMiddleware(db.GetConnection()))
	// List routes
	r.POST("/lists", listController.Create)
	r.GET("/lists", listController.GetMyLists)
	r.GET("/lists/:listId", listController.GetList)
	r.PATCH("/lists/:listId", listController.Update)
	r.DELETE("/lists/:listId", listController.Delete)

	// Item routes
	r.POST("/lists/:listId/items", listItemController.Create)
	r.POST("/lists/:listId/items/:itemId/assign", listItemController.Assign)
	r.DELETE("/lists/:listId/items/:itemId", listItemController.Delete)
	r.PATCH("/lists/:listId/items/:itemId", listItemController.Update)

	// Sharing routes
	r.POST("/lists/:listId/share", shareController.Share)
	r.GET("/lists/:listId/share", shareController.GetSharedUsers)
	r.GET("/lists/shared", shareController.GetSharedLists)
	r.DELETE("/lists/:listId/share/:userId", shareController.Unshare)

	// User routes
	r.GET("/users/logout", userController.Logout)
	r.DELETE("/users", userController.Delete)
	// r.PATCH("/users", updateUser)
}
