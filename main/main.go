package main

import (
	"fmt"
	"github.com/Douirat/lem-in/logic"
	"os"
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
	colony.DisplayColony()
	fmt.Println(colony.Start)
	fmt.Println(colony.End)
}
