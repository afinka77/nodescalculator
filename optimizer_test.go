package main

import "testing"

func TestGetCheapestFreeNode(t *testing.T) {
	t.Run("cheapest node is selected", func(t *testing.T) {
		nodes := []Node{Node{"node_8_4_A1", 8, 4, "A", 200},
			Node{"node_8_4_A2", 8, 4, "A", 100}}
		pods := []Pod{Pod{"pod1", 1, 1, "A"},
			Pod{"pod2", 1, 1, "A"},
			Pod{"pod3", 1, 1, "A"},
			Pod{"pod4", 1, 1, "A"}}
		capacity := RequestedCapacity{1, 1, "A", pods}
		fitPodsCount, bestNode := getCheapestFreeNode(nodes, capacity)
		if fitPodsCount != 4 {
			t.Fatalf("fitPodsCount %v should be equal 4", fitPodsCount)
		}

		if bestNode != nodes[1] {
			t.Fatalf("Wrond best node selected")
		}
	})

	t.Run("when pods zone is specified nodes with matching zone are selected", func(t *testing.T) {
		nodes := []Node{Node{"node_8_4_A1", 8, 4, "A", 200},
			Node{"node_8_4_A2", 8, 4, "B", 100}}
		pods := []Pod{Pod{"pod1", 1, 1, "A"}}
		capacity := RequestedCapacity{1, 1, "A", pods}
		_, bestNode := getCheapestFreeNode(nodes, capacity)

		if bestNode != nodes[0] {
			t.Fatalf("Wrond zone node is selected")
		}
	})

	t.Run("when pods zone is not specified cheapest node with any zone is selected", func(t *testing.T) {
		nodes := []Node{Node{"node_8_4_A1", 8, 4, "A", 200},
			Node{"node_8_4_A2", 8, 4, "B", 100}}
		pods := []Pod{Pod{"pod1", 1, 1, ""}}
		capacity := RequestedCapacity{1, 1, "", pods}
		_, bestNode := getCheapestFreeNode(nodes, capacity)

		if bestNode != nodes[1] {
			t.Fatalf("Wrond node is selected")
		}
	})

	t.Run("when allocation pods count is bigger then fit pods count, fit pods count is returned", func(t *testing.T) {
		nodes := []Node{Node{"node_8_4_A1", 8, 4, "A", 200},
			Node{"node_8_4_A2", 8, 4, "B", 100}}
		pods := []Pod{Pod{"pod1", 1, 1, ""},
			Pod{"pod2", 1, 1, ""},
			Pod{"pod3", 1, 1, ""},
			Pod{"pod4", 1, 1, ""},
			Pod{"pod5", 1, 1, ""}}
		capacity := RequestedCapacity{1, 1, "", pods}
		fitPodsCount, _ := getCheapestFreeNode(nodes, capacity)

		if fitPodsCount != 4 {
			t.Fatalf("Fit nodes count %v should be equal 4", fitPodsCount)
		}
	})

}

func TestAllocateFreeNodes(t *testing.T) {
	var allocatedNodes []AllocatedNode

	t.Run("allocate 5 (4 and 1) pods on two nodes", func(t *testing.T) {
		nodes := []Node{Node{"node_8_4_A1", 8, 4, "A", 200},
			Node{"node_8_4_A2", 8, 4, "B", 100}}
		pods := []Pod{Pod{"pod1", 1, 1, ""},
			Pod{"pod2", 1, 1, ""},
			Pod{"pod3", 1, 1, ""},
			Pod{"pod4", 1, 1, ""},
			Pod{"pod5", 1, 1, ""}}
		capacity := RequestedCapacity{1, 1, "", pods}
		allocatedNodes := allocateFreeNodes(nodes, capacity, allocatedNodes)

		if len(allocatedNodes) != 2 {
			t.Fatalf("Shold be allocated 2 nodes")
		}

		if allocatedNodes[0].CPUFree != 4 {
			t.Fatalf("Should be 4 free CPU on first node")
		}

		if allocatedNodes[0].MemoryFree != 0 {
			t.Fatalf("Should be 0 free Memory on first node")
		}

		if len(allocatedNodes[0].Pods) != 4 {
			t.Fatalf("Shold be 4 pods on first node")
		}

		if allocatedNodes[1].CPUFree != 7 {
			t.Fatalf("Should be 7 free CPU on second node")
		}

		if allocatedNodes[1].MemoryFree != 3 {
			t.Fatalf("Should be 3 free Memory on second node")
		}

		if len(allocatedNodes[1].Pods) != 1 {
			t.Fatalf("Shold be 1 pod on the second node")
		}
	})
}

func TestAllocatePods(t *testing.T) {
	t.Run("5 pods (4 and 1) are allocated on 2 nodes, then 3 pods are allocated on same nodes having free space", func(t *testing.T) {
		nodes := []Node{Node{"node_8_4_A1", 8, 4, "A", 200},
			Node{"node_8_4_A2", 8, 4, "B", 100}}
		pods1 := []Pod{Pod{"pod1", 1, 1, ""},
			Pod{"pod2", 1, 1, ""},
			Pod{"pod3", 1, 1, ""},
			Pod{"pod4", 1, 1, ""},
			Pod{"pod5", 1, 1, ""}}
		pods2 := []Pod{Pod{"pod6", 1, 1, ""},
			Pod{"pod7", 1, 1, ""},
			Pod{"pod8", 1, 1, ""}}
		capacity := []RequestedCapacity{RequestedCapacity{1, 1, "", pods1},
			RequestedCapacity{1, 1, "", pods2}}
		allocatedNodes := AllocatePods(nodes, capacity)

		if len(allocatedNodes) != 2 {
			t.Fatalf("Shold be allocated 2 nodes")
		}

		if len(allocatedNodes[0].Pods) != 4 {
			t.Fatalf("Shold be 4 pods on first node")
		}

		if allocatedNodes[0].CPUFree != 4 {
			t.Fatalf("Should be 4 free CPU on first node")
		}

		if allocatedNodes[0].MemoryFree != 0 {
			t.Fatalf("Should be 0 free Memory on first node")
		}

		if len(allocatedNodes[1].Pods) != 4 {
			t.Fatalf("Shold be 4 pods on the second node")
		}

		if allocatedNodes[1].CPUFree != 4 {
			t.Fatalf("Should be 4 free CPU on second node")
		}

		if allocatedNodes[1].MemoryFree != 0 {
			t.Fatalf("Should be 0 free Memory on second node")
		}
	})
}
