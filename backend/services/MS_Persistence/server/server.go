package server

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"

	api "alteroSmartTestTask/backend/services/MS_Persistence/common/api"
	"alteroSmartTestTask/backend/services/MS_Persistence/database"
	"alteroSmartTestTask/common/errors"
	logcontext "alteroSmartTestTask/common/log/context"
)

// TODO[#4]: Use MQ.
// TODO[#7]: Database replication.
var timeFormat = "2006-01-02T15:04:05.999999999Z"

// DatabaseClient is the interface with method of database.
type DatabaseClient interface {
	GetDeviceId(ctx context.Context, req *database.GetDeviceIdReq) (int64, error)
	InsertDeviceData(ctx context.Context, req *database.InsertDeviceDataReq) error
	GetDataByDeviceName(ctx context.Context,
		req *database.GetDataByDeviceNameRequest) ([]*database.DeviceData, error)
}

// This is the implementation of the MsPersistence service.
type service struct {
	DatabaseClient
}

// NewService is the method for creating MsPersistence service.
func NewService(client DatabaseClient) *service {
	return &service{DatabaseClient: client}
}

// SaveData is the method for saving chunk of the device's data.
func (s *service) SaveData(
	ctx context.Context,
	request *api.SaveDataRequest,
) (*api.SaveDataResponse, error) {
	ctx = logcontext.WithLogger(ctx, logcontext.FromContext(ctx).
		WithField("service_method", "SaveData").
		WithField("device_name", request.GetDeviceData().GetDeviceId().GetName()))
	id, err := s.DatabaseClient.GetDeviceId(ctx,
		&database.GetDeviceIdReq{
			Name: request.GetDeviceData().GetDeviceId().GetName(),
		})
	if err != nil {
		return nil, errors.ToFrontendError(ctx, err,
			codes.Internal, "Get device id error.")
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
		return nil, errors.ToFrontendError(ctx, err,
			codes.Internal, "Insert device data error.")
	}
	logcontext.FromContext(ctx).Info("success")
	return &api.SaveDataResponse{}, nil
}

// GetData is the method for loading
// any last chunks of the data of specific device.
func (s *service) GetData(
	ctx context.Context,
	request *api.GetDataRequest,
) (*api.GetDataResponse, error) {
	ctx = logcontext.WithLogger(ctx, logcontext.FromContext(ctx).
		WithField("service_method", "GetData").
		WithField("limit", request.GetLimit()).
		WithField("device_name", request.GetDeviceId().GetName()))
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
	logcontext.FromContext(ctx).Info("success")
	return &api.GetDataResponse{DeviceData: apiData}, nil
}
