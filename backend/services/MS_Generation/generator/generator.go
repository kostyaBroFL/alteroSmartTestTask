// Package generator store entity
// that emulates data generation by multiple devices.
package generator

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	api "alteroSmartTestTask/backend/services/MS_Generation/common/api"
	"alteroSmartTestTask/common/errors"
	logcontext "alteroSmartTestTask/common/log/context"
	"alteroSmartTestTask/common/syncgo"
)

// Generator is the entity that emulates data generation by multiple devices.
type generator struct {
	deviceListMutex *sync.RWMutex
	deviceList      map[string]context.CancelFunc
	wg              *sync.WaitGroup
}

// NewGenerator is the constructor for the generator.
func NewGenerator() *generator {
	return &generator{
		deviceListMutex: &sync.RWMutex{},
		deviceList:      make(map[string]context.CancelFunc),
		wg:              &sync.WaitGroup{},
	}
}

// CreateDevice is the method for creating and running device emulator.
func (g *generator) CreateDevice(
	ctx context.Context,
	device *api.Device,
) (<-chan *api.DeviceData, error) {
	ctx = logcontext.WithLogger(ctx, logcontext.FromContext(ctx).
		WithField("service_method", "CreateDevice").
		WithField("device_name", device.GetDeviceId().GetName()).
		WithField("device_frequency", device.GetFrequency()))
	g.deviceListMutex.RLock()
	if _, ok := g.deviceList[device.GetDeviceId().GetName()]; ok {
		g.deviceListMutex.RUnlock()
		err := errors.Newf("device %s already exists",
			device.GetDeviceId().GetName())
		logcontext.FromContext(ctx).WithError(err).
			Error("can not create device")
		return nil, err
	}
	g.deviceListMutex.RUnlock()

	dataChan := make(chan *api.DeviceData)
	deviceContext := logcontext.WithLogger(
		context.Background(),
		logcontext.FromContext(ctx))
	deviceContext, cancel := context.WithCancel(
		logcontext.WithLogger(deviceContext,
			logcontext.FromContext(deviceContext).
				WithField("device_name", device.GetDeviceId().GetName()),
		))
	g.deviceListMutex.Lock()
	g.deviceList[device.GetDeviceId().GetName()] = cancel
	g.deviceListMutex.Unlock()
	g.runGeneratorWithContext(deviceContext, device, dataChan)
	logcontext.FromContext(ctx).Info("generator started")
	return dataChan, nil
}

// RemoveDevice is the method for turn off and remove device by name.
func (g *generator) RemoveDevice(
	ctx context.Context, deviceId *api.DeviceId,
) error {
	ctx = logcontext.WithLogger(ctx, logcontext.FromContext(ctx).
		WithField("service_method", "RemoveDevice").
		WithField("device_name", deviceId.GetName()))
	g.deviceListMutex.RLock()
	cancel, ok := g.deviceList[deviceId.GetName()]
	g.deviceListMutex.RUnlock()
	if !ok {
		logcontext.FromContext(ctx).Error("device not found")
		return errors.New("device not found")
	}
	cancel()
	g.deviceListMutex.Lock()
	delete(g.deviceList, deviceId.GetName())
	g.deviceListMutex.Unlock()
	logcontext.FromContext(ctx).Info("success")
	return nil
}

// GetDeviceList is the method for requesting list of devices names.
func (g *generator) GetDeviceList(ctx context.Context) []string {
	ctx = logcontext.WithLogger(ctx, logcontext.FromContext(ctx).
		WithField("service_method", "GetDeviceList"))
	var output []string
	g.deviceListMutex.RLock()
	for s := range g.deviceList {
		output = append(output, s)
	}
	g.deviceListMutex.RUnlock()
	logcontext.FromContext(ctx).Info("success")
	return output
}

// Wait is the method for wait while all goroutines will be ended.
func (g *generator) Wait() {
	g.wg.Wait()
}

func (g *generator) runGeneratorWithContext(
	ctx context.Context,
	device *api.Device,
	dataChan chan<- *api.DeviceData,
) {
	// TODO[#6]: calculating is wrong.
	ticker := time.NewTicker(1000 /
		time.Duration(device.GetFrequency()) *
		time.Millisecond)
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
