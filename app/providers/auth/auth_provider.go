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
	authProvider *controllers.AuthController
	authOnce     sync.Once
)

func GetAuthProvider() *controllers.AuthController {
	authOnce.Do(func() {
		db := database.DBConn
		if db == nil {
			log.Fatal("Database connection is not initialized")
		}

		authRepo := repositories.NewAuthRepository(db)
		authService := services.NewAuthService(authRepo)
		authProvider = controllers.NewAuthController(authService)
	})
	return authProvider
}
