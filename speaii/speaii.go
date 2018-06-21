package speaii

import (
	"fmt"
	"math"
	"math/rand"
	"sort"

	"github.com/Arafatk/glot"
)

//SPEAII :: Representacao da estrutura do AG do tipo SPEAII
type SPEAII struct {
	CurrentPopulation   []Individual
	ReferencePopulation []Individual

	PopulationSize          int
	ReferencePopulationSize int

	Generation          int
	MutationProbability float64

	plot *glot.Plot

	ParetoOptimal        []Individual
	ErrorRate            float64
	GenerationalDistance float64
	ParetoSubset         float64
	Spread               float64
	MaximumSpread        float64
}

//Run :: inicializa a configuração e processa o ag
func (speaii *SPEAII) Run(Generations int, PopulationSize int, ReferencePopulationSize int, MutationProbability float64) {
	speaii.Generation = 0
	speaii.PopulationSize = PopulationSize
	speaii.ReferencePopulationSize = ReferencePopulationSize

	speaii.MutationProbability = MutationProbability

	speaii.plot, _ = glot.NewPlot(2, true, true)

	speaii.newPopulation()
	for speaii.Generation = 1; speaii.Generation <= Generations; speaii.Generation++ {
		speaii.nextPopulation()
		speaii.DoPlot()
	}
	speaii.getNonDominated()
	speaii.DoPlot()
}

//GetNonDominated :: Copia os nao dominados do arquivo, para a populacao
func (speaii *SPEAII) getNonDominated() {
	speaii.CurrentPopulation = []Individual{}
	size := len(speaii.ReferencePopulation)
	for i := 0; i < size; i++ {
		speaii.ReferencePopulation[i].Rawfitness = 0
	}
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if speaii.ReferencePopulation[i].Dominate(&speaii.ReferencePopulation[j]) {
				speaii.ReferencePopulation[j].Rawfitness++
			}
		}
	}
	for i := 0; i < size; i++ {
		if speaii.ReferencePopulation[i].Rawfitness == 0 {
			speaii.CurrentPopulation = append(speaii.CurrentPopulation, speaii.ReferencePopulation[i])
		}
	}
	speaii.PopulationSize = len(speaii.CurrentPopulation)
}

//DoPlot :: Plota a populacao atual
func (speaii *SPEAII) DoPlot() {
	if speaii.plot == nil {
		speaii.plot, _ = glot.NewPlot(2, true, true)
	}

	xaxis := make([]float64, 1)
	yaxis := make([]float64, 1)
	for i := 0; i < len(speaii.CurrentPopulation); i++ {
		xaxis = append(xaxis, speaii.CurrentPopulation[i].Goals[0])
		yaxis = append(yaxis, speaii.CurrentPopulation[i].Goals[1])
	}

	points := [][]float64{xaxis, yaxis}
	speaii.plot.AddPointGroup(" ", "points", points)

	title := fmt.Sprintf("%s%d", "Generation: ", speaii.Generation)
	speaii.plot.SetTitle(title)

	speaii.plot.SetXLabel("SUM(sin(Pi * N))")
	speaii.plot.SetYLabel("SUM(sin(Pi * N))")

	speaii.plot.ResetPlot()
}

//NewPopulation :: Cria uma população inicial aleatoria
func (speaii *SPEAII) newPopulation() {
	speaii.CurrentPopulation = make([]Individual, speaii.PopulationSize)
	speaii.ReferencePopulation = []Individual{}

	for i := 0; i < speaii.PopulationSize; i++ {
		speaii.CurrentPopulation[i].NewRandom()
		speaii.CurrentPopulation[i].Eval()
	}

	speaii.fitness()
}

func (speaii SPEAII) selectParentByTour() (int, Individual) {
	index := rand.Intn(speaii.ReferencePopulationSize)
	individual := speaii.ReferencePopulation[index]
	for i := 1; i < 2; i++ {
		if newIndex := rand.Intn(speaii.ReferencePopulationSize); speaii.ReferencePopulation[newIndex].Better(individual) {
			index = newIndex
			individual = speaii.ReferencePopulation[index]
		}
	}
	return index, individual
}

