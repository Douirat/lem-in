package logic

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Douirat/lem-in/data"
)

// Declare a global variable to hold the size of the graph:
var SIZE int

// Declare a structure to represent the room:
type Room struct {
	Name       string
	CorX, CorY int
	Next       *Room
}

// declare a structre to represent the graph using the adjacency list:
type Colony struct {
	Ants       int
	Start, End *Room
	Farm       map[string]*Room
}

// Declare a graph to ease the extraction of paths:
type Graph struct {
	Farm map[string][]string
}

// Instantiate a colony:
func NewColony() *Colony {
	colony := new(Colony)
	colony.Ants = 0
	colony.Start = nil
	colony.End = nil
	colony.Farm = make(map[string]*Room)
	return colony
}

// Instantiate a new room:
func NewRoom(str string) (*Room, error) {
	var err error
	room := new(Room)
	data := strings.Split(str, " ")
	if len(data) != 3 {
		return nil, errors.New("invalid data format")
	}
	room.Name = data[0]
	room.CorX, err = strconv.Atoi(data[1])
	if err != nil {
		return nil, err
	}
	room.CorY, err = strconv.Atoi(data[2])
	if err != nil {
		return nil, err
	}
	room.Next = nil
	return room, nil
}

// Add a new room to the colony:
func (colony *Colony) AddRoom(str string) (*Room, error) {
	room, err := NewRoom(str)
	if err != nil {
		return nil, err
	}
	colony.Farm[room.Name] = room
	SIZE++
	return room, nil
}

// Add a tunnel between tow rooms:
func (colony *Colony) AddTunnel(str string) error {
	data := strings.Split(str, "-")
	if len(data) != 2 {
		return errors.New("error data format, tunnels not valid")
	}
	roomSrc := &Room{}
	roomSrc.Name = colony.Farm[data[0]].Name
	roomSrc.CorX = colony.Farm[data[0]].CorX
	roomSrc.CorY = colony.Farm[data[0]].CorY
	roomDst := &Room{}
	roomDst.Name = colony.Farm[data[1]].Name
	roomDst.CorX = colony.Farm[data[1]].CorX
	roomDst.CorY = colony.Farm[data[1]].CorY
	roomDst.Next = colony.Farm[data[0]].Next
	if roomDst.Name == colony.Start.Name {
		colony.Start = roomDst
	} else if roomDst.Name == colony.End.Name {
		colony.End = roomDst
	}
	if roomSrc.Name == colony.Start.Name {
		colony.Start = roomSrc
	} else if roomSrc.Name == colony.End.Name {
		colony.End = roomSrc
	}
	colony.Farm[data[0]].Next = roomDst
	roomSrc.Next = colony.Farm[data[1]].Next
	colony.Farm[data[1]].Next = roomSrc
	return nil
}

// Create a graph with only room names to ease the shortest paths extraction:
func (colony *Colony) CreateFarm() *Graph {
	graph := new(Graph)
	graph.Farm = make(map[string][]string)
	for name, room := range colony.Farm {
		if _, exist := graph.Farm[name]; !exist {
			graph.Farm[name] = []string{}
		}
		temp := room.Next
		for temp != nil {
			graph.Farm[name] = append(graph.Farm[name], temp.Name)
			temp = temp.Next
		}
	}
	return graph
}

// Display the graph:
func (graph *Graph) Display() {
	for key, values := range graph.Farm {
		fmt.Printf("We are visiting the room: %s\n", key)
		for _, value := range values {
			fmt.Printf("%s <----> %s ", key, value)
		}
		fmt.Printf("\n")
	}
}

