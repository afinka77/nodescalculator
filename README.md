## nodescalculator
## Selecting optimal new nodes for unscheduled pods written in Go

@author Alina Gončarko 

#### Problem description

Given a Kubernetes cluster that is currently full, and there is currently a number of unscheduled pods. 
The cluster needs to be increased by adding new nodes to accommodate these pods.

Unscheduled pods are represented like this:
```
[ 
 { 
 "name": "pod1", 
 "cpuRequest": 0.4, 
 "memoryRequest": 0.5, 
 "zone": "B" 
 }, 
 { 
 "name": "pod2", 
 "cpuRequest": 1, 
 "memoryRequest": 0.5, 
 "zone": null 
 }, 
 ..... 
]
```

Pods will always have requested amount of CPU and memory; they may or may not require a specific zone ("A" or "B") - if zone is not provided,  pod can go into nodes on either zone.
It selects optimal nodes to add to Kubernetes cluster so that the pods get enough capacity to get scheduled. It selects from  available node types, which are defined like this:

```
[ 
 { 
 "name": "node_2_4_A", 
 "cpu": 2, 
 "memory": 4, 
 "zone": "A", 
 "cost": 9.415777348576905 
 }, 
 { 
 "name": "node_2_4_B", 
 "cpu": 2, 
 "memory": 4, 
 "zone": "B", 
 "cost": 5.050512748885945 
 }, 
... 
]
```

#### Requirements

* Each node type has amount of CPU and Memory, it’s network zone, and cost associated with it.

* Finds cheapest configuration of nodes that fits provided pod list. 

* Node type can be reused, e.g. your result can contain two nodes of type node_2_4_B and three nodes of type node_2_4_A;

* There is at least one way to fit all pods into these nodes, so that totals of requests do not exceed node capacity, e.g.: 

```
    Node 1 (node_2_4_a): pod1 (requires 1 cpu, 2 ram, requires zone A), pod2 (1 cpu, 2 ram)
    Node 2 (node_2_4_b): pod3 (requires 0.5cpu, 1ram), pod4 (0.5cpu, 1ram, zone B), pod5 (1 cpu, 2 ram)
```
* Pods that require a specific zone, must have a node for that zone;

* It is the cheapest possible configuration. There are multiple ways to fit the pods into nodes (e.g. two smaller nodes or one bigger node?),  it have to select cheapest one.

#### Result

Prints the resulting configuration like this:

```
node_2_4_B: pod1, pod99
node_2_4_B: pod2, pod3
node_2_32_B: pod5, pod7, pod85, pod33, pod55
...
```

It also pints out a list of node types and final price for selected configuration.
As well various additional stats.