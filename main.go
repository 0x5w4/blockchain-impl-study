package main

import "os"

func main() {
	// Ensure we cleanup database files if we want a fresh start manually
	// os.RemoveAll("./tmp/blocks_node_1")

	defer os.Exit(0)

	cli := CLI{}
	cli.Run()
}
