package MS_Persistence

import (
	"context"

	api "alteroSmartTestTask/backend/services/MS_Persistence/common/api"
)

type service struct {
}

func (s service) SaveData(
	ctx context.Context, request *api.SaveDataRequest) (*api.SaveDataResponse, error) {
	panic("implement me")
}
