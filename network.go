package contagiongo

import (
	"bytes"
	"contagion/utils"
	"fmt"
)

// HostNetwork interface describes a host population connected together as
// a network.
type HostNetwork interface {
	// ConnectedPopSize returns the total number of hosts in the network.
	ConnectedPopSize() int
	// GetNeighbors retrieves the unordered list of neighbors from
	// the adjacency matrix.
	GetNeighbors(ID int) (neighbors []int)
	// ConnectionExists checks if a connection a-b exists in the adjacency matrix.
	ConnectionExists(a, b int) bool
	// AddConnection adds a one way connection a-b to the adjacency matrix if
	// the connection a-b does not exists. Returns an error if given connection
	// already exists.
	AddConnection(a, b int) error
	// AddWeightedConnection adds a one way connection a-b to the adjacency
	// matrix with a given weight w.
	AddWeightedConnection(a, b int, w float64) error
	// UpdateConnectionWeight changes the weight value of an existing
	// connection. If the given connection does not exist, nothing is
	// updated.
	UpdateConnectionWeight(a, b int, w float64) error
	// UpsertConnectionWeight changes the weight value of an existing
	// connection or creates a new connection with the given weight if the
	// connection does not exist.
	UpsertConnectionWeight(a, b int, w float64)

	// AddWeightedBiConnection adds a two way reciprocal connection to the
	// adjacency matrix with a given weight.
	AddWeightedBiConnection(a, b int, w float64) error

	// DeleteConnection removes a one way connection a-b.
	DeleteConnection(a, b int) error

	// Copy returns a new copy of the adjacency matrix.
	// Changes made to the original copy will not affect the new copy
	// and changes made to the copy will likewise not affect the original.
	Copy() adjacencyMatrix

	// Dump serialized the adjacency matrix into a string stored as a byteslice.
	Dump() []byte
}

// AdjacencyMatrix is a 2D map that represents
// connections between hosts using their UIDs as index.
type adjacencyMatrix map[int]map[int]float64

func (m adjacencyMatrix) ConnectedPopSize() int {
	hostIdsSet := make(map[int]bool)
	for i, hosts := range m {
		hostIdsSet[i] = true
		for j := range hosts {
			hostIdsSet[j] = true
		}
	}
	return len(hostIdsSet)
}

func (m adjacencyMatrix) GetNeighbors(ID int) (neighbors []int) {
	for j := range m[ID] {
		neighbors = append(neighbors, j)
	}
	return
}

func (m adjacencyMatrix) ConnectionExists(a, b int) bool {
	if _, exists := m[a]; !exists {
		return false
	}
	if _, exists := m[a][b]; !exists {
		return false
	}
	return true
}

func (m adjacencyMatrix) AddConnection(a, b int) error {
	w := float64(1) // default unweighted
	return m.AddWeightedConnection(a, b, w)
}

func (m adjacencyMatrix) AddWeightedConnection(a, b int, w float64) error {
	if m.ConnectionExists(a, b) {
		return fmt.Errorf("Connection (%d,%d): %f already exists", a, b, m[a][b])
	}
	// Check if the inner map has been initialized.
	// If not, initialize before assigning a value
	if _, exists := m[a]; !exists {
		m[a] = map[int]float64{}
	}
	m[a][b] = w
	return nil
}

func (m adjacencyMatrix) UpdateConnectionWeight(a, b int, w float64) error {
	if !m.ConnectionExists(a, b) {
		return fmt.Errorf("Connection (%d,%d) does not exist", a, b)
	}
	m[a][b] = w
	return nil
}

func (m adjacencyMatrix) UpsertConnectionWeight(a, b int, w float64) {
	err := m.AddWeightedConnection(a, b, w)
	// Error means existing connection present
	// Update instead
	if err != nil {
		m[a][b] = w
	}
}

func (m adjacencyMatrix) DeleteConnection(a, b int) error {
	if !m.ConnectionExists(a, b) {
		return fmt.Errorf("connection (%d,%d) does not exist", a, b)
	}
	delete(m[a], b)
	return nil
}

func (m adjacencyMatrix) Copy() adjacencyMatrix {
	n := make(adjacencyMatrix)
	for i, nbrs := range m {
		n[i] = make(map[int]float64)
		for j, v := range nbrs {
			n[i][j] = v
		}
	}
	return n
}

func (m adjacencyMatrix) Dump() []byte {
	b := new(bytes.Buffer)
	for idf, nbrs := range m {
		for idt, weight := range nbrs {
			_, err := b.WriteString(fmt.Sprintf(`%d,%d: %f\n`, idf, idt, weight))
			utils.Check(err)
		}
	}
	return b.Bytes()
}

// Following are convenience methods for bidirectional connections

// BiConnectionExists checks if both connections a-b, and b-a exists in the
// adjacency matrix.
func (m *adjacencyMatrix) BiConnectionExists(a, b int) bool {
	if m.ConnectionExists(a, b) && m.ConnectionExists(b, a) {
		return true
	}
	return false
}

// AddBiConnection adds a two way reciprocal connection a-b and b-a to the
// adjacency matrix.
func (m adjacencyMatrix) AddBiConnection(a, b int) error {
	w := float64(1) // default unweighted
	return m.AddWeightedBiConnection(a, b, w)
}

