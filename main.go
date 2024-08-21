package main

import (
	"fmt"
	"log"

	"github.com/alecthomas/kong"
	kongtoml "github.com/alecthomas/kong-toml"
)

type Context struct {
	Debug bool
}

type RmCmd struct {
	Force     bool `help:"Force removal."`
	Recursive bool `help:"Recursively remove files."`

	Paths []string `arg:"" name:"path" help:"Paths to remove." type:"path"`
}

func (r *RmCmd) Run(ctx *Context) error {
	fmt.Println("rm", r.Paths)
	return nil
}

type LsCmd struct {
	Paths []string `arg:"" optional:"" name:"path" help:"Paths to list." type:"path"`
}

func (l *LsCmd) Run(ctx *Context) error {
	fmt.Println("ls", l.Paths)
	return nil
}

type DefaultCmd struct {
}

func (d *DefaultCmd) Run(ctx *Context) error {
	// fmt.Println("default, debug:", ctx.Debug, cli.Name)
	fmt.Println(cli.Name)
	return nil
}

var cli struct {
	Debug bool `help:"Enable debug mode." env:"DEBUG"`

	Name string `help:"Name of the user." env:"NAME"`

	Config []string `short:"c" help:"Config files." type:"path" default:"config.toml" env:"CONFIG"`

	Rm         RmCmd      `cmd:"" help:"Remove files."`
	Ls         LsCmd      `cmd:"" help:"List paths."`
	DefaultCmd DefaultCmd `cmd:"" hidden:"" default:"1"`
}

func main() {
	kong.Parse(&cli)
	log.Printf("config: %v", cli.Config)
	ctx := kong.Parse(&cli, kong.Configuration(kongtoml.Loader, cli.Config...))
	// ctx := kong.Parse(&cli)
	// Call the Run() method of the selected parsed command.
	err := ctx.Run(&Context{Debug: cli.Debug})
	ctx.FatalIfErrorf(err)
}

// 默认 cmd
