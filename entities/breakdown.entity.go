package entities

// FundBreakdown represents Vanguard's fund FundBreakdown
type FundBreakdown struct {
	Ticker     string              `json:"ticker,omitempty"`
	AssetClass string              `json:"assetClass,omitempty"`
	Countries  []*CountryBreakdown `json:"countries,omitempty"`
}

// CountryBreakdown represents fund country exposure info
type CountryBreakdown struct {
	CountryCode     string  `json:"countryCode,omitempty"`
	CountryName     string  `json:"countryName,omitempty"`
	HoldingStatCode string  `json:"holdingStatCode,omitempty"`
	FundMktPercent  float64 `json:"fundMktPercent,omitempty"`
	FundTnaPercent  float64 `json:"fundTnaPercent,omitempty"`
}
