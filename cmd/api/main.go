package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Udean777/uang-bijak-go/internal/config"
	"github.com/Udean777/uang-bijak-go/internal/handler"
	"github.com/Udean777/uang-bijak-go/internal/middleware"
	"github.com/Udean777/uang-bijak-go/internal/repository"
	"github.com/Udean777/uang-bijak-go/internal/service"
)

func main() {
	log.Println("Start application....")

	cfg := config.LoadConfig()

	dbpool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Gagal menghubungkan ke database: %v", err)
	}
	defer dbpool.Close()

	if err := dbpool.Ping(context.Background()); err != nil {
		log.Fatalf("Gagal melakukan ping ke database: %v", err)
	}
	log.Println("Berhasil terhubung ke database")

	userRepo := repository.NewUserRepository(dbpool)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	authService := service.NewAuthService(userRepo, cfg.JwtSecret, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)
	authHandler := handler.NewAuthHandler(authService)
	authMiddleware := middleware.AuthMiddleware(cfg.JwtSecret)

	categoryRepo := repository.NewCategoryRepository(dbpool)
	categoryService := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categoryService)

	walletRepo := repository.NewWalletRepository(dbpool)
	walletService := service.NewWalletService(walletRepo)
	walletHandler := handler.NewWalletHandler(walletService)

	trxRepo := repository.NewTransactionRepository(dbpool)
	trxService := service.NewTransactionService(dbpool, trxRepo, walletRepo, categoryRepo)
	trxHandler := handler.NewTransactionHandler(trxService)

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong!",
		})
	})

	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/refresh", authHandler.Refresh)
	}

	api := router.Group("/api/v1")
	api.Use(authMiddleware)
	{
		api.GET("/me", userHandler.GetMe)

		catRoutes := api.Group("/categories")
		{
			catRoutes.POST("/", categoryHandler.CreateCategory)
			catRoutes.GET("/", categoryHandler.GetUserCategories)
			catRoutes.PUT("/:id", categoryHandler.UpdateCategory)
			catRoutes.DELETE("/:id", categoryHandler.DeleteCategory)
		}

		walletRoutes := api.Group("/wallets")
		{
			walletRoutes.POST("/", walletHandler.CreateWallet)
			walletRoutes.GET("/", walletHandler.GetUserWallets)
			walletRoutes.PUT("/:id", walletHandler.UpdateWallet)
			walletRoutes.DELETE("/:id", walletHandler.DeleteWallet)
		}

		trxRoutes := api.Group("/transactions")
		{
			trxRoutes.POST("/", trxHandler.CreateTransaction)
			trxRoutes.GET("/", trxHandler.GetUserTransactions)
			// TODO: Tambahkan PUT /:id dan DELETE /:id
		}
	}

	serverAddr := ":" + cfg.AppPort
	log.Printf("Menjalankan server di %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
