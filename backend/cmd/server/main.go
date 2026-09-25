package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/campusbooks/campusbooks/internal/config"
	"github.com/campusbooks/campusbooks/internal/constants"
	"github.com/campusbooks/campusbooks/internal/database"
	"github.com/campusbooks/campusbooks/internal/router"
	"github.com/campusbooks/campusbooks/internal/util"
)

func main() {
	logger := util.NewLogger()
	cfg := config.Load()

	db, err := database.Connect(cfg, logger)
	if err != nil {
		logger.Error(constants.LogDBConnectFailed, "error", err)
		os.Exit(1)
	}
	if err := database.Migrate(db, logger); err != nil {
		logger.Error(constants.LogDBMigrateFailed, "error", err)
		os.Exit(1)
	}
	if err := database.Seed(db, logger); err != nil {
		logger.Error(constants.LogDBSeedFailed, "error", err)
		os.Exit(1)
	}

	codeStore := util.NewCodeStore(cfg, logger)
	minioClient, err := util.NewMinIOClient(cfg, logger)
	if err != nil {
		logger.Error(constants.LogMinIOConnectFailed, "error", err)
		os.Exit(1)
	}

	r := router.Setup(cfg, db, codeStore, minioClient, logger)

	srv := &http.Server{
		Addr:              ":" + cfg.ServerPort,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Info(constants.LogServerStarted, "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error(constants.LogServerStopped, "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	logger.Info(constants.LogServerShutdown)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error(constants.LogServerShutdownFailed, "error", err)
	}
}
