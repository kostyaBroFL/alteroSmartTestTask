package server

import (
	"context"
	"google.golang.org/grpc"
	"sync"

	api "alteroSmartTestTask/backend/services/MS_Generation/common/api"
	persisapi "alteroSmartTestTask/backend/services/MS_Persistence/common/api"
	"alteroSmartTestTask/common/errors"
	log_context "alteroSmartTestTask/common/log/context"
	"alteroSmartTestTask/common/syncgo"
)

type Generator interface {
	CreateDevice(ctx context.Context, device *api.Device) (
		<-chan *api.DeviceData, error)
	RemoveDevice(ctx context.Context, deviceId *api.DeviceId) error
	GetDeviceList(ctx context.Context) []string
}

type MsPersistenceClient interface {
	SaveData(ctx context.Context, in *persisapi.SaveDataRequest,
		opts ...grpc.CallOption) (*persisapi.SaveDataResponse, error)
}

type protoServer struct {
	wg        *sync.WaitGroup
	Generator // todo interface
	MsPersistenceClient
}

func NewProtoServer(
	generator Generator, client MsPersistenceClient,
) *protoServer {
	return &protoServer{
		wg:                  &sync.WaitGroup{},
		Generator:           generator,
		MsPersistenceClient: client,
	}
}

func (p *protoServer) AddDevice(
	ctx context.Context,
	request *api.AddDeviceRequest,
) (*api.AddDeviceResponse, error) {
	ctx = log_context.WithLogger(ctx,
		log_context.FromContext(ctx).
			WithField("server_method", "AddDevice").
			WithField("device_name", request.GetDevice().GetDeviceId().GetName()).
			WithField("device_freq", request.GetDevice().GetFrequency()))
	dataChan, err := p.Generator.CreateDevice(ctx, request.GetDevice())
	if err != nil {
		log_context.FromContext(ctx).WithError(err).Error("can not create device")
		return nil, errors.Newf("can not create device: %s", err.Error())
	}
	syncgo.GoWG(p.wg, func() {
		for {
			data, ok := <-dataChan
			if !ok {
				return
			}
			_, err := p.MsPersistenceClient.SaveData(
				log_context.WithLogger(
					context.Background(),
					log_context.FromContext(ctx),
				), &persisapi.SaveDataRequest{
					DeviceData: generationDeviceDataToPersistence(data),
				})
			if err != nil {
				log_context.FromContext(ctx).
					WithError(err).
					Error("send data to ms persistence client error, turn off generator")
				err = p.Generator.RemoveDevice(ctx, request.GetDevice().GetDeviceId())
				if err != nil {
					log_context.FromContext(ctx).
						WithError(err).
						Error("cen not remove device")
				}
				return
			}
		}
	})

	var output []*api.DeviceId
	deviceNameList := p.Generator.GetDeviceList(ctx)
	for _, deviceName := range deviceNameList {
		output = append(output, &api.DeviceId{Name: deviceName})
	}
	return &api.AddDeviceResponse{
		ResultedDeviceList: output,
	}, nil
}

func (p *protoServer) RemoveDevice(
	ctx context.Context,
	request *api.RemoveDeviceRequest,
) (*api.RemoveDeviceResponse, error) {
	ctx = log_context.WithLogger(ctx,
		log_context.FromContext(ctx).
			WithField("device_name", request.GetDeviceId().GetName()))
	err := p.Generator.RemoveDevice(ctx, request.GetDeviceId())
	if err != nil {
		log_context.FromContext(ctx).
			WithError(err).
			Error("cen not remove device")
		return nil, err
	}
	var output []*api.DeviceId
	deviceNameList := p.Generator.GetDeviceList(ctx)
	for _, deviceName := range deviceNameList {
		output = append(output, &api.DeviceId{Name: deviceName})
	}
	return &api.RemoveDeviceResponse{
		ResultedDeviceList: output,
	}, nil
}

func (p *protoServer) Wait() {
	p.wg.Wait()
}

func persistenceDeviceDataToGeneration(
	data *persisapi.DeviceData,
) *api.DeviceData {
	return &api.DeviceData{
		DeviceId: &api.DeviceId{
			Name: data.GetDeviceId().GetName(),
		},
		Data:      data.GetData(),
		Timestamp: data.GetTimestamp(),
	}
}

func generationDeviceDataToPersistence(
	data *api.DeviceData,
) *persisapi.DeviceData {
	return &persisapi.DeviceData{
		DeviceId: &persisapi.DeviceId{
			Name: data.GetDeviceId().GetName(),
		},
		Data:      data.GetData(),
		Timestamp: data.GetTimestamp(),
	}
}
