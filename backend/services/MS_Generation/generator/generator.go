package generator

import (
	api "alteroSmartTestTask/backend/services/MS_Generation/common/api"
	log_context "alteroSmartTestTask/common/log/context"
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math/rand"
	"sync"
	"time"
)

type generator struct {
	deviceList map[string]context.CancelFunc
	wg         sync.WaitGroup
}

func NewGenerator() *generator {
	return &generator{deviceList: make(map[string]context.CancelFunc)}
}

func (g *generator) CreateDevice(
	ctx context.Context,
	device *api.Device,
) <-chan *api.DeviceData {
	dataChan := make(chan *api.DeviceData)
	ticker := time.NewTicker(1000 / time.Duration(device.GetFrequency()) * time.Millisecond)
	done := make(chan struct{})
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		contextDevice, cancel := context.WithCancel(
			log_context.WithLogger(ctx,
				log_context.FromContext(ctx).
					WithField("device_name", device.GetDeviceId().GetName()),
			))
		g.deviceList[device.GetDeviceId().GetName()] = cancel
		for {
			select {
			case <-done:
				close(dataChan)
				return
			case <-ticker.C:
				dataChan <- &api.DeviceData{
					DeviceId:  device.DeviceId,
					Data:      rand.Float64(),
					Timestamp: timestamppb.Now(),
				}
			}
		}
	}()
	return dataChan
}

func (g *generator) runGeneratorWithContext(
	ctx context.Context,
	device *api.Device,
	dataChan chan<- *api.DeviceData,
) {
	ticker := time.NewTicker(1000 / time.Duration(device.GetFrequency()) * time.Millisecond)
	done := make(chan struct{})
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		for {
			select {
			case <-done:
				close(dataChan)
				return
			case <-ticker.C:
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

func (g *generator) Wait() {
	g.wg.Wait()
}
