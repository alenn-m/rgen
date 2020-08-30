package draft

type Draft struct {
	Models map[string]ModelOptions `yaml:"Models"`
}

type ModelOptions struct {
	Properties     map[string]string   `yaml:"Properties"`
	SkipController bool                `yaml:"SkipController"`
	Validation     map[string][]string `yaml:"Validation"`
	Actions        []string            `yaml:"Actions"`
	Relationships  map[string]string   `yaml:"Relationships"`
}
