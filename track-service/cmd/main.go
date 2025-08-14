package main

import (
	"log"
	"net/http"
	"os"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/Vovarama1992/go-utils/logger"
	actionhttp "github.com/Vovarama1992/retry/track-service/internal/delivery"
	service "github.com/Vovarama1992/retry/track-service/internal/domain"
	"github.com/Vovarama1992/retry/track-service/internal/infra/postgres"
	sessionhttp "github.com/Vovarama1992/retry/track-service/internal/session/delivery"
	sessiondomain "github.com/Vovarama1992/retry/track-service/internal/session/domain"
	sessioninfra "github.com/Vovarama1992/retry/track-service/internal/session/infra"
	visithttp "github.com/Vovarama1992/retry/track-service/internal/visit/delivery"
	visitdomain "github.com/Vovarama1992/retry/track-service/internal/visit/domain"
	visitinfra "github.com/Vovarama1992/retry/track-service/internal/visit/infra"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	_ "github.com/Vovarama1992/retry/track-service/docs"

	"go.uber.org/zap"
)

// @title Track Service API
// @version 1.0
// @description Сервис для отслеживания визитов, сессий и действий
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

	// repos
	actionRepo := postgres.NewActionRepo(db, breaker)
	visitRepo := visitinfra.NewVisitRepo(db, breaker)
	sessionRepo := sessioninfra.NewSessionRepo(db, breaker)

	// services
	visitService := visitdomain.NewVisitService(visitRepo)
	sessionService := sessiondomain.NewSessionService(sessionRepo)
	trackService := service.NewTrackService(actionRepo, visitService)

	// delivery
	visitHandler := visithttp.NewHandler(trackService, visitService, l)
	actionHandler := actionhttp.NewHandler(trackService, l)
	sessionHandler := sessionhttp.NewHandler(sessionService, l)

	// router
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://retry.school"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// routes
	visithttp.RegisterRoutes(r, visitHandler)
	actionhttp.RegisterRoutes(r, actionHandler)
	sessionhttp.RegisterRoutes(r, sessionHandler)

	// ping
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	// swagger
	r.Handle("/swagger/*", httpSwagger.WrapHandler)

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
