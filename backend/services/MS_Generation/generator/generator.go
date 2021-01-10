package generator

import (
	api "alteroSmartTestTask/backend/services/MS_Generation/common/api"
	log_context "alteroSmartTestTask/common/log/context"
	"alteroSmartTestTask/common/syncgo"
	"context"
	"errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math/rand"
	"sync"
	"time"
)

type generator struct {
	deviceListMutex *sync.RWMutex
	deviceList      map[string]context.CancelFunc
	wg              *sync.WaitGroup
}

func NewGenerator() *generator {
	return &generator{
		deviceListMutex: &sync.RWMutex{},
		deviceList:      make(map[string]context.CancelFunc),
		wg:              &sync.WaitGroup{},
	}
}

func (g *generator) CreateDevice(
	ctx context.Context,
	device *api.Device,
) (<-chan *api.DeviceData, error) {
	dataChan := make(chan *api.DeviceData)
	ctxCancel, cancel := context.WithCancel(
		log_context.WithLogger(ctx,
			log_context.FromContext(ctx).
				WithField("device_name", device.GetDeviceId().GetName()),
		))
	g.deviceList[device.GetDeviceId().GetName()] = cancel
	g.runGeneratorWithContext(ctxCancel, device, dataChan)
	return dataChan, nil
}

func (g *generator) runGeneratorWithContext(
	ctx context.Context,
	device *api.Device,
	dataChan chan<- *api.DeviceData,
) {
	// TODO bad calc
	ticker := time.NewTicker(1000 / time.Duration(device.GetFrequency()) * time.Millisecond)
	done := make(chan struct{})
	syncgo.GoWG(g.wg, func() {
		for {
			select {
			case <-done:
				close(dataChan)
				ticker.Stop()
				return
			case <-ticker.C:
				dataChan <- &api.DeviceData{
					DeviceId:  device.DeviceId,
					Data:      rand.Float64(),
					Timestamp: timestamppb.Now(),
				}
			}
		}
	})

	go func() {
		select {
		case <-ctx.Done():
			done <- struct{}{}
		}
	}()
}

func (g *generator) RemoveDevice(
	ctx context.Context, device *api.Device,
) error {
	g.deviceListMutex.Lock()
	cancel, ok := g.deviceList[device.GetDeviceId().GetName()]
	if !ok {
		return errors.New("device not found")
	}
	cancel()
	g.deviceList[device.GetDeviceId().GetName()] = nil
	g.deviceListMutex.Unlock()
	return nil
}

func (g *generator) Wait() {
	g.wg.Wait()
}
