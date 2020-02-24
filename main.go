package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/yakipote/gostnote/storage/firebase"
	"github.com/yakipote/gostnote/termbox"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"

	"github.com/urfave/cli/v2"
)

type Config struct {
	Firebase firebase.Config
}

type Storage interface {
	Upload(file *os.File)
	List() error
	GetFileList() ([]string, error)
	Download(fileName string) []byte
}

var configDir string
var storage Storage

//go:generate sh -c "go get -v github.com/jessevdk/go-assets-builder; go-assets-builder configExample.toml > configExample.go"
func init() {
	config := &Config{}
	// read config file
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	configDir = usr.HomeDir + "/gostnote/"
	// if not exists config dir then make dir
	if !exists(configDir) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			fmt.Println(err)
		}
		if err := os.MkdirAll(configDir+"tmp", 0755); err != nil {
			fmt.Println(err)
		}
		f, err := Assets.Open("/configExample.toml")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		dst, err := os.Create(configDir + "config.toml")
		if err != nil {
			panic(err)
		}
		defer dst.Close()

		_, err = io.Copy(dst, f)
		if err != nil {
			panic(err)
		}
	}

	_, err = toml.DecodeFile(configDir+"config.toml", &config)
	if err != nil {
		log.Fatalln(err)
	}
	//todo: defaultのstorageによって場合分け
	storage = firebase.NewStorage(config.Firebase)
}

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "new",
				Usage:   "new page",
				Aliases: []string{"n"},
			},
			&cli.StringFlag{
				Name:    "edit",
				Usage:   "edit page",
				Aliases: []string{"e"},
			},
			&cli.BoolFlag{
				Name:    "list",
				Usage:   "list pages",
				Aliases: []string{"l"},
			},
		},
		Action: appRun,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func appRun(c *cli.Context) error {

	// display file list
	listFlg := c.Bool("list")
	if listFlg {
		if err := storage.List(); err != nil {
			return err
		}
		return nil
	}

	// edit
	fileName := c.String("edit")
	data := []byte("")
	if fileName != "" {
		data = storage.Download(fileName)
	} else {
		fileName = c.Args().Get(0)
	}

	// no flag no args
	if c.NumFlags() == 0 && c.Args().Len() == 0 {
		list, err := storage.GetFileList()
		if err != nil {
			return err
		}
		t := termbox.NewTermbox(list)
		t.Init()
		t.Display()
		fileName = t.SelectionFile
		data = storage.Download(fileName)
	}

	// make tmp file
	fp := configDir + "tmp/" + fileName
	err := ioutil.WriteFile(fp, data, 0664)
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}
	defer func() {
		err = os.Remove(fp)
		if err != nil {
			fmt.Printf(err.Error())
		}
	}()

	cmd := exec.Command("vim", fp)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err = cmd.Run()
	if err != nil {
		fmt.Printf(err.Error())
		return err
	}

	content, _ := os.Open(fp)
	defer func() {
		err = content.Close()
		if err != nil {
			fmt.Printf(err.Error())
		}
	}()
	storage.Upload(content)
	return nil
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func copy(srcName string, dstName string) {
	src, err := os.Open(srcName)
	if err != nil {
		panic(err)
	}
	defer src.Close()

	dst, err := os.Create(dstName)
	if err != nil {
		panic(err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		panic(err)
	}
}
