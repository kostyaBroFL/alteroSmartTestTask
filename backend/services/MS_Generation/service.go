package MS_Generation

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	api "alteroSmartTestTask/backend/services/MS_Generation/common/api"

	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type service struct {
	deviceList map[string]context.CancelFunc
	wg         sync.WaitGroup
}

func NewService() *service {
	return &service{deviceList: make(map[string]context.CancelFunc)}
}

func (s *service) Wait() {
	s.wg.Wait()
}

func (s *service) AddDevice(
	ctx context.Context,
	request *api.AddDeviceRequest,
) (*api.AddDeviceResponse, error) {
	if _, ok := s.deviceList[request.GetDevice().GetDeviceId().GetName()]; ok {
		return &api.AddDeviceResponse{
			Status: api.AddDeviceStatus_DEVICE_IS_EXISTS,
		}, nil
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		fmt.Printf("стартую канал для устройства %+v\n", request.GetDevice())
		contextDevice, cancel := context.WithCancel(ctx)
		s.deviceList[request.GetDevice().GetDeviceId().GetName()] = cancel
		dataChan := make(chan *api.DeviceData)
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			for {
				data, ok := <-dataChan
				if ok {
					fmt.Printf("получил данные %+v\n\t от %+v\n", data, request.GetDevice())
				} else {
					fmt.Printf("закрываю gorutines от %+v\n", request.GetDevice())
					return
				}
			}
		}()
		s.runGeneratorWithContext(contextDevice, request.GetDevice(), dataChan)
	}()

	return &api.AddDeviceResponse{
		Status: api.AddDeviceStatus_SUCCESS,
	}, nil
}

func (s *service) RemoveDevice(
	ctx context.Context,
	request *api.RemoveDeviceRequest,
) (*api.RemoveDeviceResponse, error) {
	fmt.Printf("запускаю отключение %+v\n", request.DeviceId)
	cancel, ok := s.deviceList[request.GetDeviceId().GetName()]
	if !ok {
		fmt.Printf("устройство %+v не найдено\n", request.DeviceId)
		return nil, errors.New("not found")
	}
	cancel()
	return &api.RemoveDeviceResponse{}, nil
}

func (s *service) runGeneratorWithContext(
	ctx context.Context,
	device *api.Device,
	dataChan chan<- *api.DeviceData,
) {

	ticker := time.NewTicker(1000 / time.Duration(device.GetFrequency()) * time.Millisecond)
	done := make(chan struct{})

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			select {
			case <-done:
				fmt.Println("Выключаю генератор")
				close(dataChan)
				return
			case <-ticker.C:
				fmt.Println("Генерирую")
				dataChan <- &api.DeviceData{
					DeviceId:  device.DeviceId,
					Data:      rand.Float64(),
					Timestamp: timestamppb.Now(),
				}
			}
		}
	}()

	select {
	case <-ctx.Done():
		done <- struct{}{}
	}
}
