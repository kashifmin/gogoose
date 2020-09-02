package main

import (
	"flag"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"text/template"

	"github.com/joncalhoun/pipe"
)

type data struct {
	StructType   string
	StructImport string
	Name         string
}

func main() {
	var d data
	flag.StringVar(&d.StructType, "type", "gogoose.User", "The struct type for model being generated")
	flag.StringVar(&d.StructImport, "import", "github.com/kashifmin/gogoose", "The struct import path for model being generated")
	flag.StringVar(&d.Name, "name", "User", "The prefix name for model structs")
	flag.Parse()

	templateBytes, err := ioutil.ReadFile("templates/main.txt")
	if err != nil {
		panic(err)
	}

	t := template.Must(template.New("gogoose").Parse(string(templateBytes)))
	rc, wc, _ := pipe.Commands(
		exec.Command("gofmt"),
		exec.Command("goimports"),
	)
	t.Execute(wc, d)
	wc.Close()
	io.Copy(os.Stdout, rc)
}
