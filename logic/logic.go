package logic

import (
	"errors"
	"fmt"
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

// Formulating the colny graph based on the input extracted from the file:
func (colony *Colony) RockAndRoll(fileName string) error {
	data, err := data.ReadFile(fileName)
	if err != nil {
		return err
	}
	for _, str := range data {
		fmt.Println(str)
	}
	return nil
}
