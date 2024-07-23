package main

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/117503445/goutils"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog/log"

	// "github.com/knadh/koanf/parsers/toml"

	"github.com/knadh/koanf/parsers/toml/v2"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/providers/structs"
	flag "github.com/spf13/pflag"
)

type Config struct {
	Name string `koanf:"name"`
	Age  int    `koanf:"age"`
}

var config *Config

const DEFAULT_CONFIG = "config.toml"

func pathIsFile(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	if fileInfo.IsDir() {
		log.Warn().Str("path", path).Msg("is a directory")
		return false
	}

	return true
}

// config 是指针
func loadConfig(config interface{}) {
	v := reflect.ValueOf(config)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	} else {
		log.Fatal().Msg("config must be a pointer")
	}

	// 读取 结构体字段 标签
	t := reflect.TypeOf(v)
	log.Debug().Str("configType", t.String()).Msg("configType")

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		log.Debug().Str("field", field.Name).Str("type", field.Type.String()).Msg("field")
	}

	// koanf instance. Use "." as the key path delimiter. This can be "/" or any character.
	var k = koanf.New(".")

	if err := k.Load(structs.Provider(config, "koanf"), nil); err != nil {
		log.Fatal().Err(err).Msg("error loading default config")
	} else {
		log.Debug().Interface("config", k.All()).Msg("loading default config")
	}

	// Load environment variables.
	if err := k.Load(env.Provider("", ".", func(s string) string {
		allow := []string{"NAME", "AGE", "CONFIG"} // TODO: read from struct tags
		for _, a := range allow {
			if s == a {
				if s == "CONFIG" {
					return "config"
				}
				return s
			}
		}
		return ""
	}), nil); err != nil {
		log.Fatal().Err(err).Msg("error loading env vars")
	} else {
		log.Debug().Interface("config", k.All()).Msg("loading env vars")
	}

	// Use the POSIX compliant pflag lib instead of Go's flag lib.
	f := flag.NewFlagSet("config", flag.ContinueOnError)
	f.Usage = func() {
		fmt.Print(f.FlagUsages())
		os.Exit(0)
	}
	// Path to one or more config files to load into koanf along with some config params.
	f.StringSlice("config", []string{DEFAULT_CONFIG}, "path to one or more .toml config files")

	f.Parse(os.Args[1:])

	// Load the config files provided in the commandline.
	cFiles, err := f.GetStringSlice("config")
	if err != nil {
		log.Fatal().Err(err).Msg("error getting config files")
	}

	cFilesWithCli := false
	for _, arg := range os.Args {
		if arg == "--config" {
			cFilesWithCli = true
			break
		}
	}
	if !cFilesWithCli {
		if k.Get("config") != nil {
			cFiles = strings.Split(k.String("config"), ",")
			log.Debug().Strs("config", cFiles).Msg("config files from env")
		} else {
			log.Debug().Strs("config", cFiles).Msg("config files from default")
		}
	}

	for _, c := range cFiles {

		if !path.IsAbs(c) {
			if workingDirectory, err := os.Getwd(); err == nil {
				if pathIsFile(path.Join(workingDirectory, c)) {
					c = path.Join(workingDirectory, c)
				}
			} else {
				log.Warn().Err(err).Str("path", c).Msg("error getting working directory")
			}

			if executablePath, err := os.Executable(); err == nil {
				if pathIsFile(path.Join(path.Dir(executablePath), c)) {
					c = path.Join(path.Dir(executablePath), c)
				}
			} else {
				log.Warn().Err(err).Str("path", c).Msg("error getting executable path")
			}
		}

		if !pathIsFile(c) {
			if c == DEFAULT_CONFIG {
				log.Info().Str("file", c).Msg("config file not found")
			} else {
				log.Warn().Str("file", c).Msg("config file not found")
			}
			continue
		}

		if err := k.Load(file.Provider(c), toml.Parser()); err != nil {
			log.Fatal().Err(err).Str("file", c).Msg("error loading file")
		} else {
			log.Debug().Str("file", c).Interface("config", k.All()).Msg("loading config file")
		}

	}

	// "time" and "type" may have been loaded from the config file, but
	// they can still be overridden with the values from the command line.
	// The bundled posflag.Provider takes a flagset from the spf13/pflag lib.
	// Passing the Koanf instance to posflag helps it deal with default command
	// line flag values that are not present in conf maps from previously loaded
	// providers.
	if err := k.Load(posflag.Provider(f, ".", k), nil); err != nil {
		log.Fatal().Err(err).Msg("error loading config from cli")
	}

	// Unmarshal the loaded config into the conf struct.
	if err := k.Unmarshal("", config); err != nil {
		log.Fatal().Err(err).Msg("error unmarshaling config")
	}
}

func main() {

	goutils.InitZeroLog()

	config = &Config{
		Name: "default-name",
		Age:  18,
	}

	loadConfig(config)

	log.Info().Interface("config", config).Msg("config loaded")

}
