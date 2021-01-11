package generator

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	api "alteroSmartTestTask/backend/services/MS_Generation/common/api"
	"alteroSmartTestTask/common/errors"
	log_context "alteroSmartTestTask/common/log/context"
	"alteroSmartTestTask/common/syncgo"
)

type Generator struct {
	deviceListMutex *sync.RWMutex
	deviceList      map[string]context.CancelFunc
	wg              *sync.WaitGroup
}

func NewGenerator() *Generator {
	return &Generator{
		deviceListMutex: &sync.RWMutex{},
		deviceList:      make(map[string]context.CancelFunc),
		wg:              &sync.WaitGroup{},
	}
}

func (g *Generator) CreateDevice(
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

func (g *Generator) runGeneratorWithContext(
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

func (g *Generator) RemoveDevice(
	ctx context.Context, deviceId *api.DeviceId,
) error {
	g.deviceListMutex.Lock()
	cancel, ok := g.deviceList[deviceId.GetName()]
	if !ok {
		return errors.New("device not found")
	}
	cancel()
	g.deviceList[deviceId.GetName()] = nil
	g.deviceListMutex.Unlock()
	return nil
}

func (g *Generator) Wait() {
	g.wg.Wait()
}