// Find the shortest path;
// BFS function to find all shortest paths between two cities
func (graph *Graph) FindShortestPath(start, end string) [][]string {

	//  A queue in the BFS, each element is a room and a paths to it;
	// Why a Queue for BFS?
	// The queue in BFS is essential for the breadth-first-traversal of the graph.
	// The idea behind BFS is to explore all nodes at the current "level" (or "distance")
	// before moving on to the next level. This ensures that we explore rooms
	// in increasing order of their distance from the start room. The queue
	// helps by storing rooms in the order they are visited, while also allowing us to expand
	// paths correctly.
	// queue := NewQueue()

	// Start with the path containing just the start room
	queue := [][]string{{start}}

	// Map to track the shortest distnce to each room
	// In the algorithm, the distances map tracks the shortest number
	// of steps (or edges) from the start room to every other room (or node).
	// The goal is to find the shortest path from the starting room to the
	// destination room, which means finding the minimum number of edges
	// you need to traverse between two rooms.
	distances := make(map[string]int)
	distances[start] = 0 // When the algorithm starts, we set the distance to the start room as 0 because we're already there.

	// The pathsToRoom map is designed to keep track of all possible shortest paths from
	// the starting room to each other room. Each entry in the map stores a list of paths
	// leading to a particular room. Since we are using Breadth-First Search (BFS),
	// the idea is that BFS explores the shortest paths first, and we can collect
	// all the shortest paths as we go along.
	// Map to store all paths leading to each room
	pathsToRoom := make(map[string][][]string)
	pathsToRoom[start] = [][]string{{start}} // Start room has one path: itself

	// BFS: Exploring cities one by one:
	for len(queue) > 0 {
		// Dequeue the first element in the queue (FIFO):
		currentPath := queue[0]
		queue = queue[1:]

		// Current room is the last room in the current path:
		currentRoom := currentPath[len(currentPath)-1]

		// Explore all neighbors of the current room:
		for _, neighbor := range graph.Farm[currentRoom] {
			// Calculate the distance to the neighbor via the current city:
			newDistance := distances[currentRoom] + 1

			// If visiting the neighbor for the first time record its distance:
			if _, visited := distances[neighbor]; !visited {
				distances[neighbor] = newDistance

				// Add this new path to the queue for further exploration
				queue = append(queue, append([]string{}, append(currentPath, neighbor)...))
				// ====> Create a copy of currentPath
				// newPath := append([]string{}, currentPath...)
				// // Add the neighbor to the new path
				// newPath = append(newPath, neighbor)
				// // Append the new path to the queue
				// queue = append(queue, newPath)

				// Initialize paths to this neighbor with the current path
				pathsToRoom[neighbor] = [][]string{append([]string{}, append(currentPath, neighbor)...)}
			} else if distances[neighbor] == newDistance {
				// If we reach the neighbor with the same shortest distance, add the path
				pathsToRoom[neighbor] = append(pathsToRoom[neighbor], append([]string{}, append(currentPath, neighbor)...))
			}
		}
	}
	fmt.Println(pathsToRoom)
	fmt.Println(distances)
	// Return all shortest paths to the destination city
	return pathsToRoom[end]
}

// Formulating the colny graph based on the input extracted from the file:
func (colony *Colony) RockAndRoll(fileName string) error {
	started := false
	ended := false
	// end := false
	data, err := data.ReadFile(fileName)
	if err != nil {
		return err
	}
	fmt.Println(len(data))
	colony.Ants, err = strconv.Atoi(data[0])
	if err != nil {
		return err
	}
	data = data[1:]

	// Define regex patterns for each type of line
	patterns := map[string]*regexp.Regexp{
		"start":   regexp.MustCompile(`^##start$`),
		"room":    regexp.MustCompile(`^(.*)\s+(\d+)\s+(\d+)$`),
		"end":     regexp.MustCompile(`^##end$`),
		"tunnel":  regexp.MustCompile(`^([a-zA-Z0-9]+)-([a-zA-Z0-9]+)$`),
		"comment": regexp.MustCompile(`^#.*`),
	}

	// Iterate over the input data
	for _, str := range data {
		// Check each pattern
		for key, rg := range patterns {
			if rg.MatchString(str) {
				// Process based on the matched pattern
				switch key {
				case "start":
					// Handle the start pattern
					// fmt.Println("Start found:", str)
					started = true
					continue
				case "room":
					// Handle the room pattern
					if started && colony.Start == nil {
						colony.Start, err = colony.AddRoom(str)
						if err != nil {
							return err
						}
						started = false
						continue
					}
					if ended && colony.End == nil {
						colony.End, err = colony.AddRoom(str)
						if err != nil {
							return err
						}
						ended = false
						// end = true
						continue
					}
					_, err = colony.AddRoom(str)
					if err != nil {
						return err
					}
				case "end":
					// Handle the end pattern
					ended = true
					// fmt.Println("End found:", str)
					continue
				case "tunnel":
					// if !end {
					// 	return errors.New("wrond data format, tunnel before end flag")
					// }
					if colony.End == nil {
						return errors.New("wrond data format, tunnel before end flag")
					}
					// Handle the tunnel pattern
					colony.AddTunnel(str)
				case "comment":
					// Handle the comment pattern
					continue
				}
			}
		}
	}
	return nil
}

// Display colony:
func (colony *Colony) DisplayColony() {
	for name, room := range colony.Farm {
		fmt.Printf("We are visiting the room: %s\n", name)
		fmt.Printf("the room coordinates are [%d, %d]\n", room.CorX, room.CorY)
		temp := room.Next
		for temp != nil {
			fmt.Printf("There is a tunel between %s --- and ---> %s with the coordinates[%d, %d]\n", name, temp.Name, temp.CorX, temp.CorY)
			temp = temp.Next
		}
	}
}
