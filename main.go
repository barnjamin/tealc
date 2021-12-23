package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/algorand/go-algorand/data/transactions/logic"
)

var (
	fname       string
	disassemble bool
)

func main() {
	flag.StringVar(&fname, "f", "", "The file to read")
	flag.BoolVar(&disassemble, "d", false, "Disassemble")
	flag.Parse()

	var (
		program []byte
		outname string
	)

	if disassemble {
		program = disassembleFile(fname)
		outname = fname + ".teal"
	} else {
		program = assembleFile(fname)
		outname = fname + ".tok"
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
