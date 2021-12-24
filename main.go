package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/algorand/go-algorand/data/transactions/logic"
)

var (
	fname       string
	outname     string
	disassemble bool
	buildMap    bool
)

func main() {
	flag.StringVar(&fname, "f", "", "The file to read")
	flag.StringVar(&outname, "o", "", "The file to read")
	flag.BoolVar(&disassemble, "d", false, "Disassemble")
	flag.BoolVar(&buildMap, "m", false, "Build Map")
	flag.Parse()

	var (
		program []byte
		outname string

		am logic.AssemblyMap
	)

	if disassemble {
		program = disassembleFile(fname)
		if outname == "" {
			outname = fname + ".dis"
		}
	} else if buildMap {
		program, am = assembleWithFileMap(fname)
		if outname == "" {
			outname = fname + ".tok"
		}

		f, err := os.Create(outname + ".map.json")
		if err != nil {
			log.Fatalf("Failed to create assembly map file: %s", err)
		}
		b, err := json.Marshal(am)
		if err != nil {
			log.Fatalf("Failed to marshal assembly map json: %s", err)
		}
		f.Write(b)
		f.Close()

	} else {
		program = assembleFile(fname)
		if outname == "" {
			outname = fname + ".tok"
		}
	}

	f, err := os.Create(outname)
	if err != nil {
		log.Fatalf("Failed to create file for writing: %+v", err)
	}

	f.Write(program)
	f.Close()
}

func assembleFile(fname string) (program []byte) {
	b, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatalf("Failed to read file: %+v", err)
	}

	ops, err := logic.AssembleString(string(b))
	if err != nil {
		log.Fatalf("Failed to assemble: %+v", err)
	}

	return ops.Program
}

func disassembleFile(fname string) (disassembled []byte) {
	program, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatalf("Failed to read file: %s", err)
	}

	text, err := logic.Disassemble(program)
	if err != nil {
		log.Fatalf("Failed to disassemble: %s", err)
	}

	return []byte(text)
}

func assembleWithFileMap(fname string) (program []byte, deets logic.AssemblyMap) {
	text, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatalf("Failed to read file: %s", err)
	}

	ops, err := logic.AssembleString(string(text))
	if err != nil {
		ops.ReportProblems(fname, os.Stderr)
		log.Fatalf("Failed to assemble string: %s", err)
	}

	am := ops.GetAssemblyMap()
	am.SourceName = fname

	return ops.Program, am
}
