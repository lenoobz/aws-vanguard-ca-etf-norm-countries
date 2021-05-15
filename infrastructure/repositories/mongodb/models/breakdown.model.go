package models

import (
	"github.com/hthl85/aws-vanguard-ca-etf-norm-countries/consts"
	"github.com/hthl85/aws-vanguard-ca-etf-norm-countries/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FundBreakdownModel is the representation of individual Vanguard fund overview model
type FundBreakdownModel struct {
	ID         *primitive.ObjectID `bson:"_id,omitempty"`
	IsActive   bool                `bson:"isActive,omitempty"`
	CreatedAt  int64               `bson:"createdAt,omitempty"`
	ModifiedAt int64               `bson:"modifiedAt,omitempty"`
	Schema     string              `bson:"schema,omitempty"`
	Source     string              `bson:"source,omitempty"`
	Ticker     string              `bson:"ticker,omitempty"`
	AssetClass string              `bson:"assetClass,omitempty"`
	Countries  []*BreakdownModel   `bson:"countries,omitempty"`
}

// BreakdownModel is the representation of country the fund exposed
type BreakdownModel struct {
	CountryCode     string  `bson:"countryCode,omitempty"`
	CountryName     string  `bson:"countryName,omitempty"`
	HoldingStatCode string  `bson:"holdingStatCode,omitempty"`
	FundMktPercent  float64 `bson:"fundMktPercent,omitempty"`
	FundTnaPercent  float64 `bson:"fundTnaPercent,omitempty"`
}

// NewFundBreakdownModel create new fund exposure model
func NewFundBreakdownModel(e *entities.FundBreakdown) *FundBreakdownModel {
	var m []*BreakdownModel

	for _, v := range e.Countries {
		m = append(m, &BreakdownModel{
			CountryCode:     v.CountryCode,
			CountryName:     v.CountryName,
			HoldingStatCode: v.HoldingStatCode,
			FundMktPercent:  v.FundMktPercent,
			FundTnaPercent:  v.FundTnaPercent,
		})
	}

	return &FundBreakdownModel{
		Source:     consts.DATA_SOURCE,
		Ticker:     e.Ticker,
		AssetClass: e.AssetClass,
		Countries:  m,
	}
}
