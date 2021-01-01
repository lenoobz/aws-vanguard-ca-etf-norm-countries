package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/hthl85/aws-vanguard-ca-etf-countries/repositories/mongodb/repos"
	"github.com/hthl85/aws-vanguard-ca-etf-countries/services"
	"go.mongodb.org/mongo-driver/mongo"
)

var db *mongo.Database

func main() {
	repo, err := repos.NewFundRepo(db)
	if err != nil {
		fmt.Println("error occurred when connect to database", err)
	}

	// we won't close database connection
	db = repo.DB

	// init service
	svc := services.NewFundService(repo)

	// if err = svc.PopulateFundCountries(); err != nil {
	// 	fmt.Println("error populate fund countries")
	// }
	lambda.Start(svc.PopulateFundCountries)
}
