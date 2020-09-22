package models

import (
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
)

var (
	distanceGraph = simple.NewUndirectedGraph()
	distancePaths path.AllShortest
)

// init builds the graph that represents highway connections between regions in Dominican Republic
// see graph_do.png for visual representation
func init() {
	distanceGraph.AddNode(cibaoNordesteRegionNode)
	distanceGraph.AddNode(cibaoNoroesteRegionNode)
	distanceGraph.AddNode(cibaoNorteRegionNode)
	distanceGraph.AddNode(cibaoSurRegionNode)
	distanceGraph.AddNode(elValleRegionNode)
	distanceGraph.AddNode(enriquilloRegionNode)
	distanceGraph.AddNode(higuamoRegionNode)
	distanceGraph.AddNode(ozamaRegionNode)
	distanceGraph.AddNode(valdesiaRegionNode)
	distanceGraph.AddNode(yumaRegionNode)

	distanceGraph.SetEdge(distanceGraph.NewEdge(cibaoNordesteRegionNode, cibaoNorteRegionNode))
	distanceGraph.SetEdge(distanceGraph.NewEdge(cibaoNordesteRegionNode, cibaoSurRegionNode))
	distanceGraph.SetEdge(distanceGraph.NewEdge(cibaoNordesteRegionNode, ozamaRegionNode))
	distanceGraph.SetEdge(distanceGraph.NewEdge(cibaoNoroesteRegionNode, cibaoNorteRegionNode))
	distanceGraph.SetEdge(distanceGraph.NewEdge(cibaoNorteRegionNode, cibaoSurRegionNode))
	distanceGraph.SetEdge(distanceGraph.NewEdge(cibaoSurRegionNode, ozamaRegionNode))
	distanceGraph.SetEdge(distanceGraph.NewEdge(elValleRegionNode, enriquilloRegionNode))
	distanceGraph.SetEdge(distanceGraph.NewEdge(elValleRegionNode, valdesiaRegionNode))
	distanceGraph.SetEdge(distanceGraph.NewEdge(higuamoRegionNode, ozamaRegionNode))
	distanceGraph.SetEdge(distanceGraph.NewEdge(higuamoRegionNode, yumaRegionNode))
	distanceGraph.SetEdge(distanceGraph.NewEdge(ozamaRegionNode, valdesiaRegionNode))

	distancePaths = path.DijkstraAllPaths(distanceGraph)
}

// GetDistance returns the distance between two regions for billing
func GetDistance(origin, destination Region) float64 {
	originNode := RegionToRegionNode[origin]
	destinationNode := RegionToRegionNode[destination]

	// need to add 1 for the purpose of billing inside the same region
	return distancePaths.Weight(originNode.ID(), destinationNode.ID()) + 1
}
