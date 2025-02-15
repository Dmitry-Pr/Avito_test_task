package services

import "merch-shop/internal/app/repositories"

type MerchServiceInterface interface {
	GetAllMerch() ([]string, error)
}

type MerchService struct {
	repo repositories.MerchRepositoryInterface
}

func NewMerchService(repo repositories.MerchRepositoryInterface) MerchServiceInterface {
	return &MerchService{repo: repo}
}

func (s *MerchService) GetAllMerch() ([]string, error) {
	return s.repo.GetAll()
}
