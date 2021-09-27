package main

type Node struct {
	Name   string
	CPU    float64
	Memory float64
	Zone   string
	Cost   float64
}

type Pod struct {
	Name          string
	CPURequest    float64
	MemoryRequest float64
	Zone          string
}

type RequestedCapacity struct {
	CPURequest    float64
	MemoryRequest float64
	Zone          string
	Pods          []Pod
}

type NodeStatsForPod struct {
	Node       Node
	CostForPod float64
	Count      float64
}

type AllocatedNode struct {
	Node       Node
	CPUFree    float64
	MemoryFree float64
	Pods       []Pod
}
