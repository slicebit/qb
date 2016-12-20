package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type constant struct {
	Name  string
	Value int
}

func readConstants(filename string, prefix string) []constant {
	var constants []constant
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	splitter := regexp.MustCompile("[ \t]+")
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "#define "+prefix) {
			tokens := splitter.Split(line, -1)
			if tokens[0] != "#define" {
				panic(tokens[0])
			}
			value, err := strconv.Atoi(tokens[2])
			if err != nil {
				os.Stderr.WriteString("Skipped " + line + "\n")
			}
			constants = append(constants, constant{
				Name:  tokens[1],
				Value: value,
			})
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return constants
}

func writeConstant(w io.Writer, c constant) {
	fmt.Fprintf(w, "    %s = %d\n", c.Name, c.Value)
}

func main() {
	f, err := os.OpenFile("errors.go", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString(`package mysql

// Errors defined in errmsg.h
const (
`)

	for _, c := range readConstants("/usr/include/mysql/errmsg.h", "CR_") {
		writeConstant(f, c)
	}
	f.WriteString(`)

// Errors defined in mysqld_error.h
const (
`)
	for _, c := range readConstants("/usr/include/mysql/mysqld_error.h", "ER_") {
		writeConstant(f, c)
	}
	for _, c := range readConstants("/usr/include/mysql/mysqld_error.h", "WARN_") {
		writeConstant(f, c)
	}
	f.WriteString(")\n")
}
