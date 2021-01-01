package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FundOverviewModel is the representation of individual Vanguard fund overview model
type FundOverviewModel struct {
	ID              *primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
	Schema          int                     `json:"schema,omitempty" bson:"schema,omitempty"`
	IsActive        bool                    `json:"isActive,omitempty" bson:"isActive,omitempty"`
	CreatedAt       int64                   `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	ModifiedAt      int64                   `json:"modifiedAt,omitempty" bson:"modifiedAt,omitempty"`
	Ticker          string                  `json:"ticker,omitempty" bson:"ticker,omitempty"`
	CountryExposure []*CountryExposureModel `json:"countryExposure,omitempty" bson:"countryExposure,omitempty"`
}

// CountryExposureModel is the representation of country the fund exposed
type CountryExposureModel struct {
	CountryCode     string  `json:"countryCode,omitempty" bson:"countryCode,omitempty"`
	CountryName     string  `json:"countryName,omitempty" bson:"countryName,omitempty"`
	HoldingStatCode string  `json:"holdingStatCode,omitempty" bson:"holdingStatCode,omitempty"`
	FundMktPercent  float64 `json:"fundMktPercent,omitempty" bson:"fundMktPercent,omitempty"`
}