// AddWeightedBiConnection adds a two way reciprocal connection to the
// adjacency matrix with a given weight.
func (m adjacencyMatrix) AddWeightedBiConnection(a, b int, w float64) error {
	// a-b must not create a self loop
	if a == b {
		return fmt.Errorf("start and end nodes are the same")
	}
	// Both a-b and b-a connections must not exist
	// If either forward or reverse orientation exists, returns an error
	if abExists := m.ConnectionExists(a, b); abExists {
		return fmt.Errorf("connection (%d,%d): %f already exists", a, b, m[a][b])
	}
	if baExists := m.ConnectionExists(b, a); baExists {
		return fmt.Errorf("connection (%d,%d): %f already exists", b, a, m[b][a])
	}
	// Guaranteed that both a-b and b-a do not exist
	m.AddWeightedConnection(a, b, w)
	m.AddWeightedConnection(b, a, w)
	return nil
}

// UpdateBiConnectionWeight changes the weight value of an existing
// two-way connection. Both connections must exist for any change to occur.
func (m adjacencyMatrix) UpdateBiConnectionWeight(a, b int, w float64) error {
	// a-b must not create a self loop
	if a == b {
		return fmt.Errorf("start and end nodes are the same")
	}
	// Both a-b and b-a connections must exist in order to update
	// If either forward or reverse orientation does not exist,
	// returns an error
	if !m.ConnectionExists(a, b) {
		return fmt.Errorf("Connection (%d,%d) does not exist", a, b)
	}
	if !m.ConnectionExists(b, a) {
		return fmt.Errorf("connection (%d,%d) does not exist", b, a)
	}
	// Guaranteed that both a-b and b-a exist
	m[a][b] = w
	m[b][a] = w
	return nil
}

// UpsertBiConnectionWeight changes the weight value of an existing
// connection or creates a new two-way connection with the given weight if
// both connection do not exist. Panics if only one connection exists.
func (m adjacencyMatrix) UpsertBiConnectionWeight(a, b int, w float64) {
	m.UpsertConnectionWeight(a, b, w)
	m.UpsertConnectionWeight(b, a, w)
}

// DeleteBiConnection removes the connection a-b and b-a.
func (m adjacencyMatrix) DeleteBiConnection(a, b int) error {
	// a-b must not create a self loop
	if a == b {
		return fmt.Errorf("start and end nodes are the same")
	}
	if !m.ConnectionExists(a, b) {
		return fmt.Errorf("connection (%d,%d) does not exist", a, b)
	}
	if !m.ConnectionExists(b, a) {
		return fmt.Errorf("connection (%d,%d) does not exist", b, a)
	}
	delete(m[a], b)
	delete(m[b], a)
	return nil
}

// EmptyAdjacencyMatrix creates a new 2D mapping with no contents.
func EmptyAdjacencyMatrix() adjacencyMatrix {
	m := make(adjacencyMatrix)
	return m
}

// // GNPAdjacencyMatrix generates an Erdos-Renyi graph using the given number
// // of nodes n, probability p and a seed integer.
// func GNPAdjacencyMatrix(n int, p float64, seed int64) adjacencyMatrix {
// 	// http://homepage.divms.uiowa.edu/~sriram/196/spring12/lectureNotes/Lecture4.pdf
// 	rand.Seed(int64(seed))

// 	g := make(adjacencyMatrix)
// 	for i := 0; i < n; i++ {
// 		for j := i; j < n; j++ {
// 			if rv.Binomial(1, p) == 1 {
// 				g.AddConnection(i, j)
// 				g.AddConnection(j, i)
// 			}
// 		}
// 	}
// 	return g
// }

// // BAAdjacencyMatrix generates an Barabasi-Alberts preferential connection
// // graph using the given number of nodes n, parameter m and a seed integer.
// func BAAdjacencyMatrix(n, m int, seed int64) adjacencyMatrix {
// 	// TODO: TEST!!!
// 	rand.Seed(int64(seed))
// 	// Create initial fully connected graph
// 	g := make(adjacencyMatrix)

// 	nodeIDs := utils.MakeRange(0, n)
// 	utils.PermuteIntsInplace(&nodeIDs)

// 	// Start with m nodes in a line
// 	prevI := 0
// 	nodeCount := 0
// 	for i := 1; i < m; i++ {
// 		g.AddConnection(nodeIDs[prevI], nodeIDs[i])
// 		g.AddConnection(nodeIDs[i], nodeIDs[prevI])
// 		prevI = i
// 		nodeCount++
// 	}

// 	getProb := func(ids ...int) (prob []float64) {
// 		for _, i := range ids {
// 			prob = append(prob, float64(len(g.GetNeighbors(i))))
// 		}
// 		sum := float64(utils.SumFloat64s(prob...))

// 		for i, p := range prob {
// 			prob[i] = math.Log1p(p / sum)
// 		}
// 		return
// 	}

// 	// Add new nodes and connect to node with most neighbors
// 	for i := m; i < n; i++ {
// 		connectedTo := []int{}
// 		for c := 0; c < m; c++ {
// 			currentProb := getProb(utils.MakeRangeWithout(0, n, connectedTo...)...)
// 			j := MultinomialWhereInt(1, currentProb, 1)[0]

// 			g.AddConnection(nodeIDs[i], j)
// 			g.AddConnection(j, nodeIDs[i])
// 			connectedTo = append(connectedTo, j)
// 		}
// 	}

// 	return g
// }
