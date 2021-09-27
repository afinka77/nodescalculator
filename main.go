package main

import (
	"fmt"
)

var (
	nodesJsonPath = "data/nodes.json"
	podsJsonPath  = "data/pods.json"
)

func main() {
	var nodes []Node
	JsonToStruct(nodesJsonPath, &nodes)

	var pods []Pod
	JsonToStruct(podsJsonPath, &pods)
	printPodsStats(pods)

	capacityMap := GroupPodsByRequestedCapacity(pods)
	printCapacityMap(capacityMap)

	allocatedNodes := AllocatePods(nodes, capacityMap)
	//here is main task result
	printAllocatedNodesStats(allocatedNodes)
	printDetailedPodsAllocationStats(allocatedNodes)
}

func printCapacityMap(capacityMap []RequestedCapacity) {
	fmt.Printf("\n\nPods grouped by requested capacity and zones:")
	podsCountTotal := 0
	for _, capacity := range capacityMap {
		podsCountTotal += len(capacity.Pods)
		fmt.Printf("\n %.1f requested CPU %.1f requested Memory %1v zone: %v pods", capacity.CPURequest, capacity.MemoryRequest, capacity.Zone, len(capacity.Pods))
	}
	fmt.Printf("\nTOTAL %v pods", podsCountTotal)
}

func printPodsStats(pods []Pod) {
	totalCPURequest := 0.0
	totalMemoryRequest := 0.0
	for _, pod := range pods {
		totalCPURequest += pod.CPURequest
		totalMemoryRequest += pod.MemoryRequest
	}
	fmt.Printf("\n\nPods count: %v \nTotal pods CPU request %.2f \nTotal pods memory request %.2f", len(pods), totalCPURequest, totalMemoryRequest)
}

func printAllocatedNodesStats(allocatedNodes []AllocatedNode) {
	fmt.Printf("\n\nAllocated nodes:")
	totalCost := 0.0
	totalCPU := 0.0
	freeCPU := 0.0
	totalMemory := 0.0
	freeMemory := 0.0
	totalPods := 0
	for _, node := range allocatedNodes {
		totalCost += node.Node.Cost
		totalCPU += node.Node.CPU
		freeCPU += node.CPUFree
		totalMemory += node.Node.Memory
		freeMemory += node.MemoryFree
		totalPods += len(node.Pods)
		podNames := ""
		for i, pod := range node.Pods {
			if i > 0 {
				podNames += ", "
			}
			podNames += pod.Name
		}
		fmt.Printf("\n%-11v (%-2v pods): %v ", node.Node.Name, len(node.Pods), podNames)
	}
	fmt.Printf("\n\nAllocated Nodes Count: %v \nTotalCost: %.2f \nTotalCPU %.2f \nFreeCPU (%.2f %%) %.2f "+
		"\nTotalMemory %.2f \nFreeMemory (%.2f %%) %.2f \nTotalPods %v",
		len(allocatedNodes), totalCost, totalCPU, freeCPU/totalCPU*100, freeCPU,
		totalMemory, freeMemory/totalMemory*100, freeMemory, totalPods)
}

func printDetailedPodsAllocationStats(allocatedNodes []AllocatedNode) {
	fmt.Printf("\n\nDetailed stats by nodes:")
	totalNodesCPU := 0.0
	totalNodesMemory := 0.0
	totalNodesCPUUsed := 0.0
	totalNodesMemoryUsed := 0.0
	for _, aNode := range allocatedNodes {
		fmt.Printf("\nNode %v zone %v cpu %.2f (%.2f free) memory %.2f (%.2f free)", aNode.Node.Name, aNode.Node.Zone, aNode.Node.CPU, aNode.CPUFree, aNode.Node.Memory, aNode.MemoryFree)
		totalCPU := 0.0
		totalMemory := 0.0
		podsZones := ""
		for _, pod := range aNode.Pods {
			totalCPU += pod.CPURequest
			totalMemory += pod.MemoryRequest
			podsZones += pod.Zone
		}
		totalNodesCPU += aNode.Node.CPU
		totalNodesCPUUsed += totalCPU
		totalNodesMemory += aNode.Node.Memory
		totalNodesMemoryUsed += totalMemory
		fmt.Printf("\n%v pods with zones %v (%.2f total/%.2f used/%.2f free) CPU (%.2f total/%.2f used/%.2f free) memory",
			len(aNode.Pods), podsZones, aNode.Node.CPU, totalCPU, aNode.Node.CPU-totalCPU, aNode.Node.Memory,
			totalMemory, aNode.Node.Memory-totalMemory)
	}
	fmt.Printf("\n\nTotal nodes CPU (%.2f total/%.2f used/%.2f free) total nodes memory (%.2f total/%.2f used/%.2f free)", totalNodesCPU, totalNodesCPUUsed, totalNodesCPU-totalNodesCPUUsed, totalNodesMemory, totalNodesMemoryUsed, totalNodesMemory-totalNodesMemoryUsed)
}
