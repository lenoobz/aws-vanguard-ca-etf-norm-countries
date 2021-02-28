package entities

// Exposure represents Vanguard's fund Exposure
type Exposure struct {
	Ticker          string             `json:"ticker,omitempty"`
	AssetClass      string             `json:"assetClass,omitempty"`
	CountryExposure []*CountryExposure `json:"countryExposure,omitempty"`
}

// CountryExposure represents fund country exposure info
type CountryExposure struct {
	CountryCode     string  `json:"countryCode,omitempty"`
	CountryName     string  `json:"countryName,omitempty"`
	HoldingStatCode string  `json:"holdingStatCode,omitempty"`
	FundMktPercent  float64 `json:"fundMktPercent,omitempty"`
}

// Country represents country info
type Country struct {
	Name string `json:"name"`
	Code string `json:"alpha3Code"`
}
