package providers

import (
	"log"
	"sync"
	database "template-auth/app/config"
	controllers "template-auth/app/http/controllers/auth"
	repositories "template-auth/app/repositories/auth"
	services "template-auth/app/services/auth"
)

var (
	forgotProvider *controllers.ForgotController
	forgotOnce     sync.Once
)

func GetForgotProvider() *controllers.ForgotController {
	forgotOnce.Do(func() {
		db := database.DBConn
		if db == nil {
			log.Fatal("Database connection is not initialized")
		}

		forgotRepo := repositories.NewForgotRepository(db)
		authRepo := repositories.NewAuthRepository(db)
		forgotService := services.NewForgotService(forgotRepo, authRepo)
		forgotProvider = controllers.NewForgotController(forgotService)
	})
	return forgotProvider
}
