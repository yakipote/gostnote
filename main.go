package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/urfave/cli/v2"
)

const NoteBook string = "/home/bun/Documents/"
func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "new",
				Usage:   "new file",
				Aliases: []string{"n"},
			},
			&cli.StringFlag{
				Name:  "edit",
				Usage: "edit file",
				Aliases: []string{"e"},
			},
		},
		//Action: func(c *cli.Context) error {
		//	name := "someone"
		//	if c.NArg() > 0 {
		//		name = c.Args().Get(0)
		//	}
		//	if language == "spanish" {
		//		fmt.Println("Hola", name)
		//	} else {
		//		fmt.Println("Hello", name)
		//	}
		//	return nil
		//},
		Action: appRun,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func appRun(c *cli.Context) error {
	fileName := c.Args().Get(0)
	fmt.Printf("go-note\n")
	cmd := exec.Command("/usr/bin/vim", NoteBook + fileName)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}
	fmt.Printf("go-note\n")
	return nil
}
