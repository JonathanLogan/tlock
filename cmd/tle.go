package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/drand/tlock"
	"github.com/drand/tlock/cmd/commands"
	"github.com/drand/tlock/networks/http"
)

func main() {
	log := log.New(os.Stderr, "", 0)

	if len(os.Args) == 1 {
		commands.PrintUsage(log)
		return
	}

	if err := run(log); err != nil {
		log.Fatal(err)
	}
}

func run(log *log.Logger) error {
	flags, err := commands.Parse()
	if err != nil {
		return fmt.Errorf("parse commands: %v", err)
	}

	var src io.Reader = os.Stdin
	if name := flag.Arg(0); name != "" && name != "-" {
		f, err := os.OpenFile(name, os.O_RDONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open input file %q: %v", name, err)
		}
		defer f.Close()
		src = f
	}

	var dst io.Writer = os.Stdout
	if name := flags.Output; name != "" && name != "-" {
		f, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return fmt.Errorf("failed to open output file %q: %v", name, err)
		}
		defer f.Close()
		dst = f
	}

	network := http.NewNetwork(flags.Network, flags.Chain)

	switch {
	case flags.Decrypt:
		return tlock.NewDecrypter(network).Decrypt(dst, src)
	default:
		return commands.Encrypt(flags, dst, src, network)
	}
}
