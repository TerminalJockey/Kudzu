package scripts

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"github.com/traefik/yaegi/stdlib/syscall"
	"github.com/traefik/yaegi/stdlib/unrestricted"
	"github.com/traefik/yaegi/stdlib/unsafe"
)

//ScriptList lists scripts, can be modified with a filter
func ScriptList(input ...string) {
	switch len(input) {
	case 1:
		fmt.Println("listing all scripts: ")
		scriptdir, err := os.Open("Scripts\\")
		if err != nil {
			log.Println(err)
		}
		files, err := scriptdir.Readdir(0)
		if err != nil {
			log.Println(err)
		}
		for _, file := range files {
			fmt.Println(file.Name())
		}
	case 2:
		if input[1] == "help" || input[1] == "-h" {
			fmt.Println("Lists scripts available, filter with -f")
		}
	case 3:
		if input[1] == "-f" {
			fmt.Println("filtering by:", input[2])
		}
	}
}

//ScriptRun runs scripts
func ScriptRun(scropts ScriptOps, input ...string) {
	fmt.Println("Output:")
	if strings.HasSuffix(input[1], ".kzs") == true {
		//get template bytes
		scriptbytes, err := ioutil.ReadFile("Scripts\\" + input[1])
		if err != nil {
			log.Println(err)
		}
		//setup template buffer
		tmplbuf := new(bytes.Buffer)
		tmpl, err := template.New("").Parse(string(scriptbytes))
		tmpl.Execute(tmplbuf, scropts)

		cmdbuf := new(bytes.Buffer)
		i := interp.New(interp.Options{
			Stdout: cmdbuf,
		})
		i.Use(stdlib.Symbols)
		i.Use(unrestricted.Symbols)
		i.Use(syscall.Symbols)
		i.Use(unsafe.Symbols)

		_, err = i.Eval(string(tmplbuf.Bytes()))
		if err != nil {
			log.Println(err)
		}
		fmt.Println(cmdbuf.String())
	}
}

//ScriptGetOpts gets scripts from .kzs head, and returns struct for filling
func ScriptGetOpts(input string) {
	scriptbytes, err := ioutil.ReadFile("Scripts\\" + input)
	if err != nil {
		log.Println(err)
	}
	intro := bytes.Index(scriptbytes, []byte("/*{"))
	outtro := bytes.Index(scriptbytes, []byte("}*/"))
	if intro == -1 || outtro == -1 {
		fmt.Println("Check the formatting of your kz script")
		return
	}
	out := strings.Split(string(scriptbytes[intro+3:outtro]), "\n")
	for _, i := range out {
		if strings.TrimSpace(i) != "" && len(strings.TrimSpace(i)) > 1 {
			fmt.Println(i)
		}
	}
}

//ScriptOps will hold all the options for our scripts
type ScriptOps struct {
	LHOST, RHOST, CMD, LPORT, RPORT string
}
