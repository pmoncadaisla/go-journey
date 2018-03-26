package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func (suite *ServiceTestSuite) TestGaugeInc() {
	suite.service.GaugeInc("test_gauge_inc")
	assert.Equal(suite.T(), 1, suite.service.Gauges["test_gauge_inc"].Value)
	suite.service.GaugeInc("test_gauge_inc")
	assert.Equal(suite.T(), 2, suite.service.Gauges["test_gauge_inc"].Value)
	suite.service.GaugeInc("test_gauge_inc_other")
	assert.Equal(suite.T(), 2, suite.service.Gauges["test_gauge_inc"].Value)
	assert.Equal(suite.T(), 1, suite.service.Gauges["test_gauge_inc_other"].Value)
}

func (suite *ServiceTestSuite) TestGaugeDec() {
	suite.service.GaugeInc("test_gauge_dec")
	suite.service.GaugeInc("test_gauge_dec")
	assert.Equal(suite.T(), 2, suite.service.Gauges["test_gauge_dec"].Value)
	suite.service.GaugeDec("test_gauge_dec")
	assert.Equal(suite.T(), 1, suite.service.Gauges["test_gauge_dec"].Value)
	suite.service.GaugeDec("test_gauge_dec")
	assert.Equal(suite.T(), 0, suite.service.Gauges["test_gauge_dec"].Value)
	suite.service.GaugeDec("test_gauge_dec")
	suite.service.GaugeDec("test_gauge_dec")
	assert.Equal(suite.T(), 0, suite.service.Gauges["test_gauge_dec"].Value)
}

func (suite *ServiceTestSuite) TestCounterInc() {
	suite.service.CounterInc("test_counter_inc")
	assert.Equal(suite.T(), 1, suite.service.Counters["test_counter_inc"].Value)
	suite.service.CounterInc("test_counter_inc")
	suite.service.CounterInc("test_counter_inc")
	suite.service.CounterInc("test_counter_inc")
	assert.Equal(suite.T(), 4, suite.service.Counters["test_counter_inc"].Value)

}

func (suite *ServiceTestSuite) SetupTest() {
	suite.service = Instance()

}

type ServiceTestSuite struct {
	suite.Suite
	service *Service
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
