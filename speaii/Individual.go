package speaii

import (
	"math"
	"math/rand"
)

//Individual :: Representacao do individuo do AG
type Individual struct {
	DNA     *([]float64)
	DNASize int

	Goals     []float64
	GoalsSize int

	Fitness    float64
	Strength   float64
	Rawfitness float64
	Density    float64

	//Utilizados como auxiliares no calculo de fitness
	//nao necessariamente corresponde a uma caracteristica do inviduo
	Distance    float64
	DominatedBy []int
}

func randomValidValue() float64 {
	return random(0, 6)
}

//Initialize :: Inicializa as variaveis  do individuo
func (individual *Individual) Initialize() {
	individual.DNASize, individual.GoalsSize = 5, 2

	dna := make([]float64, individual.DNASize)

	individual.DNA = &dna
	individual.Goals = make([]float64, individual.GoalsSize)

	individual.Fitness = 0
	individual.Strength = 0
	individual.Rawfitness = 0
	individual.Density = 0
}

//NewRandom :: Gera um individuo
func (individual *Individual) NewRandom() {
	individual.Initialize()

	for i := 0; i < individual.DNASize; i++ {
		(*individual.DNA)[i] = randomValidValue()
	}
}

//Eval :: Calcula os valores para cada objetivo com base no DNA
//deve mudar para cada problema.
func (individual *Individual) Eval() {
	for i := 0; i < individual.GoalsSize; i++ {
		individual.Goals[i] = 0.0
	}

	for i := 0; i < individual.DNASize; i++ {
		individual.Goals[0] += math.Sin(math.Pi * (*individual.DNA)[i])
		individual.Goals[1] += math.Cos(math.Pi * (*individual.DNA)[i])
	}

	individual.Fitness = 0
	individual.Strength = 0
	individual.Rawfitness = 0
	individual.Density = 0
}

func random(min float64, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

//Mutation :: Realiza uma mutação no individuo
//se o individuo sofrer mutação, sorteia uma posição e altera o valor dela
//por um outro valor no dominio valido
func (individual *Individual) Mutation(probability float64) {
	for i := 0; i < individual.DNASize; i++ {
		if rand.Float64() > probability {
			return
		}

		(*individual.DNA)[i] += random(-2, 2)
		if (*individual.DNA)[i] > 6 {
			(*individual.DNA)[i] = 6
		} else if (*individual.DNA)[i] < 0 {
			(*individual.DNA)[i] = 0
		}
	}
}

//Dominate :: Checa se um individuo é dominado por outro
// se for menor em ao menos um objeto é considerado falso
// se for maior em ao menos um e indiferente aos outros é true
func (individual *Individual) Dominate(other *Individual) bool {
	hasOneLessThan := false
	for i := 0; i < individual.GoalsSize; i++ {
		if individual.Goals[i] > other.Goals[i] {
			return false
		}
		if individual.Goals[i] < other.Goals[i] {
			hasOneLessThan = true
		}
	}
	return hasOneLessThan
}

//Better :: Checa se um individuo é melhor que o outro baseado no rank e crowding distance
func (individual *Individual) Better(other Individual) bool {
	return individual.Fitness > other.Fitness
}

//GoalsDistance :: Calcula a soma da distancia euclidiana de cada objetivo, entre os dois individuos
func (individual *Individual) GoalsDistance(other Individual) float64 {
	individual.Distance = 0

	for i := 0; i < individual.GoalsSize; i++ {
		individual.Distance += math.Pow(individual.Goals[i]-other.Goals[i], 2.0)
	}

	return individual.Distance
}

//ResetValues :: Reinicia os valores de parametros, como fitness, distance, etc..
func (individual *Individual) ResetValues() {
	individual.Fitness = 0
	individual.Strength = 0
	individual.Rawfitness = 0
	individual.Density = 0

	individual.Distance = 0
	individual.DominatedBy = []int{}
}

//Crossover :: Realiza o crosoover entre dois pais "parent1" e "parent2"
//e coloca o resultado em "child1" e "child2"
func Crossover(parent1 Individual, parent2 Individual, child1 *Individual, child2 *Individual) {
	for i := 0; i < parent1.DNASize; i++ {
		if rand.Intn(2) == 0 {
			(*child1.DNA)[i], (*child2.DNA)[i] = (*parent1.DNA)[i], (*parent2.DNA)[i]
		} else {
			(*child1.DNA)[i], (*child2.DNA)[i] = (*parent2.DNA)[i], (*parent1.DNA)[i]
		}
	}

	for i := 0; i < parent1.GoalsSize; i++ {
		child1.Goals[i], child2.Goals[i] = 0, 0
	}

	child1.ResetValues()
	child2.ResetValues()
}
