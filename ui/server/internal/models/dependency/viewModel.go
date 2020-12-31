package dependency

type NodeViewModel struct {
	name  string
	value int64
}

type EdgeViewModel struct {
	source string
	target string
	value  int64
}

type ViewModel struct {
	nodes []NodeViewModel
	edges []EdgeViewModel
}
