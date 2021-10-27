package services

import (
	"fmt"
	"io/ioutil"

	"github.com/DannyMassa/dead-link-linter/types"
	flag "github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

type ParameterService interface {
	SetConfig() *types.Config
}

type ParameterServiceImpl struct{}

func (p ParameterServiceImpl) SetConfig() *types.Config {
	// Initialize config object and set defaults
	defaultConfig := p.parseDefaultConfiguration()

	// Overwrite with .deadlink file values
	fileConfig := p.parseFileConfiguration()
	defaultConfig = p.rectifyConfigurations(fileConfig, defaultConfig)

	// Overwrite with CLI flags
	cliConfig := p.parseCLIConfiguration()
	defaultConfig = p.rectifyConfigurations(cliConfig, defaultConfig)

	return &defaultConfig
}

func (p ParameterServiceImpl) rectifyConfigurations(priorityConfig types.Config, otherConfig types.Config) types.Config { //nolint
	if len(priorityConfig.Directories) > 0 {
		otherConfig.Directories = priorityConfig.Directories
	}
	if len(priorityConfig.FileExtensions) > 0 {
		otherConfig.FileExtensions = priorityConfig.FileExtensions
	}
	if priorityConfig.GoldenURL != "" {
		otherConfig.GoldenURL = priorityConfig.GoldenURL
	}
	if len(priorityConfig.Ignored) > 0 {
		otherConfig.Ignored = priorityConfig.Ignored
	}
	if priorityConfig.IndividualTimeout != -1 {
		otherConfig.IndividualTimeout = priorityConfig.IndividualTimeout
	}
	if priorityConfig.MaxFailures != -1 {
		otherConfig.MaxFailures = priorityConfig.MaxFailures
	}
	if priorityConfig.LogVerbosity != -1 {
		otherConfig.LogVerbosity = priorityConfig.LogVerbosity
	}

	return otherConfig
}

func (p ParameterServiceImpl) parseCLIConfiguration() types.Config {
	cliConfig := types.NewConfig()
	flag.StringArrayVar(&cliConfig.Directories, "directories", []string{}, "Fill in Later")
	flag.StringArrayVar(&cliConfig.FileExtensions, "file_extensions ", []string{}, "Fill in Later")
	flag.StringVar(&cliConfig.GoldenURL, "golden_url", "", "Fill in Later")
	flag.StringArrayVar(&cliConfig.Ignored, "ignored", []string{}, "Fill in Later")
	flag.IntVar(&cliConfig.IndividualTimeout, "individual_timeout", -1, "")
	flag.IntVar(&cliConfig.MaxFailures, "max_failures", -1, "")
	flag.IntVar(&cliConfig.LogVerbosity, "log_verbosity", -1, "")
	flag.Parse()
	return cliConfig
}

func (p ParameterServiceImpl) parseDefaultConfiguration() types.Config {
	return types.Config{
		Directories: []string{"./"},
		FileExtensions: []string{".markdown", ".mdown", ".mkdn", ".md", ".mkd", ".mdwn", ".mdtxt", ".mdtext", ".text",
			".txt", ".rmd", ".rst"},
		GoldenURL:         "https://google.com",
		Ignored:           []string{},
		IndividualTimeout: 10,
		MaxFailures:       0,
		LogVerbosity:      1,
	}
}

func (p ParameterServiceImpl) parseFileConfiguration() types.Config {
	fileConfig := types.NewConfig()
	buf, err := ioutil.ReadFile(".deadlink")

	if err != nil {
		fmt.Printf("Could not find .deadlink file\n")
	}
	err = yaml.Unmarshal(buf, &fileConfig)
	if err != nil {
		panic("Could not parse .deadlink file")
	}

	return fileConfig
}
