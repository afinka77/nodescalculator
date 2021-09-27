package main

import (
	"testing"
)

var (
	nodesTestJsonPath = "testdata/nodes.json"
	podsTestJsonPath  = "testdata/pods.json"
)

func TestJsonToStruct(t *testing.T) {
	t.Run("nodes parsed", func(t *testing.T) {
		var nodes []Node
		JsonToStruct(nodesTestJsonPath, &nodes)

		if nodes == nil {
			t.Fatalf("Didn't parse nodes json")
		}

		if len(nodes) != 2 {
			t.Fatalf("Got: %v instead of 2\n", len(nodes))
		}

		if nodes[0].Name != "node_2_4_A" || nodes[1].Name != "node_2_4_B" {
			t.Fatalf("Error when parsing nodes name")
		}

		if nodes[0].CPU != 2 || nodes[1].CPU != 3 {
			t.Fatalf("Error when parsing nodes CPU")
		}

		if nodes[0].Memory != 4 || nodes[1].Memory != 5 {
			t.Fatalf("Error when parsing nodes memory")
		}

		if nodes[0].Zone != "A" || nodes[1].Zone != "B" {
			t.Fatalf("Error when parsing nodes zone")
		}

		if nodes[0].Cost != 9.415777348576905 || nodes[1].Cost != 5.050512748885945 {
			t.Fatalf("Error when parsing nodes cost")
		}
	})

	t.Run("pods parsed", func(t *testing.T) {
		var pods []Pod
		JsonToStruct(podsTestJsonPath, &pods)

		if pods == nil {
			t.Fatalf("Didn't parse pods json")
		}

		if len(pods) != 2 {
			t.Fatalf("Got: %v instead of 2\n", len(pods))
		}

		if pods[0].Name != "pod1" || pods[1].Name != "pod2" {
			t.Fatalf("Error when parsing pods name")
		}

		if pods[0].CPURequest != 0.4 || pods[1].CPURequest != 1 {
			t.Fatalf("Error when parsing pods CPU request")
		}

		if pods[0].MemoryRequest != 0.5 || pods[1].MemoryRequest != 1 {
			t.Fatalf("Error when parsing pods memory request")
		}

		if pods[0].Zone != "B" || pods[1].Zone != "" {
			t.Fatalf("Error when parsing pods zone")
		}
	})

}
