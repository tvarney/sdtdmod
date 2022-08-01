package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/tvarney/maputil/errctx"
	"github.com/tvarney/sdtdmod/pkg/node/load"
)

func main() {
	os.Exit(run())
}

func run() int {
	log.SetOutput(os.Stdout)
	log.Printf("Loading ./config.json")

	n, err := load.LoadFile("config.json", &errctx.ErrorPrinter{Stream: os.Stderr})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return 1
	}
	switch len(n) {
	case 0:
		fmt.Fprintf(os.Stdout, "No config nodes loaded\n")
	case 1:
		d, _ := json.MarshalIndent(n[0].Serialize(), "", "  ")
		fmt.Fprintf(os.Stdout, "Config:\n%s\n", string(d))
	default:
		nodes := make([]map[string]interface{}, 0, len(n))
		for _, node := range n {
			nodes = append(nodes, node.Serialize())
		}
		d, _ := json.MarshalIndent(nodes, "", "  ")
		fmt.Fprintf(os.Stdout, "Config:\n%s\n", string(d))
	}
	return 0
}
