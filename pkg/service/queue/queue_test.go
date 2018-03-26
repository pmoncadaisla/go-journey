package service

import (
	"testing"

	"github.com/pmoncadaisla/go-journey/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func (suite *ServiceTestSuite) TestLen() {
	assert.Equal(suite.T(), 0, suite.queue.Len())
	suite.queue.Push(domain.Journey{})
	assert.Equal(suite.T(), 1, suite.queue.Len())
	suite.queue.Get()
	assert.Equal(suite.T(), 1, suite.queue.Len())
	suite.queue.Pop()
	assert.Equal(suite.T(), 0, suite.queue.Len())
}

func (suite *ServiceTestSuite) TestPush() {
	suite.queue.Push(domain.Journey{ID: 1})
	assert.Equal(suite.T(), 1, suite.queue.Len())
	suite.queue.Push(domain.Journey{ID: 2})
	assert.Equal(suite.T(), 2, suite.queue.Len())

}

func (suite *ServiceTestSuite) TestPop() {
	suite.queue.Push(domain.Journey{ID: 1})
	suite.queue.Push(domain.Journey{ID: 2})
	assert.Equal(suite.T(), 2, suite.queue.Len())
	assert.Equal(suite.T(), 1, suite.queue.Pop().ID)
	assert.Equal(suite.T(), 2, suite.queue.Pop().ID)

}

func (suite *ServiceTestSuite) TestGet() {
	suite.queue.Push(domain.Journey{ID: 1})
	suite.queue.Push(domain.Journey{ID: 2})
	assert.Equal(suite.T(), 2, suite.queue.Len())
	assert.Equal(suite.T(), 1, suite.queue.Get().ID)
	assert.Equal(suite.T(), 1, suite.queue.Get().ID)
	assert.Equal(suite.T(), 2, suite.queue.Len())
	suite.queue.Pop()
	suite.queue.Pop()
}

func (suite *ServiceTestSuite) SetupTest() {
	suite.queue = Instance()

}

type ServiceTestSuite struct {
	suite.Suite
	queue Interface
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
