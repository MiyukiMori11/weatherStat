package domain

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

type Service struct {
	Name       string  `yaml:"name"`
	Root       string  `yaml:"root"`
	Host       string  `yaml:"host"`
	Port       string  `yaml:"port"`
	Scheme     string  `yaml:"scheme"`
	HealthPath string  `yaml:"healthPath"`
	Routes     []Route `yaml:"routes"`
}

type Route struct {
	Path    string   `yaml:"path"`
	Method  string   `yaml:"method"`
	Payload []string `yaml:"payload,omitempty"`
	Params  []string `yaml:"params,omitempty"`
}
