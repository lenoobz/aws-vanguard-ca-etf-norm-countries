package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/hthl85/aws-vanguard-ca-etf-norm-countries/config"
	"github.com/hthl85/aws-vanguard-ca-etf-norm-countries/infrastructure/logger"
	"github.com/hthl85/aws-vanguard-ca-etf-norm-countries/infrastructure/repositories/mongodb/repo"
	"github.com/hthl85/aws-vanguard-ca-etf-norm-countries/usecase/exposure"
)

func main() {
	appConf := config.AppConf

	// create new logger
	logger, err := logger.NewZapLogger()
	if err != nil {
		log.Fatal("create app logger failed")
	}
	defer logger.Close()

	// create new repository
	repo, err := repo.NewExposureMongo(nil, logger, &appConf.Mongo)
	if err != nil {
		log.Fatal("create fund mongo repo failed")
	}
	defer repo.Close()

	// create new service
	svc := exposure.NewExposureService(repo, logger)

	lambda.Start(svc.PopulateFundExposures)
}
