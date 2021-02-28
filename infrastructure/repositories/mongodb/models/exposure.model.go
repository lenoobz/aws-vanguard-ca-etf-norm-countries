package models

import (
	"github.com/hthl85/aws-vanguard-ca-etf-norm-countries/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FundExposureModel is the representation of individual Vanguard fund overview model
type FundExposureModel struct {
	ID              *primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Schema          string              `json:"schema,omitempty" bson:"schema,omitempty"`
	IsActive        bool                `json:"isActive,omitempty" bson:"isActive,omitempty"`
	CreatedAt       int64               `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	ModifiedAt      int64               `json:"modifiedAt,omitempty" bson:"modifiedAt,omitempty"`
	Ticker          string              `json:"ticker,omitempty" bson:"ticker,omitempty"`
	AssetClass      string              `json:"assetClass,omitempty" bson:"assetClass,omitempty"`
	CountryExposure []*ExposureModel    `json:"countryExposure,omitempty" bson:"countryExposure,omitempty"`
}

// ExposureModel is the representation of country the fund exposed
type ExposureModel struct {
	CountryCode     string  `json:"countryCode,omitempty" bson:"countryCode,omitempty"`
	CountryName     string  `json:"countryName,omitempty" bson:"countryName,omitempty"`
	HoldingStatCode string  `json:"holdingStatCode,omitempty" bson:"holdingStatCode,omitempty"`
	FundMktPercent  float64 `json:"fundMktPercent,omitempty" bson:"fundMktPercent,omitempty"`
}

// NewFundExposureModel create new fund exposure model
func NewFundExposureModel(fund *entities.Exposure) *FundExposureModel {
	var countryExposures []*ExposureModel

	for _, exposure := range fund.CountryExposure {
		countryExposures = append(countryExposures, &ExposureModel{
			CountryCode:     exposure.CountryCode,
			CountryName:     exposure.CountryName,
			HoldingStatCode: exposure.HoldingStatCode,
			FundMktPercent:  exposure.FundMktPercent,
		})
	}

	return &FundExposureModel{
		Ticker:          fund.Ticker,
		AssetClass:      fund.AssetClass,
		CountryExposure: countryExposures,
	}
}
