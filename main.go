package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/codegangsta/cli"
)

var (
	app                    *cli.App
	workDir, ignoreFileDir string
)

func init() {
	app = cli.NewApp()
	setWokrkDirPaths()
}

func setSubCommands() {
	app.Commands = []cli.Command{
		{
			Name:    "update",
			Aliases: []string{"u"},
			Usage:   "update gitignore files",
			Action:  update,
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "list up available files",
			Action:  list,
		},
		{
			Name:    "get",
			Aliases: []string{"g"},
			Usage:   "get gitignore file into current directory",
			Action:  get,
		},
	}
}

// update updates gitignore files.
func update(c *cli.Context) {
	if err := os.RemoveAll(ignoreFileDir); err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}

	cloneIgnoreFilesIfNotExist()
}

func setWokrkDirPaths() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	workDir = filepath.Join(usr.HomeDir, ".getignore")
	ignoreFileDir = filepath.Join(workDir, "gitignore_files")
}

// exists checks a file exists or not.
func exists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

// cloneIgnoreFilesIfNotExist get ignore files from repository on GitHub.
// Since it just call Git clone command internally, Git is required.
func cloneIgnoreFilesIfNotExist() {
	if exists(ignoreFileDir) {
		return
	}

	fmt.Println("Cloning gitginore files. This may take a while...")

	if err := os.Mkdir(workDir, 0777); err != nil {
		if !os.IsExist(err) {
			log.Fatal(err)
		}
	}

	cmd := exec.Command("git", "clone", "git@github.com:github/gitignore.git", ignoreFileDir)
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Done.")
}

func list(c *cli.Context) {
	cloneIgnoreFilesIfNotExist()

	fileInfos, err := ioutil.ReadDir(ignoreFileDir)
	if err != nil {
		log.Fatal(err)
	}

	pattern := "\\.gitignore$"
	query := c.Args().First()
	if len(query) > 0 {
		pattern = fmt.Sprintf(".*%s.*\\.gitignore$", strings.ToLower(query))
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		log.Fatal(err)
	}

	for _, fi := range fileInfos {
		var name = fi.Name()
		matched := re.MatchString(strings.ToLower(name))
		if err != nil {
			log.Fatal(err)
		}
		if matched {
			fmt.Printf("%s\n", name)
		}
	}
}

func get(c *cli.Context) {
	cloneIgnoreFilesIfNotExist()
	srcName := c.Args().First() + ".gitignore"
	srcPath := filepath.Join(ignoreFileDir, srcName)

	if !exists(srcPath) {
		fmt.Printf("%s is not found. Please enter correct name. (It is case sensitive. e.g. Scala, CakePHP.)\n", srcName)
		return
	}

	src, err := os.Open(srcPath)
	if err != nil {
		log.Fatal(src)
	}
	defer src.Close()

	dstPath := ".gitignore"
	if exists(dstPath) {
		fmt.Println("Warn! .gitignore already exists in this directory. Remove it before execute.")
		return
	}

	dst, err := os.Create(dstPath)
	if err != nil {
		log.Fatal(src)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(src)
	}

	fmt.Printf("Done: %s\n", srcName)
}

func main() {
	app.Name = "getignore"
	app.Usage = "get .gitignore file from https://github.com/github/gitignore"
	setSubCommands()
	app.Run(os.Args)
}
