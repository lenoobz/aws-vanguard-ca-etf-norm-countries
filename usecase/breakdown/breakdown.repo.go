package breakdown

import (
	"context"

	"github.com/hthl85/aws-vanguard-ca-etf-norm-countries/entities"
)

///////////////////////////////////////////////////////////
// Exposure Repository Interface
///////////////////////////////////////////////////////////

// Reader interface
type Reader interface {
	FindCountriesBreakdown(context.Context) ([]*entities.FundBreakdown, error)
}

// Writer interface
type Writer interface {
	UpdateCountriesBreakdown(context.Context, []*entities.FundBreakdown) error
}

// Repo interface
type Repo interface {
	Reader
	Writer
}
