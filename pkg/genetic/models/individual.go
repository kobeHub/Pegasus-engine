/*
`models` package defaines the nessary struct to implements
genetic algorithm. including:
   + Individual
   + Task
   + Node
   + Resources
*/
package models

// The item of genetic algorithm, one individual is
// is one entry during the evolution
type Individual struct {
	ID string
	// The taskId as key and nodeId as value, one `Assignment`
	// is a schedule scheme
	Assignment map[string]string
	// The pods layout before rescheduling
	OrignalAssignment    map[string]string
	AllNodes             map[string]Node
	AllTasks             map[string]Task
	NumOfUnassignedTasks int
	NumofUnassignedNodes int

	// Objective valus, all the objective values to be optimized
	ObjectiveValues           []float64
	TranslatedObjectiveValues []float64
	NormalizedObjectiveValues []float64
	PerpendicularDistance     float64
	ReferencePoint            ReferencePoint

	// NSGA-II
	IndividualsDominatedByThis   []*Individual
	NumOfIndividualsDominateThis int
	Rank                         int
	CrowdingDistance             float64

	// NSGA-III
	ConstrainedViolationValue float64
	IsFeasible                bool
}

// The initial state of current cluster before rescheduling
// init the individual
//
// every `Node` contains tasks list and resources usage state
// `num_tasks` is the number of current pods
func (info *Individual) init(nodes []Node, num_tasks int) {
	info.OrignalAssignment = make(map[string]string, num_tasks)
	info.AllTasks = make(map[string]Task, num_tasks)
	info.AllNodes = make(map[string]Node, len(nodes))

	for _, node := range nodes {
		info.AllNodes[node.ID] = node
		for key, task := range node.Tasks {
			info.AllTasks[key] = task
			info.OrignalAssignment[node.ID] = task.TaskID
		}
	}
}
