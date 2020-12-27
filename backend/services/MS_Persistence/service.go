package MS_Persistence

import (
	"context"

	api "alteroSmartTestTask/backend/services/MS_Persistence/common/api"
	database "alteroSmartTestTask/backend/services/MS_Persistence/database"
)

type DatabaseClient interface {
	GetDeviceId(ctx context.Context, req *database.GetDeviceIdReq) (int64, error)
	InsertDeviceData(ctx context.Context, req *database.InsertDeviceDataReq) error
}

type service struct {
	DatabaseClient
}

func NewService(client DatabaseClient) *service {
	return &service{DatabaseClient: client}
}

func (s *service) SaveData(
	ctx context.Context,
	request *api.SaveDataRequest,
) (*api.SaveDataResponse, error) {
	id, err := s.DatabaseClient.GetDeviceId(ctx,
		&database.GetDeviceIdReq{
			Name: request.GetDeviceData().GetDeviceId().GetName(),
		})
	if err != nil {
		return nil, err
	}
	err = s.DatabaseClient.InsertDeviceData(ctx,
		&database.InsertDeviceDataReq{
			DeviceId:         id,
			Data:             request.GetDeviceData().GetData(),
			TimestampSeconds: request.GetDeviceData().GetTimestamp().GetSeconds(),
			TimestampNanos:   request.GetDeviceData().GetTimestamp().GetNanos(),
		})
	if err != nil {
		return nil, err
	}
	return &api.SaveDataResponse{}, nil
}
