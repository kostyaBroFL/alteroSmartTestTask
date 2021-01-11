package generator

import (
	api "alteroSmartTestTask/backend/services/MS_Generation/common/api"
	log_context "alteroSmartTestTask/common/log/context"
	"alteroSmartTestTask/common/syncgo"
	"context"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"sync"
	"testing"
	"time"
)

var (
	testDeviceName      = "testDeviceName"
	testDeviceFrequency = int32(500)
	testDeviceId        = &api.DeviceId{
		Name: testDeviceName,
	}
	testDevice = &api.Device{
		DeviceId:  testDeviceId,
		Frequency: testDeviceFrequency,
	}
)

type generatorTest struct {
	suite.Suite
	generator *Generator
	ctx       context.Context
	cancel    func()
	wg        *sync.WaitGroup
}

func (g *generatorTest) SetupTest() {
	g.wg = &sync.WaitGroup{}
	g.generator = NewGenerator()
	cancelContext, cancelFunc := context.WithCancel(context.Background())
	g.ctx = log_context.WithLogger(
		cancelContext,
		logrus.NewEntry(logrus.New()),
	)
	g.cancel = cancelFunc
}

func (g *generatorTest) TestCreateDeviceFreq() {
	deviceDataChan, err := g.generator.CreateDevice(g.ctx, testDevice)
	g.Require().Nil(err)
	endTimerChan := time.NewTimer(251 * time.Millisecond).C
	var counter int32
	syncgo.GoWG(g.wg, func() {
		for {
			select {
			case <-deviceDataChan:
				counter += 1
			case <-endTimerChan:
				g.cancel()
				return
			}
		}
	})
	g.wg.Wait()
	g.Assert().Equal(testDeviceFrequency/4, counter)
	err = g.generator.RemoveDevice(g.ctx, testDeviceId)
	g.Require().Nil(err)
	g.generator.Wait()
}

func (g *generatorTest) TestCreateDeviceDuplicateError() {
	_, err := g.generator.CreateDevice(g.ctx, testDevice)
	g.Require().Nil(err)
	_, err = g.generator.CreateDevice(g.ctx, testDevice)
	g.Require().NotNil(err)
	g.Assert().Contains(err.Error(), "already exists")
	err = g.generator.RemoveDevice(g.ctx, testDeviceId)
	g.Require().Nil(err)
	g.generator.Wait()
}

func (g *generatorTest) TestRemoveDevice() {
	dataChan, err := g.generator.CreateDevice(g.ctx, testDevice)
	g.Require().Nil(err)
	err = g.generator.RemoveDevice(g.ctx, testDeviceId)
	g.Require().Nil(err)
	g.generator.Wait()
	_, ok := <-dataChan
	g.Require().False(ok)
}

func TestGenerator(t *testing.T) {
	suite.Run(t, new(generatorTest))
}
