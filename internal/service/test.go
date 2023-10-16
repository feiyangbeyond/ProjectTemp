package service

import (
	"context"

	"template/internal/model"
)

type TestRepo interface {
	Save(ctx context.Context, g *model.Test) (*model.Test, error)
}

type TestService struct {
	repo TestRepo
}

func NewTestService(repo TestRepo) *TestService {
	return &TestService{repo: repo}
}

func (t *TestService) TestXxx(ctx context.Context) error {
	if _, err := t.repo.Save(ctx, &model.Test{}); err != nil {
		return err
	}
	return nil
}
