package scripts

import (
	"bytes"
	"encoding/json"
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
		scriptdir, err := os.Open("Scripts/")
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
func ScriptRun(scropts interface{}, input ...string) {
	fmt.Println("Output:")
	if strings.HasSuffix(input[1], ".kzs") == true {
		//get template bytes
		scriptbytes, err := ioutil.ReadFile("Scripts/" + input[1])
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
	scriptbytes, err := ioutil.ReadFile("Scripts/" + input)
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
			if strings.HasPrefix(strings.TrimSpace(i), "Type:") == true {
				scrtypearr := strings.Split(strings.TrimSpace(i), ":")
				scrtype := scrtypearr[1]
				fmt.Println(scrtype)
				getjson(scriptbytes[intro+3:outtro], scrtype)
			}
		}
	}
}

//ScriptOps will hold all the options for our scripts
type ScriptOps struct {
	LHOST, RHOST, CMD, LPORT, RPORT string
}

func GetJsonStruct(input string) (winlocal WinLocal, winremote WinRemote) {
	scriptbytes, err := ioutil.ReadFile("Scripts/" + input)
	if err != nil {
		log.Println(err)
	}
	intro := bytes.Index(scriptbytes, []byte("/*{"))
	outtro := bytes.Index(scriptbytes, []byte("}*/"))
	if intro == -1 || outtro == -1 {
		fmt.Println("Check the formatting of your kz script")
		return
	}
	var scrtype string
	out := strings.Split(string(scriptbytes[intro+3:outtro]), "\n")
	for _, i := range out {
		if strings.TrimSpace(i) != "" && len(strings.TrimSpace(i)) > 1 {
			if strings.HasPrefix(strings.TrimSpace(i), "Type") == true {

				scrtypearr := strings.Split(strings.TrimSpace(i), ":")
				scrtype = scrtypearr[1]
				fmt.Println(scrtype)
				break
			}
		}
	}
	begin := bytes.Index(scriptbytes[intro+3:outtro], []byte("{"))
	end := bytes.Index(scriptbytes[intro+3:outtro], []byte("}"))
	fmt.Println("getjson", string(scriptbytes[begin:end+1]))
	switch scrtype {
	case "WinLocal":
		test := WinLocal{}
		err := json.Unmarshal(scriptbytes[begin:end+1], &test)
		if err != nil {
			log.Println("unmarshal_winlocal:", err)
		}
		test.Use = true
		fmt.Println(test)
		return test, WinRemote{Use: false}
	case "WinRemote":
		test := WinRemote{}
		err := json.Unmarshal(scriptbytes[begin:end+1], &test)
		if err != nil {
			log.Println("unmarshal_winremote:", err)
		}
		test.Use = true
		fmt.Println(test)
		return WinLocal{Use: false}, test
	}
	return
}

func getjson(in []byte, scrtype string) (winlocal WinLocal, winremote WinRemote) {
	begin := bytes.Index(in, []byte("{"))
	end := bytes.Index(in, []byte("}"))
	fmt.Println("getjson", string(in[begin:end+1]))
	switch scrtype {
	case "WinLocal":
		test := WinLocal{}
		err := json.Unmarshal(in[begin:end+1], &test)
		if err != nil {
			log.Println(err)
		}
		test.Use = true
		fmt.Println(test)
		return test, WinRemote{Use: false}
	case "WinRemote":
		test := WinRemote{}
		err := json.Unmarshal(in[begin:end+1], &test)
		if err != nil {
			log.Println(err)
		}
		test.Use = true
		fmt.Println(test)
		return WinLocal{Use: false}, test
	}
	return
}

//WinLocal holds options for local windows scripts
type WinLocal struct {
	Use      bool   `json:"Use"`
	NodeID   string `json:"NodeID"`
	Username string `json:"Username"`
	Password string `json:"Password"`
	Domain   string `json:"Domain"`
	Cmd      string `json:"Cmd"`
	Lhost    string `json:"Lhost"`
	Lport    string `json:"Lport"`
	Rhost    string `json:"Rhost"`
	Rport    string `json:"Rport"`
	Hostname string `json:"Hostname"`
}

//WinRemote holds options for remote windows scripts
type WinRemote struct {
	Use      bool   `json:"Use"`
	NodeID   string `json:"NodeID"`
	Username string `json:"Username"`
	Password string `json:"Password"`
	Domain   string `json:"Domain"`
	Cmd      string `json:"Cmd"`
	Lhost    string `json:"Lhost"`
	Lport    string `json:"Lport"`
	Rhost    string `json:"Rhost"`
	Rport    string `json:"Rport"`
	Hostname string `json:"Hostname"`
}
