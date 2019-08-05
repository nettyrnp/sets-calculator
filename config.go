package main

import (
	"github.com/namsral/flag"
	"github.com/pkg/errors"
	"os"
	"strings"
	"sync"
)

var flagErrorHandling = flag.ContinueOnError

type Config struct {
	Folder string `json:"folder"`
	Task   string `json:"task"`
}

func (c Config) Validate() []error {
	var errs []error
	if strings.TrimSpace(c.Folder) == "" {
		errs = append(errs, errors.New("Config requires a non-empty folder value"))
	}
	if strings.TrimSpace(c.Task) == "" {
		errs = append(errs, errors.New("Config requires a non-empty task value"))
	}
	return errs
}

func initConfig() (*Config, error) {
	var config Config
	flagset := flag.NewFlagSetWithEnvPrefix("sets-calculator", "", flagErrorHandling)

	flagset.StringVar(&config.Folder, "folder", "files", "Sub-folder with files containing sets of integers.")
	flagset.StringVar(&config.Task, "task", "", "Formulation of the arithmetic task, e.g. `[ SUM [ DIF a.txt b.txt c.txt ] [ INT b.txt c.txt ] ]`.")

	if err := flagset.Parse(os.Args[1:]); err != nil {
		return nil, errors.Wrap(err, "parsing flags")
	}

	// Validate the config.
	if errs := config.Validate(); len(errs) > 0 {
		return nil, errors.Errorf("invalid flag(s): %s", errs)
	}
	return &config, nil
}

var cfg *Config
var once sync.Once

func GetConfig() (*Config, error) {
	err := error(nil)
	if cfg == nil {
		once.Do(func() {
			cfg, err = initConfig()
		})
	}
	return cfg, err
}
