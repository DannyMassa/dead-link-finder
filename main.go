package main

import (
	"fmt"
	"github.com/DannyMassa/dead-link-finder/controllers"
	"github.com/DannyMassa/dead-link-finder/services"
	"github.com/DannyMassa/dead-link-finder/types"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"time"
)

var (
	urlService services.URLService = &services.URLServiceImpl{}
)

func main() {
	start := time.Now()

	c := &types.Config{}
	buf, err := ioutil.ReadFile(".deadlink")

	if err != nil {
		fmt.Printf("Could not find .deadlink file, using defaults\n")
		c.FileExtensions = []string{".markdown", ".mdown", ".mkdn", ".md", ".mkd", ".mdwn", ".mdtxt", ".mdtext",
			".text", ".txt", ".rmd", ".rst"}
		c.Directories = []string{"./"}
		c.Ignored = []string{}
	}

	err = yaml.Unmarshal(buf, c)
	if err != nil {
		panic("Could not parse .deadlink file")
	}

	// Testing against golden URL
	if !urlService.LinkLivenessChecker(c.GoldenURL) {
		os.Exit(2)
	}

	controllers.Controller.Run(c)
	err = controllers.Controller.PrintResults(c)

	elapsed := time.Since(start)
	fmt.Printf("\nLinter finished in %s\n", elapsed)
	if err != nil {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
