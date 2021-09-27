package main

import (
	"math"
	"sort"
)

func GroupPodsByRequestedCapacity(pods []Pod) []RequestedCapacity {
	var out []RequestedCapacity
	for _, pod := range pods {
		index, found := findRequestedCapacity(out, pod)
		if found {
			out[index].Pods = append(out[index].Pods, pod)
		} else {
			out = append(out, RequestedCapacity{pod.CPURequest, pod.MemoryRequest, pod.Zone, []Pod{pod}})
		}
	}

	sortByCountAndCapacityDesc(out)

	return out
}

func AllocatePods(nodes []Node, requestedCapacity []RequestedCapacity) []AllocatedNode {
	var allocatedNodes []AllocatedNode

	for i, capacity := range requestedCapacity {
		// allocate on used nodes
		for j, allocatedNode := range allocatedNodes {
			if capacity.Zone == "" || capacity.Zone == allocatedNode.Node.Zone {
				maxPodsFit := math.Floor(math.Min(allocatedNode.CPUFree/capacity.CPURequest, allocatedNode.MemoryFree/capacity.MemoryRequest))
				podsFit := math.Round(math.Min(maxPodsFit, float64(len(capacity.Pods))))
				if podsFit > 0 {
					allocatedNodes[j].CPUFree -= capacity.CPURequest * podsFit
					allocatedNodes[j].MemoryFree -= capacity.MemoryRequest * podsFit
					allocatedNode.Pods = append(allocatedNode.Pods, capacity.Pods[0:int(podsFit)]...)
					allocatedNodes[j].Pods = allocatedNode.Pods
					capacity.Pods = capacity.Pods[int(podsFit):len(capacity.Pods)]
					requestedCapacity[i].Pods = capacity.Pods
				}
			}
		}

		// allocate free nodes
		allocatedNodes = allocateFreeNodes(nodes, capacity, allocatedNodes)
	}

	return allocatedNodes
}

func allocateFreeNodes(nodes []Node, capacity RequestedCapacity, allocatedNodes []AllocatedNode) []AllocatedNode {
	if len(capacity.Pods) > 0 {
		fitPodsCount, bestNode := getCheapestFreeNode(nodes, capacity)
		nodesCount := math.Ceil(float64(len(capacity.Pods)) / fitPodsCount)
		for i := 0; i < int(nodesCount); i++ {
			podsFit := math.Min(fitPodsCount, float64(len(capacity.Pods))-(fitPodsCount*float64(i)))
			allocatedNodes = append(allocatedNodes, AllocatedNode{
				bestNode,
				bestNode.CPU - podsFit*capacity.CPURequest,
				bestNode.Memory - podsFit*capacity.MemoryRequest,
				capacity.Pods[int(fitPodsCount)*i : int(math.Min(fitPodsCount*float64(i+1), float64(len(capacity.Pods))))]})

		}
	}
	return allocatedNodes
}

func getCheapestFreeNode(nodes []Node, requestedCapacity RequestedCapacity) (float64, Node) {
	var out []NodeStatsForPod

	//walk through nodes of fitting zone to calculate cost per pod
	for _, node := range nodes {
		if requestedCapacity.Zone == "" || requestedCapacity.Zone == node.Zone {
			maxPodsFit := math.Floor(math.Min(node.CPU/requestedCapacity.CPURequest, node.Memory/requestedCapacity.MemoryRequest))
			podsCost := node.Cost / maxPodsFit
			podsFit := math.Min(maxPodsFit, float64(len(requestedCapacity.Pods)))
			out = append(out, NodeStatsForPod{node, podsCost, podsFit})
		}
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].CostForPod < out[j].CostForPod
	})

	return out[0].Count, out[0].Node
}

func findRequestedCapacity(out []RequestedCapacity, pod Pod) (index int, found bool) {
	index, found = -1, false
	for i, capacity := range out {
		if pod.CPURequest == capacity.CPURequest && pod.MemoryRequest == capacity.MemoryRequest && pod.Zone == capacity.Zone {
			index, found = i, true
			break
		}
	}
	return index, found
}

func sortByCountAndCapacityDesc(out []RequestedCapacity) {
	sort.Slice(out, func(i, j int) bool {
		if len(out[i].Pods) != len(out[j].Pods) {
			return len(out[i].Pods) > len(out[j].Pods)
		} else if out[i].MemoryRequest != out[j].MemoryRequest {
			return out[i].MemoryRequest > out[i].MemoryRequest
		} else if out[i].CPURequest != out[j].CPURequest {
			return out[i].CPURequest > out[j].CPURequest
		} else {
			return false
		}
	})
}
