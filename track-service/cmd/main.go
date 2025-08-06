package main

import (
	"log"
	"net/http"
	"os"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/Vovarama1992/go-utils/logger"
	service "github.com/Vovarama1992/retry/track-service/internal/domain"
	"github.com/Vovarama1992/retry/track-service/internal/infra/postgres"
	delivery "github.com/Vovarama1992/retry/track-service/internal/visit/delivery"
	visitdomain "github.com/Vovarama1992/retry/track-service/internal/visit/domain"
	visitinfra "github.com/Vovarama1992/retry/track-service/internal/visit/infra"
	"github.com/go-chi/chi/v5"

	_ "github.com/Vovarama1992/retry/track-service/docs"

	"go.uber.org/zap"
)

// @title Track Service API
// @version 1.0
// @description Сервис для отслеживания визитов и действий
// @BasePath /
func main() {
	// logger
	zapBase, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("cannot init zap: %v", err)
	}
	defer zapBase.Sync()

	l := logger.NewZapLogger(zapBase.Sugar())
	l.Log(logger.LogEntry{
		Level:   "info",
		Message: "track service starting",
		Service: "track",
		Method:  "main",
	})

	// db + breaker
	db := postgres.NewPgConn()
	defer db.Close()

	breaker := postgres.NewPgBreaker()

	// repos + services
	actionRepo := postgres.NewActionRepo(db, breaker)
	visitRepo := visitinfra.NewVisitRepo(db, breaker)

	visitService := visitdomain.NewVisitService(visitRepo)
	trackService := service.NewTrackService(actionRepo, visitService)

	// delivery
	handler := delivery.NewHandler(trackService, visitService)
	r := chi.NewRouter()
	delivery.RegisterRoutes(r, handler)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	// swagger
	r.Handle("/swagger/", httpSwagger.WrapHandler)

	// run
	addr := ":" + os.Getenv("TRACK_SERVICE_PORT")
	l.Log(logger.LogEntry{
		Level:   "info",
		Message: "http listening at " + addr,
		Service: "track",
		Method:  "main",
	})

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
