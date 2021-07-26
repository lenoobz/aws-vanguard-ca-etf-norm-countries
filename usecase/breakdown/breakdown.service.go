package breakdown

import (
	"context"

	logger "github.com/lenoobz/aws-lambda-logger"
)

// Service exposure
type Service struct {
	repo Repo
	log  logger.ContextLog
}

// NewService create new service
func NewService(r Repo, l logger.ContextLog) *Service {
	return &Service{
		repo: r,
		log:  l,
	}
}

// PopulateFundBreakdown populate fund exposures data
func (s *Service) PopulateFundBreakdown(ctx context.Context) error {
	s.log.Info(ctx, "populate fund country exposure")

	exposures, err := s.repo.FindCountriesBreakdown(ctx)
	if err != nil {
		s.log.Error(ctx, "find all exposures failed", "error", err)
		return err
	}

	if err := s.repo.UpdateCountriesBreakdown(ctx, exposures); err != nil {
		s.log.Error(ctx, "update all exposures failed", "error", err)
		return err
	}

	return nil
}
