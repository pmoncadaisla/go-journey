package service

import (
	"testing"

	"github.com/pmoncadaisla/go-journey/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func (suite *ServiceTestSuite) TestGetLast() {
	assert.Equal(suite.T(), 1, suite.service.GetLast().ID)

}

func (suite *ServiceTestSuite) TestGetNextID() {
	assert.Equal(suite.T(), 2, suite.service.GetNextID())
}

func (suite *ServiceTestSuite) SetupTest() {
	suite.service = Instance()
	suite.service.SetLast(&domain.Journey{ID: 1})

}

type ServiceTestSuite struct {
	suite.Suite
	service *Service
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
