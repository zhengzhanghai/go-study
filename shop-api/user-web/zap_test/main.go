package main

import (
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	url := "http://localhost:8080"
	sugar := logger.Sugar()
	sugar.Infow("failed to fetch URL", "url", url, "attempt", 3)
	sugar.Infof("failed to fetch URL %s", "url")
}
