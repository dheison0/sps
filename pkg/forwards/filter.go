package forwards

var Filter = map[string]bool{}

func AddFilter(t string) {
	Filter[t] = true
}
