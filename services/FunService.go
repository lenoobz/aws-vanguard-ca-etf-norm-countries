package services

import (
	"fmt"

	"github.com/hthl85/aws-vanguard-ca-etf-countries/repositories"
)

// FundService struct
type FundService struct {
	fundRepo repositories.IFundRepository
}

// NewFundService create as new service
func NewFundService(fundRepo repositories.IFundRepository) *FundService {
	fmt.Println("Create new Fund Service")

	return &FundService{
		fundRepo: fundRepo,
	}
}

// PopulateFundCountries find fund country exposure
func (svc *FundService) PopulateFundCountries() error {
	fmt.Println("Populate Fund Countries")

	funds, err := svc.fundRepo.GetAllFundsOverview()
	if err != nil {
		return err
	}

	return svc.fundRepo.UpdateAllFundsOverview(funds)
}
