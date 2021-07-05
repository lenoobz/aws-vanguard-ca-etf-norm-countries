package models

import (
	"context"
	"time"

	logger "github.com/hthl85/aws-lambda-logger"
	"github.com/hthl85/aws-vanguard-ca-etf-norm-countries/consts"
	"github.com/hthl85/aws-vanguard-ca-etf-norm-countries/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FundBreakdownModel struct
type FundBreakdownModel struct {
	ID         *primitive.ObjectID `bson:"_id,omitempty"`
	CreatedAt  int64               `bson:"createdAt,omitempty"`
	ModifiedAt int64               `bson:"modifiedAt,omitempty"`
	Enabled    bool                `bson:"enabled"`
	Deleted    bool                `bson:"deleted"`
	Schema     string              `bson:"schema,omitempty"`
	Source     string              `bson:"source,omitempty"`
	Ticker     string              `bson:"ticker,omitempty"`
	AssetClass string              `bson:"assetClass,omitempty"`
	Countries  []*BreakdownModel   `bson:"countries,omitempty"`
}

// BreakdownModel struct
type BreakdownModel struct {
	CountryCode     string  `bson:"countryCode,omitempty"`
	CountryName     string  `bson:"countryName,omitempty"`
	HoldingStatCode string  `bson:"holdingStatCode,omitempty"`
	FundMktPercent  float64 `bson:"fundMktPercent,omitempty"`
	FundTnaPercent  float64 `bson:"fundTnaPercent,omitempty"`
}

// NewFundBreakdownModel create new fund exposure model
func NewFundBreakdownModel(ctx context.Context, log logger.ContextLog, e *entities.FundBreakdown, schemaVersion string) *FundBreakdownModel {
	var breakdownModel []*BreakdownModel

	for _, country := range e.Countries {
		breakdownModel = append(breakdownModel, &BreakdownModel{
			CountryCode:     country.CountryCode,
			CountryName:     country.CountryName,
			HoldingStatCode: country.HoldingStatCode,
			FundMktPercent:  country.FundMktPercent,
			FundTnaPercent:  country.FundTnaPercent,
		})
	}

	return &FundBreakdownModel{
		ModifiedAt: time.Now().UTC().Unix(),
		Enabled:    true,
		Deleted:    false,
		Schema:     schemaVersion,
		Source:     consts.DATA_SOURCE,
		Ticker:     e.Ticker,
		AssetClass: e.AssetClass,
		Countries:  breakdownModel,
	}
}
