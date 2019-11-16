package models

type Vertex struct {
	Id       int
	PageRank float64
	Edges    []Edge
}

type Edge struct {
	Src_id   int
	Dest_id  int
	PageRank float64
}

type InLinks struct {
	Id           int
	SumPageRanks float64
}

type Line struct {
	Id  int
	Out []int
}
