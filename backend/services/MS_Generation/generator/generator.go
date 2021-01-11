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
	g.deviceListMutex.RLock()
	if _, ok := g.deviceList[device.GetDeviceId().GetName()]; ok {
		g.deviceListMutex.RUnlock()
		return nil, errors.Newf(
			"device %s already exists",
			device.GetDeviceId().GetName(),
		)
	}
	g.deviceListMutex.RUnlock()

	dataChan := make(chan *api.DeviceData)
	ctxCancel, cancel := context.WithCancel(
		log_context.WithLogger(ctx,
			log_context.FromContext(ctx).
				WithField("device_name", device.GetDeviceId().GetName()),
		))
	g.deviceListMutex.Lock()
	g.deviceList[device.GetDeviceId().GetName()] = cancel
	g.deviceListMutex.Unlock()
	g.runGeneratorWithContext(ctxCancel, device, dataChan)
	return dataChan, nil
}

func (g *generator) RemoveDevice(
	ctx context.Context, deviceId *api.DeviceId,
) error {
	g.deviceListMutex.RLock()
	cancel, ok := g.deviceList[deviceId.GetName()]
	g.deviceListMutex.RUnlock()
	if !ok {
		return errors.New("device not found")
	}
	cancel()
	g.deviceListMutex.Lock()
	g.deviceList[deviceId.GetName()] = nil
	g.deviceListMutex.Unlock()
	return nil
}

func (g *generator) GetDeviceList(ctx context.Context) []string {
	var output []string
	g.deviceListMutex.RLock()
	for s := range g.deviceList {
		output = append(output, s)
	}
	g.deviceListMutex.RUnlock()
	return output
}

func (g *generator) Wait() {
	g.wg.Wait()
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
