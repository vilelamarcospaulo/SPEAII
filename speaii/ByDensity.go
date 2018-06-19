package speaii

//ByDensity ::
type ByDensity []Individual

func (s ByDensity) Len() int {
	return len(s)
}

func (s ByDensity) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByDensity) Less(i, j int) bool {
	return s[i].Density > s[j].Density
}
