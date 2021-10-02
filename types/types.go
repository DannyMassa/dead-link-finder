package types

type Config struct {
	FileExtensions    []string `yaml:"fileExtensions"`
	Directories       []string `yaml:"directories"`
	Ignored           []string `yaml:"ignored"`
	SuccessLogs       bool     `yaml:"successLogs"`
	MaxFailures       int      `yaml:"maxFailures"`
	GoldenURL         string   `yaml:"goldenURL"`
	IndividualTimeout int      `yaml:"individualTimeout"`
}

type URL struct {
	Link       string
	Result     string
	File       string
	Directory  string
}
