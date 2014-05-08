package coretest

import (
	"github.com/stretchr/testify/assert"
	"github.com/lexlapax/graveldb/core"
	"github.com/lexlapax/graveldb/util"
)

type GraphKeyIndexTestSuite struct {
	BaseTestSuite
}

func (suite *GraphKeyIndexTestSuite) TestKeyIndex(){
	if suite.TestGraph.Capabilities().KeyIndex() == false {
		suite.T().Skip()
	}

	stringset := util.NewStringSet()
	atomset := core.NewAtomSet()

	assert.Equal(suite.T(), 0, len(suite.TestGraph.IndexedKeys(core.VertexType)))
	assert.Equal(suite.T(), 0, len(suite.TestGraph.IndexedKeys(core.EdgeType)))
	//create index
	suite.TestGraph.CreateKeyIndex("name", core.VertexType)
	suite.TestGraph.CreateKeyIndex("name", core.EdgeType)

	keynames := suite.TestGraph.IndexedKeys(core.VertexType)
	assert.Equal(suite.T(), 1, len(keynames))
	assert.Equal(suite.T(), "name", keynames[0])

	keynames = suite.TestGraph.IndexedKeys(core.EdgeType)
	assert.Equal(suite.T(), 1, len(keynames))
	assert.Equal(suite.T(), "name", keynames[0])

	suite.TestGraph.CreateKeyIndex("name", core.VertexType)
	keynames = suite.TestGraph.IndexedKeys(core.VertexType)
	assert.Equal(suite.T(), 1, len(keynames))
	assert.Equal(suite.T(), "name", keynames[0])

	suite.TestGraph.CreateKeyIndex("weight", core.EdgeType)
	keynames = suite.TestGraph.IndexedKeys(core.EdgeType)
	assert.Equal(suite.T(), 2, len(keynames))
	stringset.Clear()
	stringset.AddArray(keynames)
	assert.Equal(suite.T(), 2, stringset.Count())
	assert.True(suite.T(), stringset.Contains("name"))
	assert.True(suite.T(), stringset.Contains("weight"))

	// this is what we will test
	// v1 points to v2 and v3 which both point to v4 which points back to v1

	// 		v2
	// 	/		\
	// v1			v4 - v1
	// 	\		/
	// 		v3
	
	//

	vid1 := []byte("vertex1")
	vid2 := []byte("vertex2")
	vid3 := []byte("vertex3")
	vid4 := []byte("vertex4")

	vertex1,_ := suite.TestGraph.AddVertex(vid1)
	vertex2,_ := suite.TestGraph.AddVertex(vid2)
	vertex3,_ := suite.TestGraph.AddVertex(vid3)
	vertex4,_ := suite.TestGraph.AddVertex(vid4)

	vertex1.SetProperty("name", []byte("vertex-1"))
	vertex2.SetProperty("name", []byte("vertex-2"))
	vertex3.SetProperty("name", []byte("vertex-3"))
	vertex4.SetProperty("name", []byte("vertex-4"))

	vertices := suite.TestGraph.VerticesWithProp("name", "vertex-1")
	assert.True(suite.T(), vertices != nil)
	assert.Equal(suite.T(), 1, len(vertices))
	atomset.Clear()
	atomset.AddVertexArray(vertices)
	assert.True(suite.T(), atomset.Contains(vertex1))


	eid1 := []byte("edge1")
	eid2 := []byte("edge2")
	eid3 := []byte("edge3")
	eid4 := []byte("edge4")
	eid5 := []byte("edge5")

	edge1, _ := suite.TestGraph.AddEdge(eid1, vertex1, vertex2, "1 to 2")
	edge2, _ := suite.TestGraph.AddEdge(eid2, vertex1, vertex3, "1 to 3")
	edge3, _ := suite.TestGraph.AddEdge(eid3, vertex2, vertex4, "2 to 4")
	edge4, _ := suite.TestGraph.AddEdge(eid4, vertex3, vertex4, "3 to 4")
	edge5, _ := suite.TestGraph.AddEdge(eid5, vertex4, vertex1, "4 to 1")

	edge1.SetProperty("name", []byte("edge-1"))
	edge2.SetProperty("name", []byte("edge-2"))
	edge3.SetProperty("name", []byte("edge-3"))
	edge4.SetProperty("name", []byte("edge-4"))
	edge5.SetProperty("name", []byte("edge-5"))

	edges := suite.TestGraph.EdgesWithProp("name", "edge-5")
	assert.True(suite.T(), edges != nil)
	assert.Equal(suite.T(), 1, len(edges))
	atomset.Clear()
	atomset.AddEdgeArray(edges)
	assert.True(suite.T(), atomset.Contains(edge5))

	edge1.SetProperty("weight", []byte("10"))
	edge2.SetProperty("weight", []byte("10"))
	edge3.SetProperty("weight", []byte("10"))
	edge4.SetProperty("weight", []byte("20"))
	edge5.SetProperty("weight", []byte("20"))

	atomset.Clear()
	atomset.AddEdgeArray(suite.TestGraph.EdgesWithProp("weight", "10"))
	assert.Equal(suite.T(), 3, atomset.Count())
	assert.True(suite.T(), atomset.Contains(edge1))
	assert.True(suite.T(), atomset.Contains(edge2))
	assert.True(suite.T(), atomset.Contains(edge3))

	atomset.Clear()
	atomset.AddEdgeArray(suite.TestGraph.EdgesWithProp("weight", "20"))
	assert.Equal(suite.T(), 2, atomset.Count())
	assert.True(suite.T(), atomset.Contains(edge4))
	assert.True(suite.T(), atomset.Contains(edge5))

	suite.TestGraph.DelEdge(edge5)

	atomset.Clear()
	atomset.AddEdgeArray(suite.TestGraph.EdgesWithProp("weight", "20"))
	assert.Equal(suite.T(), 1, atomset.Count())
	assert.True(suite.T(), atomset.Contains(edge4))

	suite.TestGraph.DelVertex(vertex1)
	atomset.Clear()
	atomset.AddVertexArray(suite.TestGraph.VerticesWithProp("name", "vertex-1"))
	assert.Equal(suite.T(), 0, atomset.Count())
	atomset.Clear()
	atomset.AddEdgeArray(suite.TestGraph.EdgesWithProp("weight", "10"))
	assert.Equal(suite.T(), 1, atomset.Count())
	assert.True(suite.T(), atomset.Contains(edge3))
}
