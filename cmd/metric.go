package main

type NetworkMetric struct {
	NodeCount     int     `json:"node_count"`
	Broadcast     string  `json:"broadcast"`
	AvgDegree     float64 `json:"avg_degree"`
	Delay         int     `json:"delay"`
	DuplicateRate float64 `json:"duplicate_rate"`
	ReceivingRate float64 `json:"receiving_rate"`
}
