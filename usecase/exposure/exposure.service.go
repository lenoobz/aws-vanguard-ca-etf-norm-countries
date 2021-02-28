package exposure

import (
	"context"

	"github.com/hthl85/aws-vanguard-ca-etf-norm-countries/usecase/logger"
)

// Service exposure
type Service struct {
	repo IExposureRepo
	log  logger.IAppLogger
}

// NewExposureService create new service
func NewExposureService(r IExposureRepo, l logger.IAppLogger) *Service {
	return &Service{
		repo: r,
		log:  l,
	}
}

// PopulateFundExposures populate fund exposures data
func (s *Service) PopulateFundExposures(ctx context.Context) error {
	s.log.Info(ctx, "populate fund country exposure")

	exposures, err := s.repo.FindAllExposure(ctx)
	if err != nil {
		s.log.Error(ctx, "find all exposures failed")
		return err
	}

	if err := s.repo.UpdateAllExposure(ctx, exposures); err != nil {
		s.log.Error(ctx, "update all exposures failed")
		return err
	}

	return nil
}
