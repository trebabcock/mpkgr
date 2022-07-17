package format

type Package struct {
	Name        string            `yaml:"name"`
	Version     string            `yaml:"version"`
	Release     int               `yaml:"release"`
	Source      map[string]string `yaml:"source"`
	License     []string          `yaml:"license"`
	Summary     string            `yaml:"summary"`
	Group       string            `yaml:"group"`
	Description string            `yaml:"description"`
	Builddeps   []string          `yaml:"builddeps"`
	Rundeps     []string          `yaml:"rundeps"`
	Homepage    string            `yaml:"homepage"`
	Setup       string            `yaml:"setup"`
	Build       string            `yaml:"build"`
	Install     string            `yaml:"install"`
}
