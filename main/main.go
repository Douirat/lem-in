package main

import (
	"fmt"
	"os"

	"github.com/Douirat/lem-in/logic"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Not enough arguments!")
		return
	}
	colony := logic.NewColony()
	err := colony.RockAndRoll(args[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("The number of ants is: %d\n", colony.Ants)
	// colony.DisplayColony()
	graph := colony.CreateFarm()
	graph.Display()
	pathsToRoom := graph.FindShortestPath(colony.Start.Name, colony.End.Name)
	allPaths := graph.FindAllPathsToDestination(colony.Start.Name, colony.End.Name)
	fmt.Println(allPaths)
	fmt.Println(pathsToRoom)
}
