package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"github.com/krhoda/wasmaster/asset"
)

var (
	npmCmd   = "npm"
	cargoCmd = "cargo"
	wasmCmd  = "wasm-bindgen"

	pathExecList = []string{
		"npm",
		"cargo",
		"wasm-bindgen",
	}

	fileTemplateList = []string{
		"Cargo.toml",
		"index.html",
		"index.js",
		"lib.rs",
		"webpack.config.js",
	}

	templateMap = map[string][]byte{}

	topLevelDirs = []string{
		"build",
		"dist",
	}

	webDeps = []string{
		"react",
		"react-dom",
	}

	webDevDeps = []string{
		"babel-core",
		"babel-loader",
		"babel-preset-env",
		"babel-preset-react",
		"babel-plugin-syntax-dynamic-import",
		"webpack",
		"webpack-cli",
	}
)

func preflight() {
	for _, target := range pathExecList {
		_, err := exec.LookPath(target)
		maybeFailWith(fmt.Sprintf("Could not find required executable on PATH: %s", target), err)
	}

	for _, template := range fileTemplateList {
		data, err := asset.Asset(fmt.Sprintf("data/%s", template))
		maybeFailWith(fmt.Sprintf("Could not find template file: %s", template), err)
		templateMap[template] = data
	}
}

func maybeFailWith(message string, err error) {
	if err == nil {
		return
	}

	log.Fatalf("%s\nError: %s\n", message, err.Error())
}

func createProjectDirectories() {
	for _, target := range topLevelDirs {
		err := os.Mkdir(target, 0755)
		msg := fmt.Sprintf("Could not create top-level directory: %s", target)
		maybeFailWith(msg, err)
	}
}

func makeAndPipeCmd(cmdString string, cmdArgs []string) *exec.Cmd {
	nextCommand := exec.Command(cmdString, cmdArgs...)
	nextCommand.Stderr = os.Stderr
	nextCommand.Stdout = os.Stdout

	return nextCommand
}

func firstCargo(td string) {
	firstCargoArgs := []string{"new", td, "--lib"}
	firstCargoCmd := makeAndPipeCmd(cargoCmd, firstCargoArgs)
	err := firstCargoCmd.Run()
	maybeFailWith("Failed to initialize project with cargo", err)
}

func npmania() {
	initArgs := []string{"init", "-y"}
	nInit := makeAndPipeCmd(npmCmd, initArgs)

	err := nInit.Run()
	maybeFailWith("Failed to run NPM init", err)

	depArgs := append([]string{"install", "--save"}, webDeps...)
	nDeps := makeAndPipeCmd(npmCmd, depArgs)

	err = nDeps.Run()
	maybeFailWith("Failed to save dependecies with NPM.", err)

	devDepArgs := append([]string{"install", "--save-dev"}, webDevDeps...)
	nDevDeps := makeAndPipeCmd(npmCmd, devDepArgs)

	err = nDevDeps.Run()
	maybeFailWith("Failed to save dev-dependecies with NPM.", err)
}

func boilerplate() {
	for filename, contents := range templateMap {
		switch filename {
		case ".babelrc", "Cargo.toml", "webpack.config.js":
			makeFile(filename, contents)
		case "index.js", "lib.rs":
			fname := filepath.Join("src", filename)
			makeFile(fname, contents)
		case "index.html":
			fname := filepath.Join("dist", filename)
			makeFile(fname, contents)
		}
	}
}

func makeFile(name string, contents []byte) {
	f, err := os.Create(name)
	maybeFailWith(fmt.Sprintf("Could not create %s", name), err)
	defer f.Close()

	_, err = f.Write(contents)
	maybeFailWith(fmt.Sprintf("Could not write contents to %s", name), err)

	f.Sync()
}

func createProject() {
	createProjectDirectories()
	log.Println("INITIAL DIRECTORY STRUCTURE COMPLETE. STARTING NPM...")
	npmania()
	log.Println("NPM INIT/DEPS COMPLETE. STARTING BOILERPLATE...")
	boilerplate()
}

func main() {
	log.Println("BEHOLD WASMASTER: LIKE CREATE REACT APP WITHOUT DOCUMENTATION...")
	log.Println("... BUT WITH RUST-BASED WASM!")

	preflight()

	argLen := len(os.Args)
	if argLen != 2 {
		maybeFailWith("Wasmaster accepts a single argument, a non-existing directory name. Any thing more or less results in this error", fmt.Errorf("Too many or few arguments"))
	}

	td := os.Args[1]
	_, err := os.Stat(td)
	if err == nil || !os.IsNotExist(err) {
		if err == nil {
			err = fmt.Errorf("Directory exists")
		}

		maybeFailWith("Wasmaster accepts a single argument, a non-existing directory name. Any thing more or less results in this error", err)
	}

	log.Println("PREFLIGHT CHECKS PASSED. CREATING PROJECT...")

	firstCargo(td)

	err = os.Chdir(td)
	maybeFailWith(fmt.Sprintf("Failed to change into newly created directory: %s", td), err)

	log.Println("RUST SKELETON PROJECT COMPLETE. STARTING NPM...")

	createProject()
}
