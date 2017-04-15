package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

func toStrings(args []interface{}) []string {
	result := []string{}
	for _, v := range args {
		result = append(result, v.(string))
	}
	return result
}

func getEnvVariable(args ...interface{}) string {
	r := []string{}
	for _, s := range toStrings(args) {
		r = append(r, os.Getenv(s))
	}
	return strings.Join(r, " ")
}

func execCommand(args ...interface{}) string {
	arguments := toStrings(args)
	if len(arguments) <= 0 {
		fmt.Fprintf(os.Stderr, "ERROR: exec found but no command was given")
		os.Exit(1)
	}
	totalCmd := strings.Join(arguments, " ")
	cmd := exec.Command("bash", "-c", totalCmd)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: cmd '%s' failed to execute", totalCmd)
		os.Exit(1)
	}
	return stripEol(out.String())
}

func stripEol(s string) string {
	s1 := strings.TrimSuffix(s, "\r\n")
	if s1 != s {
		return s1
	}
	return strings.TrimSuffix(s1, "\n")
}

func main() {
	var err error

	flag.Parse()

	if len(flag.Args()) <= 0 {
		fmt.Fprintf(os.Stderr, "ERROR: must supply the template file\n")
		os.Exit(1)
	}
	infile := flag.Args()[0]

	funcMap := template.FuncMap{
		"env":  getEnvVariable,
		"exec": execCommand,
	}

	content, err := ioutil.ReadFile(infile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to read file '%s' because '%s'\n", infile, err.Error())
		os.Exit(3)
	}

	templ, err := template.New("Error").Funcs(funcMap).Parse(string(content[:]))
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to parse the template file '%s' because '%s'\n", infile, err.Error())
		os.Exit(2)
	}

	err = templ.Execute(os.Stdout, struct{}{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Template execution failed '%s'\n", err.Error())
	}
}
