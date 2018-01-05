// +build ignore

package main

import (
	"goddbench"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"text/template"
	"time"
)

func TimeString(t time.Time) string {
	return t.Format(time.RFC3339Nano)
}

func main() {
	data := struct {
		Name string
		Out  string
		Fmt  string
		Seed int
		Len  int
	}{}

	flag.StringVar(&data.Name, "name", "", "Service name")
	flag.StringVar(&data.Out, "out", "", "Output directory")
	flag.StringVar(&data.Fmt, "fmt", "gofmt", "Formatter")
	flag.IntVar(&data.Seed, "seed", 0, "Random seed")
	flag.IntVar(&data.Len, "len", 1000, "Length of data list")
	flag.Parse()

	rand.Seed(int64(data.Seed))

	if data.Name == "" {
		data.Name = fmt.Sprintf("gob%d", data.Len)
	}

	if data.Out == "" {
		data.Out = fmt.Sprintf("app/%s", data.Name)
	}

	os.MkdirAll(data.Out, 0755)

	templ := template.New("app.yaml.txt")
	templ.Funcs(template.FuncMap{"GenN": goddbench.GenN, "TimeString": TimeString})
	_, err := templ.ParseFiles("app.yaml.txt", "app.go.txt", "gob_tmp.go.txt")
	if err != nil {
		log.Fatal(err)
	}

	{
		cmd := exec.Command(data.Fmt)
		file, err := cmd.StdinPipe()
		if err != nil {
			log.Fatal(err)
		}

		cmd.Stderr = os.Stderr
		path := fmt.Sprintf("%s/list.go", data.Out)
		cmd.Stdout, err = os.Create(path)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Writing to", path)

		err = cmd.Start()
		if err != nil {
			log.Fatal(err)
		}

		err = templ.ExecuteTemplate(file, "gob_tmp.go.txt", data)
		if err != nil {
			log.Fatal(err)
		}

		file.Close()

		err = cmd.Wait()
		if err != nil {
			log.Fatal(err)
		}
	}

	{
		path := fmt.Sprintf("%s/list.gob", data.Out)
		file, err := os.Create(path)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		log.Println("Writing to", path)
		e := gob.NewEncoder(file)
		e.Encode(goddbench.GenN(data.Len))
	}

	{
		path := fmt.Sprintf("%s/app.yaml", data.Out)
		file, err := os.Create(path)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		log.Println("Writing to", path)
		err = templ.ExecuteTemplate(file, "app.yaml.txt", data)
		if err != nil {
			log.Fatal(err)
		}

	}

	{
		path := fmt.Sprintf("%s/app.go", data.Out)
		file, err := os.Create(path)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		log.Println("Writing to", path)
		err = templ.ExecuteTemplate(file, "app.go.txt", data)
		if err != nil {
			log.Fatal(err)
		}
	}
}
