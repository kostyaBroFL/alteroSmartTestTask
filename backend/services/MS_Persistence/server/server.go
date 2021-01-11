package server

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"

	api "alteroSmartTestTask/backend/services/MS_Persistence/common/api"
	database "alteroSmartTestTask/backend/services/MS_Persistence/database"
	"alteroSmartTestTask/common/errors"
	logcontext "alteroSmartTestTask/common/log/context"
)

var timeFormat = "2006-01-02T15:04:05.999999999Z"

type DatabaseClient interface {
	GetDeviceId(ctx context.Context, req *database.GetDeviceIdReq) (int64, error)
	InsertDeviceData(ctx context.Context, req *database.InsertDeviceDataReq) error
	GetDataByDeviceName(ctx context.Context,
		req *database.GetDataByDeviceNameRequest) ([]*database.DeviceData, error)
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
	ctx = logcontext.WithLogger(ctx, logcontext.FromContext(ctx).
		WithField("device_name", request.GetDeviceData().GetDeviceId().GetName()))
	logcontext.FromContext(ctx).Info("Save data request start")
	id, err := s.DatabaseClient.GetDeviceId(ctx,
		&database.GetDeviceIdReq{
			Name: request.GetDeviceData().GetDeviceId().GetName(),
		})
	if err != nil {
		return nil, errors.ToFrontendError(ctx, err, codes.Internal, "Get device id error.")
	}
	err = s.DatabaseClient.InsertDeviceData(ctx,
		&database.InsertDeviceDataReq{
			DeviceId: id,
			Data:     request.GetDeviceData().GetData(),
			Timestamp: time.Unix(
				request.GetDeviceData().GetTimestamp().GetSeconds(),
				int64(request.GetDeviceData().GetTimestamp().GetNanos()),
			).Format(timeFormat),
		})
	if err != nil {
		return nil, errors.ToFrontendError(ctx, err, codes.Internal, "Insert device data error.")
	}
	return &api.SaveDataResponse{}, nil
}

func (s *service) GetData(
	ctx context.Context,
	request *api.GetDataRequest,
) (*api.GetDataResponse, error) {
	dataList, err := s.DatabaseClient.GetDataByDeviceName(
		ctx, &database.GetDataByDeviceNameRequest{
			DeviceName: request.GetDeviceId().GetName(),
			Limit:      request.GetLimit(),
		})
	if err != nil {
		return nil, errors.ToFrontendError(ctx, err, codes.Internal, "Database error")
	}
	var apiData []*api.DeviceData
	for _, data := range dataList {
		dataTime, err := time.Parse(timeFormat, data.Timestamp)
		if err != nil {
			return nil, errors.ToFrontendError(ctx, err, codes.Internal, "Data error")
		}
		apiData = append(apiData, &api.DeviceData{
			DeviceId:  request.DeviceId,
			Data:      data.Data,
			Timestamp: timestamppb.New(dataTime),
		})
	}
	return &api.GetDataResponse{DeviceData: apiData}, nil
}