//NextPopulation :: Gera a população t + 1, com base na atual (t)
func (speaii *SPEAII) nextPopulation() {
	newPopulation := make([]Individual, speaii.PopulationSize)

	for i := 0; i < speaii.PopulationSize; i += 2 {
		indexParent1, parent1 := speaii.selectParentByTour()
		indexParent2, parent2 := speaii.selectParentByTour()

		for indexParent1 == indexParent2 {
			indexParent2, parent2 = speaii.selectParentByTour()
		}

		var child1, child2 Individual
		child1.Initialize()
		child2.Initialize()
		Crossover(parent1, parent2, &child1, &child2)

		child1.Mutation(speaii.MutationProbability)
		child2.Mutation(speaii.MutationProbability)

		//Avalia os filhos gerados de acordo com o novo DNA
		child1.Eval()
		child2.Eval()

		newPopulation[i] = child1
		newPopulation[i+1] = child2
	}

	speaii.CurrentPopulation = newPopulation
	speaii.fitness()
}

//AppendPopulation :: Copia uma lista para uma nova lista de ponteiros
func appendPopulation(union *([]*Individual), population *[]Individual) {
	size := len(*population)
	for i := 0; i < size; i++ {
		(*population)[i].ResetValues()
		*union = append(*union, &(*population)[i])
	}
}

//Fitness :: calcula o fitness para cada individual nas populações
func (speaii *SPEAII) fitness() {
	var union []*Individual

	// Concatena a referencia de todos os individuos tando na populacao T quanto na populaco de referencia
	appendPopulation(&union, &speaii.CurrentPopulation)
	appendPopulation(&union, &speaii.ReferencePopulation)

	dominatedBy(union)
	density(union)
	speaii.mangeReferencePopulation(union)
}

//DominatedBy :: Calcula para cada individuo, os parametros de strength e rawfitness, em relacao a uniao da pop atual e a pop de referencia
func dominatedBy(union []*Individual) {
	//Conta quantos individuos B um individuo A domina,
	//e adiciona A na lista de dominantes de B
	size := len(union)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if union[i].Dominate(union[j]) {
				union[i].Strength++
				union[j].DominatedBy = append(union[j].DominatedBy, i)
			}
		}
	}

	//percorre a lista de dominantes de um individuo, e soma a "força" de seus dominantes
	for i := 0; i < size; i++ {
		for j := 0; j < len(union[i].DominatedBy); j++ {
			union[i].Rawfitness += union[union[i].DominatedBy[j]].Strength
		}
	}
}

//Density :: Calcula o parametro density de cada individuo na população,
//sendo esse, a somatoria da distancia de acordo a um parametro
func density(union []*Individual) {
	size := len(union)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			union[j].GoalsDistance(*union[i])
		}

		ordered := make([]*Individual, size)
		copy(ordered, union)

		sort.Sort(ByDistance(ordered))
		k := int(math.Sqrt(float64(size)))

		union[i].Density = 1 / (ordered[k].Distance + 2)
		union[i].Fitness = union[i].Rawfitness + union[i].Density
	}
}

//MangeReferencePopulation :: Faz a copia dos individuos nao dominados para populacao de referencia, e controla o tamanho dela
func (speaii *SPEAII) mangeReferencePopulation(union []*Individual) {
	speaii.ReferencePopulation = []Individual{}
	size := len(union)
	sizeReference := 0
	for i := 0; i < size; i++ {
		if union[i].Rawfitness == 0 { // Rawfitness == 0, individuo não dominado, copiar para o novo vetor
			sizeReference++
			speaii.ReferencePopulation = append(speaii.ReferencePopulation, *union[i])
		}
	}

	if sizeReference > speaii.ReferencePopulationSize { //Extrapolou o limite.
		for i := 0; i < sizeReference; i++ {
			ordered := make([]*Individual, sizeReference)

			speaii.ReferencePopulation[i].Density = 0

			for j := 0; j < sizeReference; j++ {
				speaii.ReferencePopulation[j].GoalsDistance(speaii.ReferencePopulation[i])
				ordered[j] = &speaii.ReferencePopulation[j]
			}

			sort.Sort(ByDistance(ordered))
			k := int(math.Sqrt(float64(sizeReference)))

			speaii.ReferencePopulation[i].Density = ordered[k].Distance
		}

		sort.Sort(ByDensity(speaii.ReferencePopulation))
		speaii.ReferencePopulation = speaii.ReferencePopulation[:speaii.ReferencePopulationSize]

	} else if sizeReference < speaii.ReferencePopulationSize { //Completar com os melhores dominados.
		sort.Sort(ByFitness(union))

		for i := 0; i < size && sizeReference < speaii.ReferencePopulationSize; i++ {
			if union[i].Rawfitness != 0 { // Rawfitness != 0, individuo deve ser dominado, pois os não dominados já foram copiados
				sizeReference++
				speaii.ReferencePopulation = append(speaii.ReferencePopulation, *union[i])
			}
		}
	}
}

