package generator

import (
	api "alteroSmartTestTask/backend/services/MS_Generation/common/api"
	log_context "alteroSmartTestTask/common/log/context"
	"alteroSmartTestTask/common/syncgo"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"sync"
	"testing"
	"time"
)

var (
	testDeviceName      = "testDeviceName"
	testDeviceFrequency = int32(10)
	testDevice          = &api.Device{
		DeviceId: &api.DeviceId{
			Name: testDeviceName,
		},
		Frequency: testDeviceFrequency,
	}
)

type generatorTest struct {
	suite.Suite
	generator *generator
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
	deviceDataChan := g.generator.CreateDevice(g.ctx, testDevice)
	endTimerChan := time.NewTimer(time.Second + 10*time.Millisecond).C
	var counter int32
	syncgo.GoWG(g.wg, func() {
		for {
			select {
			case <-deviceDataChan:
				counter += 1
				fmt.Printf("+1\n")
			case <-endTimerChan:
				fmt.Println("cancel")
				g.cancel()
				return
			}
		}
	})
	g.wg.Wait()
	g.Assert().Equal(testDeviceFrequency, counter)
}

func TestGenerator(t *testing.T) {
	suite.Run(t, new(generatorTest))
}
