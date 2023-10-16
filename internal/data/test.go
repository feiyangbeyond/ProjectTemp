package data

import (
	"context"

	"template/internal/model"
	"template/internal/service"
	"template/pkg/log"
)

var _ service.TestRepo = (*testRepo)(nil)

type testRepo struct {
	data *Data
	log  *log.Helper
}

func NewTestRepo(data *Data, logger log.Logger) service.TestRepo {
	return &testRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r testRepo) Save(ctx context.Context, g *model.Test) (*model.Test, error) {
	return nil, nil
}
