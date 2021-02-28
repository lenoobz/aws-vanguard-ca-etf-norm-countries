package exposure

import (
	"context"

	"github.com/hthl85/aws-vanguard-ca-etf-norm-countries/entities"
)

///////////////////////////////////////////////////////////
// Exposure Repository Interface
///////////////////////////////////////////////////////////

// IExposureReader interface
type IExposureReader interface {
	FindAllExposure(ctx context.Context) ([]*entities.Exposure, error)
}

// IExposureWriter interface
type IExposureWriter interface {
	UpdateAllExposure(ctx context.Context, exposure []*entities.Exposure) error
}

// IExposureRepo interface
type IExposureRepo interface {
	IExposureReader
	IExposureWriter
}

///////////////////////////////////////////////////////////
// Exposure Service Interface
///////////////////////////////////////////////////////////

// IExposureService define business rule of sector
type IExposureService interface {
	PopulateFundExposures() error
}
