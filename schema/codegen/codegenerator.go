package main

import (
	"io/ioutil"
	"os"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/elodina/go-avro"
	"github.com/urfave/cli"
)

var (
	in = cli.StringFlag{
		Name:  "schema",
		Usage: "Input file name of the avro schema",
		Value: "",
	}
	out = cli.StringFlag{
		Name:  "out",
		Usage: "Output file name",
		Value: "",
	}
	log = logger.GetOrCreate("main")
)

func main() {
	app := cli.NewApp()
	app.Name = "Avro code generator"
	app.Usage = "Generate go code from avro schema"
	app.Flags = []cli.Flag{in, out}
	app.Action = func(c *cli.Context) error {
		return startProcess(c)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
		return
	}

}

func startProcess(c *cli.Context) error {
	input := c.GlobalString(in.Name)
	output := c.GlobalString(out.Name)

	schema, err := ioutil.ReadFile(input)
	if err != nil {
		return err
	}

	codeGenerator := avro.NewCodeGenerator([]string{string(schema)})
	code, err := codeGenerator.Generate()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(output, []byte(code), 0664)
}
