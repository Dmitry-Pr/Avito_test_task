package services

import "merch-store/internal/app/repositories"

type MerchService struct {
	repo *repositories.MerchRepository
}

func NewMerchService(repo *repositories.MerchRepository) *MerchService {
	return &MerchService{repo: repo}
}

func (s *MerchService) GetAllMerch() ([]string, error) {
	return s.repo.GetAll()
}
