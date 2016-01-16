package main

import (
	"flag"
	"fmt"
	"go/build"
	"io/ioutil"
	"log"
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

func main() {
	flag.Usage = usage
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println("missing arguments")
		usage()
		os.Exit(-1)
	}

	pkgs := gotool.ImportPaths(flag.Args())

	mainpkg, err := build.Import(pkgs[0], ".", 0)
	if err != nil {
		log.Printf("unable to import main package %s: %s", pkgs[0], err)
		os.Exit(-2)
	}

	tmpFile, err := ioutil.TempFile(mainpkg.Dir, "pluggo")
	if err != nil {
		log.Printf("unable to create temporary file: %s", err)
		os.Exit(-3)
	}

	fmt.Fprintf(tmpFile, "package %s\n", mainpkg.Name)
	for _, pkgname := range os.Args[2:] {
		fmt.Fprintf(tmpFile, "import _ \"%s\"\n", pkgname)
	}

	tmpFile.Close()
	os.Rename(tmpFile.Name(), tmpFile.Name()+".go")
	defer os.Remove(tmpFile.Name() + ".go")

	output, err := exec.Command("go", "build", flag.Args()[0]).CombinedOutput()
	if err != nil {
		fmt.Print(string(output))
		os.Exit(-4)
	}
}
