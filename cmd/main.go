package main

import (
	"context"
	"log"

	"github.com/google/uuid"
	corid "github.com/lenoobz/aws-lambda-corid"
	logger "github.com/lenoobz/aws-lambda-logger"
	"github.com/lenoobz/aws-vanguard-ca-etf-norm-countries/config"
	"github.com/lenoobz/aws-vanguard-ca-etf-norm-countries/infrastructure/repositories/mongodb/repos"
	"github.com/lenoobz/aws-vanguard-ca-etf-norm-countries/usecase/breakdown"
)

func main() {
	appConf := config.AppConf

	// create new logger
	zap, err := logger.NewZapLogger()
	if err != nil {
		log.Fatal("create app logger failed")
	}
	defer zap.Close()

	// create new repository
	repo, err := repos.NewBreakdownMongo(nil, zap, &appConf.Mongo)
	if err != nil {
		log.Fatal("create fund mongo repo failed")
	}
	defer repo.Close()

	// create new service
	svc := breakdown.NewService(repo, zap)

	// try correlation context
	id, _ := uuid.NewRandom()
	ctx := corid.NewContext(context.Background(), id)
	svc.PopulateFundBreakdown(ctx)
}
