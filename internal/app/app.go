package app

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/egor-denisov/wallet-infotecs/config"
	v1 "github.com/egor-denisov/wallet-infotecs/internal/controller/http/v1"
	repo "github.com/egor-denisov/wallet-infotecs/internal/repository/postgres"
	"github.com/egor-denisov/wallet-infotecs/internal/usecase"
	"github.com/egor-denisov/wallet-infotecs/pkg/logger"
	"github.com/egor-denisov/wallet-infotecs/pkg/postgres"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Connect postgres db
	pg, err := postgres.New(fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.PG.User, cfg.PG.Password, cfg.PG.Host, cfg.PG.Port, cfg.PG.DB))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	// Migrate database schema
	err = pg.Migrate("./migrations/up.sql");
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - pg.Migrate: %w", err))
	}
	defer pg.DB.Close()

	// Use case
	walletUseCase := usecase.New(
		repo.NewWalletRepo(pg),
		cfg.App.DefaultBalance,
	)

	// HTTP Server
	httpServer := gin.New()
	v1.NewRouter(httpServer, l, walletUseCase)
	
	httpServer.Run(fmt.Sprintf(":%s", cfg.HTTP.Port))
	
}
