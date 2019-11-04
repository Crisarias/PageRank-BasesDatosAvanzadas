package models

type Vertex struct {
	id       int
	pageRank float64
	edges    []Edge
}

type Edge struct {
	src_id   int
	dest_id  int
	pageRank float64
}
