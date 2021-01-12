package server

import (
	"context"
	"sync"

	"google.golang.org/grpc"

	api "alteroSmartTestTask/backend/services/MS_Generation/common/api"
	persisapi "alteroSmartTestTask/backend/services/MS_Persistence/common/api"
	"alteroSmartTestTask/common/errors"
	logcontext "alteroSmartTestTask/common/log/context"
	"alteroSmartTestTask/common/syncgo"
)

// TODO[#4]: Use MQ.

// Generator is the interface with methods for manipulating devices.
type Generator interface {
	CreateDevice(ctx context.Context, device *api.Device) (
		<-chan *api.DeviceData, error)
	RemoveDevice(ctx context.Context, deviceId *api.DeviceId) error
	GetDeviceList(ctx context.Context) []string
}

// MsPersistenceClient is the interface with methods wor save data from devices.
type MsPersistenceClient interface {
	SaveData(ctx context.Context, in *persisapi.SaveDataRequest,
		opts ...grpc.CallOption) (*persisapi.SaveDataResponse, error)
}

// This is the implementation of the MsGeneration service.
type server struct {
	wg *sync.WaitGroup
	Generator
	MsPersistenceClient
}

// NewServer is the constructor forMsGeneration service.
func NewServer(
	generator Generator, client MsPersistenceClient,
) *server {
	return &server{
		wg:                  &sync.WaitGroup{},
		Generator:           generator,
		MsPersistenceClient: client,
	}
}

// AddDevice is the method for creating and running device emulator.
func (p *server) AddDevice(
	ctx context.Context,
	request *api.AddDeviceRequest,
) (*api.AddDeviceResponse, error) {
	ctx = logcontext.WithLogger(ctx,
		logcontext.FromContext(ctx).
			WithField("server_method", "AddDevice").
			WithField("device_name", request.GetDevice().GetDeviceId().GetName()).
			WithField("device_freq", request.GetDevice().GetFrequency()))
	dataChan, err := p.Generator.CreateDevice(ctx, request.GetDevice())
	if err != nil {
		logcontext.FromContext(ctx).WithError(err).
			Error("can not create device")
		return nil, errors.Newf("can not create device: %s", err.Error())
	}
	syncgo.GoWG(p.wg, func() {
		for {
			data, ok := <-dataChan
			if !ok {
				return
			}
			_, err := p.MsPersistenceClient.SaveData(
				logcontext.WithLogger(
					context.Background(),
					logcontext.FromContext(ctx),
				), &persisapi.SaveDataRequest{
					DeviceData: generationDeviceDataToPersistence(data),
				})
			if err != nil {
				logcontext.FromContext(ctx).
					WithError(err).
					Error("send data to ms persistence client error, turn off generator")
				err = p.Generator.RemoveDevice(ctx, request.GetDevice().GetDeviceId())
				if err != nil {
					logcontext.FromContext(ctx).
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
	logcontext.FromContext(ctx).Info("success")
	return &api.AddDeviceResponse{
		ResultedDeviceList: output,
	}, nil
}

// RemoveDevice is the method for turn off and remove device.
func (p *server) RemoveDevice(
	ctx context.Context,
	request *api.RemoveDeviceRequest,
) (*api.RemoveDeviceResponse, error) {
	ctx = logcontext.WithLogger(ctx,
		logcontext.FromContext(ctx).
			WithField("service_method", "RemoveDevice").
			WithField("device_name", request.GetDeviceId().GetName()))
	err := p.Generator.RemoveDevice(ctx, request.GetDeviceId())
	if err != nil {
		logcontext.FromContext(ctx).
			WithError(err).
			Error("cen not to remove device")
		return nil, err
	}
	var output []*api.DeviceId
	deviceNameList := p.Generator.GetDeviceList(ctx)
	for _, deviceName := range deviceNameList {
		output = append(output, &api.DeviceId{Name: deviceName})
	}
	logcontext.FromContext(ctx).Info("success")
	return &api.RemoveDeviceResponse{
		ResultedDeviceList: output,
	}, nil
}

// Wait is the method for wait while all goroutines will be ended.
func (p *server) Wait() {
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
