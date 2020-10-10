package draft

type Draft struct {
	Models map[string]ModelOptions `yaml:"Models"`
}

type ModelOptions struct {
	Properties     map[string]string   `yaml:"Properties"`
	Validation     map[string][]string `yaml:"Validation"`
	Actions        []string            `yaml:"Actions"`
	Relationships  map[string]string   `yaml:"Relationships"`
	SkipController bool                `yaml:"SkipController"`
	Public         bool                `yaml:"Public"`
}
