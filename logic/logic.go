package logic

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Douirat/lem-in/data"
)

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
					fmt.Println("Start found:", str)
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
					fmt.Println("End found:", str)
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
