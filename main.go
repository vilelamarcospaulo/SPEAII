package main

import (
	"SPEAII/speaii"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func readOptimal() []speaii.Individual {
	optimal := []speaii.Individual{}
	file, _ := os.Open("pareto.dat")

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		buffer := scanner.Text()
		values := strings.Split(buffer, " ")

		ind := speaii.Individual{}
		ind.NewRandom()

		for i := 0; i < ind.GoalsSize; i++ {
			ind.Goals[i], _ = strconv.ParseFloat(values[i], 64)
		}

		optimal = append(optimal, ind)
	}

	return optimal
}

func main() {
	// agAux := speaii.SPEAII{}
	// agAux.CurrentPopulation = readOptimal()
	// agAux.DoPlot()

	ag := speaii.SPEAII{}
	start := time.Now()
	ag.Run(500, 1000, 600, 0.02)
	elapsed := time.Since(start)
	optimal := readOptimal()

	ag.ParetoOptimal = optimal
	fmt.Println("Time: ", elapsed)
	fmt.Println("Pareto size: ", ag.PopulationSize)
	fmt.Println("Error rate: ", ag.CalcErrorRate())
	fmt.Println("Pareto subset: ", ag.CalcParetoSubset())
	fmt.Println("Generational distance: ", ag.CalcGenerationalDistance())
	fmt.Println("Spread : ", ag.CalcSpread())
	fmt.Println("Maximum Spread (m3): ", ag.CalcMaximumSpread())
}
