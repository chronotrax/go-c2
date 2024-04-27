package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/chronotrax/go-c2/internal/server/config"
	"github.com/chronotrax/go-c2/internal/server/embed"
	"github.com/chronotrax/go-c2/internal/server/handler"
	"github.com/chronotrax/go-c2/internal/server/logging"
	"github.com/chronotrax/go-c2/internal/server/model/sqliteDB"
	"github.com/chronotrax/go-c2/internal/server/route"
	"github.com/chronotrax/go-c2/pkg/msgqueue"

	"github.com/labstack/echo/v4"
)

func main() {
	// Config
	conf, err := config.InitConfig()
	if err != nil {
		slog.Error("failed getting config, using defaults", slog.String("error", err.Error()))
	}
	//goland:noinspection GoDfaErrorMayBeNotNil
	slog.Info("using config:", slog.String("conf", conf.String()))

	// Echo router
	e := echo.New()
	e.Debug = conf.Debug

	// Logging
	logging.InitLogging(e)

	// Dependencies
	agentDB := sqliteDB.Connect(embed.EmbededMigrations, conf.DBName)
	msgQueue := msgqueue.NewMsgQueue()
	d := handler.NewDepends(sqliteDB.NewAgentDB(agentDB), msgQueue)

	// Routes
	route.InitRoutes(e, d)

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Start server
	go func() {
		if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
