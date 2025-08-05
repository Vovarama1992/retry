package main

import (
	"log"

	"github.com/Vovarama1992/go-utils/logger"
	"go.uber.org/zap"
)

func main() {
	zapBase, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("cannot init zap: %v", err)
	}
	defer zapBase.Sync()

	l := logger.NewZapLogger(zapBase.Sugar())

	l.Log(logger.LogEntry{
		Level:   "info",
		Message: "retry backend up",
		Service: "retry",
		Method:  "main",
	})
}
