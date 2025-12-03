package service

import "futuremarket/repository"

type BlacklistService struct {
	Repo repository.BlacklistRepository
}

func (s BlacklistService) BlacklistToken(token string) error {
	return s.Repo.Add(token)
}

func (s BlacklistService) IsTokenBlacklisted(token string) (bool, error) {
	return s.Repo.Exists(token)
}
