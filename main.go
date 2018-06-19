package main

import (
	"SPEAII/speaii"
	"fmt"
)

func main() {
	ag := speaii.SPEAII{}
	ag.Run(500, 1000, 600, 0.2)

	for _, ind := range ag.CurrentPopulation {
		fmt.Println(ind.DNA, "||", ind.Goals)
	}
}
