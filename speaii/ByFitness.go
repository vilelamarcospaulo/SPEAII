package speaii

//ByFitness ::
type ByFitness []*Individual

func (s ByFitness) Len() int {
	return len(s)
}

func (s ByFitness) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByFitness) Less(i, j int) bool {
	return s[i].Fitness < s[j].Fitness
}
