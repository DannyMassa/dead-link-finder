package types

type Config struct {
	Directories       []string `yaml:"directories"`
	FileExtensions    []string `yaml:"fileExtensions"`
	GoldenURL         string   `yaml:"goldenURL"`
	Ignored           []string `yaml:"ignored"`
	IndividualTimeout int      `yaml:"individualTimeout"`
	MaxFailures       int      `yaml:"maxFailures"`
	LogVerbosity      int      `yaml:"logVerbosity"`
}

func NewConfig() Config {
	return Config{
		Directories:       nil,
		FileExtensions:    nil,
		GoldenURL:         "",
		Ignored:           nil,
		IndividualTimeout: -1,
		MaxFailures:       -1,
		LogVerbosity:      -1,
	}
}

type URL struct {
	Link      string
	Result    string
	File      string
	Directory string
}
