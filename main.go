package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/tvarney/maputil/errctx"
	"github.com/tvarney/sdtdmod/pkg/node"
	"github.com/tvarney/sdtdmod/pkg/node/load"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	os.Exit(run())
}

func run() int {
	debug := kingpin.Flag("debug", "enable debug output").Short('D').Bool()
	cnfgfile := kingpin.Flag("config", "path to the configuration file").Short('c').Default("./config.json").String()

	_ = kingpin.Command("validate", "validate the configuration file")
	apply := kingpin.Command("apply", "apply the configuration to the data")
	applydir := apply.Arg("xmldir", "the directory containing XML files to update").Required().String()
	dryrun := kingpin.Command("dry-run", "display changes that would be made")
	dryrundir := dryrun.Arg("xmldir", "the directory continaing XML files to update").Required().String()

	cmd := kingpin.Parse()
	if *debug {
		log.SetOutput(os.Stderr)
	} else {
		log.SetOutput(io.Discard)
	}

	log.Printf("Loading config file %q", *cnfgfile)
	n, err := load.LoadFile(*cnfgfile, &errctx.ErrorPrinter{Stream: os.Stderr})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		return 1
	}

	switch cmd {
	case "validate":
		return Validate(n)
	case "apply":
		return Apply(n, false, *applydir)
	case "dry-run":
		return Apply(n, true, *dryrundir)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %q", cmd)
		return -1
	}
}

func Apply(config []*node.Node, dryrun bool, dir string) int {
	return 0
}

func Validate(config []*node.Node) int {
	switch len(config) {
	case 0:
		fmt.Fprintf(os.Stdout, "No config nodes loaded\n")
	case 1:
		d, _ := json.MarshalIndent(config[0].Serialize(), "", "  ")
		fmt.Fprintf(os.Stdout, "Config:\n%s\n", string(d))
	default:
		nodes := make([]map[string]interface{}, 0, len(config))
		for _, node := range config {
			nodes = append(nodes, node.Serialize())
		}
		d, _ := json.MarshalIndent(nodes, "", "  ")
		fmt.Fprintf(os.Stdout, "Config:\n%s\n", string(d))
	}
	return 0
}