//CalcErrorRate :: Calcula o error rate da populacao atual
func (speaii *SPEAII) CalcErrorRate() float64 {
	speaii.ErrorRate = 0.0

	for _, ind := range speaii.CurrentPopulation {
		for _, indRef := range speaii.ParetoOptimal {
			if indRef.Dominate(&ind) {
				speaii.ErrorRate++
				break
			}
		}
	}

	speaii.ErrorRate /= float64(speaii.PopulationSize)
	return speaii.ErrorRate
}

//CalcGenerationalDistance :: Calculaa generational distance
func (speaii *SPEAII) CalcGenerationalDistance() float64 {
	speaii.GenerationalDistance = 0.0

	nonDominated := make([]*Individual, 0)
	for i := 0; i < speaii.PopulationSize; i++ {
		nonDominated = append(nonDominated, &speaii.CurrentPopulation[i])
	}

	sizeRef := len(speaii.ParetoOptimal)
	for i := 0; i < speaii.PopulationSize; i++ {
		for j := 0; j < sizeRef; j++ {
			speaii.ParetoOptimal[j].GoalsDistance(speaii.CurrentPopulation[i])
		}
		sort.Sort(ByDistance(nonDominated))

		speaii.GenerationalDistance += nonDominated[0].Distance
	}

	speaii.GenerationalDistance = math.Sqrt(speaii.GenerationalDistance)
	speaii.GenerationalDistance /= float64(speaii.PopulationSize)

	return speaii.GenerationalDistance
}

//CalcParetoSubset :: Calcula o paretosubset
func (speaii *SPEAII) CalcParetoSubset() float64 {
	speaii.ParetoSubset = (1 - speaii.ErrorRate) * float64(speaii.PopulationSize)
	return speaii.ParetoSubset
}

//CalcSpread :: Calcula o spread
func (speaii *SPEAII) CalcSpread() float64 {
	nonDominated := make([]*Individual, 0)

	for i := 0; i < speaii.PopulationSize; i++ {
		speaii.CurrentPopulation[i].CurrentGoal = 0
		nonDominated = append(nonDominated, &speaii.CurrentPopulation[i])
	}
	sort.Sort(ByGoal(nonDominated)) //ORDERNA PELO VETOR X FIXO

	dIE := 0.0

	byGoalSorted := make([]*Individual, speaii.PopulationSize)
	copy(byGoalSorted, nonDominated)
	for goal := 0; goal < speaii.CurrentPopulation[0].GoalsSize; goal++ {
		for i := 0; i < speaii.PopulationSize; i++ {
			byGoalSorted[i].CurrentGoal = goal
		}
		sort.Sort(ByGoal(byGoalSorted))

		dIE += byGoalSorted[0].GoalsDistance(*nonDominated[0])                                             // soma os extremos superiores em realacao ao objetivo
		dIE += byGoalSorted[speaii.PopulationSize-1].GoalsDistance(*nonDominated[speaii.PopulationSize-1]) // soma os extremos inferirores
	}

	dAverage := 0.0
	for i := 0; i < speaii.PopulationSize-1; i++ {
		dAverage += nonDominated[i].GoalsDistance(*nonDominated[i+1])
	}
	dAverage /= float64(speaii.PopulationSize - 1)

	sum := 0.0
	for i := 0; i < speaii.PopulationSize-1; i++ {
		di := nonDominated[i].GoalsDistance(*nonDominated[i+1])
		sum += math.Abs(di - dAverage)
	}

	a := dIE + sum
	b := dIE + ((dAverage) * float64(speaii.PopulationSize-1))
	speaii.Spread = a / b

	return speaii.Spread
}

//CalcMaximumSpread :: Calcula o MaximumSpread
func (speaii *SPEAII) CalcMaximumSpread() float64 {
	nonDominated := make([]*Individual, 0)

	for i := 0; i < speaii.PopulationSize; i++ {
		nonDominated = append(nonDominated, &speaii.CurrentPopulation[i])
	}
	speaii.MaximumSpread = 0.0

	for goal := 0; goal < speaii.CurrentPopulation[0].GoalsSize; goal++ {
		for i := 0; i < speaii.PopulationSize; i++ {
			nonDominated[i].CurrentGoal = goal
		}
		sort.Sort(ByGoal(nonDominated))

		speaii.MaximumSpread += nonDominated[0].GoalsDistance(*nonDominated[speaii.PopulationSize-1])
	}

	speaii.MaximumSpread = math.Sqrt(speaii.MaximumSpread)
	return speaii.MaximumSpread
}
