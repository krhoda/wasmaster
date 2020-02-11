package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/krhoda/wasmaster/asset"
	"github.com/krhoda/wasmaster/tmpstr"
)

var (
	npmCmd   = "npm"
	npxCmd   = "npx"
	cargoCmd = "cargo"
	wasmCmd  = "wasm-bindgen"

	projectName = ""

	pathExecList = []string{
		"npm",
		"npx",
		"cargo",
		"wasm-bindgen",
	}

	fileAssetList = []string{
		".babelrc",
		"index.html",
		"index.js",
		"lib.rs",
		"webpack.config.js",
	}

	fileTemplateMap = map[string]string{
		"Cargo.toml":     tmpstr.CargoToml,
		"package.json":   tmpstr.PackageJson,
		"wasm.worker.js": tmpstr.WebWorker,
	}

	assetMap = map[string][]byte{}

	topLevelDirs = []string{
		"js",
		"build",
		"dist",
	}

	webDeps = []string{
		"react",
		"react-dom",
	}

	webDevDeps = []string{
		"@babel/core",
		"babel-loader",
		"@babel/preset-env",
		"@babel/preset-react",
		"webpack",
		"webpack-cli",
	}
)

func preflight() {
	for _, target := range pathExecList {
		_, err := exec.LookPath(target)
		maybeFailWith(fmt.Sprintf("Could not find required executable on PATH: %s", target), err)
	}

	for _, fileAsset := range fileAssetList {
		data, err := asset.Asset(fmt.Sprintf("data/%s", fileAsset))
		maybeFailWith(fmt.Sprintf("Could not find fileAsset file: %s", fileAsset), err)
		assetMap[fileAsset] = data
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
	// initArgs := []string{"init", "-y"}
	// nInit := makeAndPipeCmd(npmCmd, initArgs)

	// err := nInit.Run()
	// maybeFailWith("Failed to run NPM init", err)

	// depArgs := append([]string{"install", "--save"}, webDeps...)
	// nDeps := makeAndPipeCmd(npmCmd, depArgs)
	nArgs := []string{"i"}
	nDeps := makeAndPipeCmd(npmCmd, nArgs)

	err := nDeps.Run()
	maybeFailWith("Failed to save dependecies with NPM.", err)

	// devDepArgs := append([]string{"install", "--save-dev"}, webDevDeps...)
	// nDevDeps := makeAndPipeCmd(npmCmd, devDepArgs)

	// err = nDevDeps.Run()
	// maybeFailWith("Failed to save dev-dependecies with NPM.", err)
}

func makeFile(name string, contents []byte) {
	f, err := os.Create(name)
	maybeFailWith(fmt.Sprintf("Could not create %s", name), err)
	defer f.Close()

	_, err = f.Write(contents)
	maybeFailWith(fmt.Sprintf("Could not write contents to %s", name), err)

	f.Sync()
}

func makeTemplate(name string, temp *template.Template) {
	f, err := os.Create(name)
	maybeFailWith(fmt.Sprintf("Could not create %s", name), err)
	defer f.Close()

	pn := projectTemplate{
		ProjectName: projectName,
	}

	err = temp.Execute(f, pn)
	maybeFailWith(fmt.Sprintf("Could not write %s", name), err)
}

func boilerplate() {
	for filename, contents := range assetMap {
		switch filename {

		case ".babelrc", "webpack.config.js":
			log.Println(filename)
			makeFile(filename, contents)

		case "index.html":
			fname := filepath.Join("dist", filename)
			makeFile(fname, contents)

		case "index.js":
			fname := filepath.Join("js", filename)
			makeFile(fname, contents)

		case "lib.rs":
			fname := filepath.Join("src", filename)
			makeFile(fname, contents)
		}
	}

	for filename, temp := range fileTemplateMap {
		t, err := template.New(filename).Parse(temp)
		maybeFailWith("Failed to parse %s template", err)

		switch filename {
		case "Cargo.toml", "package.json":
			makeTemplate(filename, t)
		case "wasm.worker.js":
			fname := filepath.Join("dist", filename)
			makeTemplate(fname, t)
		}
	}
}

type projectTemplate struct {
	ProjectName string
}

// TODO: CHANGE TO PACKAGE JSON TASKS AND RUN:
func testBuild() {
	wasmArgs := []string{"build", "--target", "wasm32-unknown-unknown"}
	wasmBuildCmd := makeAndPipeCmd(cargoCmd, wasmArgs)

	err := wasmBuildCmd.Run()
	maybeFailWith("Could not run build for WASM with cargo", err)

	bindgenArgs := []string{
		// TODO: ADD AFTER TEMPLATE.
		fmt.Sprintf("target/wasm32-unknown-unknown/debug/%s.wasm", projectName),
		"--no-typescript",
		"--no-modules",
		"--out-dir",
		"dist",
	}
	wasmBindCmd := makeAndPipeCmd(wasmCmd, bindgenArgs)

	err = wasmBindCmd.Run()
	maybeFailWith("Could not run build for WASM with wasm-bindgen", err)

	webpackCmd := makeAndPipeCmd(npxCmd, []string{"webpack"})

	err = webpackCmd.Run()
	maybeFailWith("Could not run webpack from NPX", err)
}

func createProject() {
	createProjectDirectories()
	log.Println("INITIAL DIRECTORY STRUCTURE COMPLETE. STARTING NPM...")

	boilerplate()
	log.Println("BOILERPLATE WRITTEN. STARTING BUILD OF TOTAL PROJECT...")

	npmania()
	log.Println("NPM INIT/DEPS COMPLETE. STARTING BOILERPLATE...")

	// TODO: REPLACE WITH PACKAGE JSON FIDDLING.
	testBuild()
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

	projectName = td
	createProject()
}
