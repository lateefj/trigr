package trigr

// Process represents a step in the entire pipeline
type Process struct {
	Name string
}

// Pipeline represents the entire development process
type Pipeline struct {
	Configuration map[string]Process
	Prepare       map[string]Process
	Build         map[string]Process
	Package       map[string]Process
	Deploy        map[string]Process
	Running       []Process
}

func NewPipeline() *Pipeline {
	return &Pipeline{
		Configuration: make(map[string]Process),
		Prepare:       make(map[string]Process),
		Build:         make(map[string]Process),
		Package:       make(map[string]Process),
		Deploy:        make(map[string]Process),
		Running:       make([]Process, 0),
	}
}
