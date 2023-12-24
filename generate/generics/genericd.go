package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// go:generate go run main.go
func main() {

	n := 10

	x, err := filepath.Abs("generics_gen.go")
	if err != nil {
		fmt.Printf("Error getting absolute path: %s\n", err)
		os.Exit(1)
	}

	file, err := os.Create(x)
	if err != nil {
		fmt.Printf("Error creating file: %s\n", err)
		os.Exit(1)
	}
	defer file.Close()

	fmt.Fprintf(file, "// Code generated by go generate; DO NOT EDIT.\n")
	fmt.Fprintf(file, "package snake\n\n")

	for i := 0; i <= n; i++ {
		genCommandRunGen(file, i)
		genCommandRunWithOutputGen(file, i)
	}

	for i := 0; i <= n; i++ {
		for o := 0; o <= n; o++ {
			genResolverRunWithOutputGen(file, i, o)
		}
	}

	fmt.Println("Done, wrote to ", file.Name())
}
func xy(file *os.File, n, o int, includeAny bool, includeBrackets bool, insert string) {
	if (includeBrackets && n+o != 0) || insert != "" {
		fmt.Fprintf(file, "[")
	}
	for i := 1; i <= n; i++ {
		fmt.Fprintf(file, "X%d", i)
		if includeAny {
			fmt.Fprintf(file, " any")
		}
		if i != n || o != 0 {
			fmt.Fprintf(file, ", ")
		}
	}
	for i := 1; i <= o; i++ {
		fmt.Fprintf(file, "Y%d", i)
		if includeAny {
			fmt.Fprintf(file, " any")
		}
		if i != o {
			fmt.Fprintf(file, ", ")
		}
	}
	if insert != "" {
		if n+o != 0 {
			fmt.Fprintf(file, ", ")
		}
		fmt.Fprintf(file, "L %s", insert)
		xy(file, n, o, false, true, "")
	}
	if (includeBrackets && n+o != 0) || insert != "" {
		fmt.Fprintf(file, "]")
	}
}

func genCommandRunGen(file *os.File, n int) {
	work := fmt.Sprintf("RunCommand_In%.2d_Out01", n)
	fmt.Fprintf(file, "type gen%s", work)
	xy(file, n, 0, true, true, "")
	fmt.Fprintf(file, " interface { NamedMethod; Run (")
	xy(file, n, 0, false, false, "")
	fmt.Fprintf(file, ") error }\n")

	fmt.Fprintf(file, "func Gen%s", work)
	xy(file, n, 0, true, true, fmt.Sprintf("gen%s", work))
	fmt.Fprintf(file, "(l L) TypedNamedRunner[L] { return &namedrund[L]{&rund[L]{l}} }\n\n")
	// xy(file, n, 0, false, true)
	// fmt.Fprintf(file, ") bool { return true }\n\n")
}

func genCommandRunWithOutputGen(file *os.File, n int) {
	work := fmt.Sprintf("RunCommand_In%.2d_Out02", n)
	fmt.Fprintf(file, "type gen%s", work)
	xy(file, n, 0, true, true, "")
	fmt.Fprintf(file, " interface { NamedMethod; Run (")
	xy(file, n, 0, false, false, "")
	fmt.Fprintf(file, ") (Output, error) }\n")

	fmt.Fprintf(file, "func Gen%s", work)
	xy(file, n, 0, true, true, fmt.Sprintf("gen%s", work))
	// fmt.Fprintf(file, "(gen%s", work)
	fmt.Fprintf(file, "(l L) TypedNamedRunner[L] { return &namedrund[L]{&rund[L]{l}} }\n\n")

	// xy(file, n, 0, false, true)
	// fmt.Fprintf(file, ") bool { return true }\n\n")
}

func genResolverRunWithOutputGen(file *os.File, n int, o int) {
	work := fmt.Sprintf("RunResolver_In%.2d_Out%.2d", n, o+1)
	fmt.Fprintf(file, "type gen%s", work)
	xy(file, n, o, true, true, "")
	fmt.Fprintf(file, " interface { Run (")
	xy(file, n, 0, false, false, "")
	fmt.Fprintf(file, ") (")
	xy(file, 0, o, false, false, "")
	if o != 0 {
		fmt.Fprintf(file, ", ")
	}
	fmt.Fprintf(file, "error) }\n")

	fmt.Fprintf(file, "func Gen%s", work)
	xy(file, n, o, true, true, fmt.Sprintf("gen%s", work))
	fmt.Fprintf(file, "(l L) TypedRunner[L] { return &rund[L]{l} }\n\n")
	// xy(file, n, o, false, true)
	// fmt.Fprintf(file, ") bool { return true }\n\n")
}