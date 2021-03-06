package mem

import (
	"testing"
	"github.com/stretchr/testify/suite"
	"github.com/lexlapax/graveldb/core"
	"github.com/lexlapax/graveldb/coretest"
)

func init() {
	Register()
}

type GraphTestSuite struct {
	coretest.GraphTestSuite
}

func (suite *GraphTestSuite) SetupSuite() {
    suite.TestGraph = core.GetGraph(GraphImpl)
}

func (suite *GraphTestSuite) TearSuite() {
    suite.TestGraph = nil
}

func TestGraphTestSuite(t *testing.T) {
	//t.Skip()
    suite.Run(t, new(GraphTestSuite))
}
