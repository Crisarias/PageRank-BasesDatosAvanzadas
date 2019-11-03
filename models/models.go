package models

type Node struct {
	pageRank float64
	uid      int
	outLinks map[int]Node
	inLinks  map[int]Node
}
