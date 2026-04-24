package main

import (
	"fmt"
	"os"
)

func logf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage:")
	fmt.Fprintln(os.Stderr, "  bfilt build <fpr> <input_collection> <output_path>")
	fmt.Fprintln(os.Stderr, "  bfilt query <bloom_file> <query_file>")
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	switch os.Args[1] {
	case "build":
		if len(os.Args) != 5 {
			usage()
			os.Exit(2)
		}
		if err := runBuild(os.Args[2], os.Args[3], os.Args[4]); err != nil {
			logf("build: %v\n", err)
			os.Exit(1)
		}
	case "query":
		if len(os.Args) != 4 {
			usage()
			os.Exit(2)
		}
		if err := runQuery(os.Args[2], os.Args[3]); err != nil {
			logf("query: %v\n", err)
			os.Exit(1)
		}
	default:
		usage()
		os.Exit(2)
	}
}
