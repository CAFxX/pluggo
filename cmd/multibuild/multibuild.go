package main

import (
	"flag"
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/kisielk/gotool"
)

func usage() {
	fmt.Println("multibuild builds one or more packages and links them together")
	fmt.Println("multibuild <mainPkg> [<pkg>*]")
	fmt.Println("  <mainPkg>   main package to build")
	fmt.Println("  <pkg>       additional package to link with <mainPkg>")
	fmt.Println("              it is possible to specify more than one <pkg>")
	fmt.Println("              <pkg> can't be a package named \"main\"")
	fmt.Println("The resulting build will have the package name of <mainPkg>")
}

func multibuild(packages []string) (error, int) {
	// parse the names of the paths/packages provided on the command-line
	pkgs := gotool.ImportPaths(packages)

	// parse the main package: we need the package name and the path
	mainpkg, err := build.Import(pkgs[0], ".", 0)
	if err != nil {
		return fmt.Errorf("unable to import main package %s: %s", pkgs[0], err), -2
	}

	// create the linker file in the main package directory, so that it will be
	// compiled and linked with the main package when we `go build` is invoked
	// the linker file is a go file in the same package of the main package that
	// imports all additional packages (normally for their side-effects)
	tmpFile, err := ioutil.TempFile(mainpkg.Dir, "pluggo")
	if err != nil {
		return fmt.Errorf("unable to create temporary file: %s", err), -3
	}

	fmt.Fprintf(tmpFile, "package %s\n", mainpkg.Name)
	for _, pkgname := range pkgs[1:] {
		fmt.Fprintf(tmpFile, "import _ \"%s\"\n", pkgname)
	}

	tmpFile.Close()
	os.Rename(tmpFile.Name(), tmpFile.Name()+".go")
	defer os.Remove(tmpFile.Name() + ".go")

	// run go build on the main package
	output, err := exec.Command("go", "build", packages[0]).CombinedOutput()
	if err != nil {
		return fmt.Errorf("error executing go build: %s\ngo build output:\n%s", err, string(output)), -4
	}

	return nil, 0
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println("missing arguments")
		usage()
		os.Exit(-1)
	}

	err, rv := multibuild(flag.Args())
	if err != nil {
		fmt.Println(err)
		os.Exit(rv)
	}
}
