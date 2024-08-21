package main

import (
	"fmt"

	"github.com/alecthomas/kong"
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

var cli struct {
	Debug bool `help:"Enable debug mode."`

	Rm RmCmd `cmd:"" help:"Remove files."`
	Ls LsCmd `cmd:"" help:"List paths."`
}

func main() {
	ctx := kong.Parse(&cli)
	// Call the Run() method of the selected parsed command.
	err := ctx.Run(&Context{Debug: cli.Debug})
	ctx.FatalIfErrorf(err)
}
