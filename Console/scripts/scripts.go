package scripts

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"

	nodes "github.com/TerminalJockey/Kudzu/Nodes"
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
			if file.IsDir() == true {
				curdir, err := os.Open("Scripts/" + file.Name())
				if err != nil {
					log.Println(err)
				}
				curfiles, err := curdir.ReadDir(0)
				for _, f := range curfiles {
					fmt.Println(file.Name() + "/" + f.Name())
				}
			} else {
				fmt.Println(file.Name())
			}

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

//ScriptCompileAndRun temporarily compiles script to tmp/UID.go, uses go run to execute, then clears the go source from tmp dir
func ScriptCompileAndRun(scropts interface{}, input ...string) {
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
		tmpdir := "tmp/" + GenUID() + ".go"
		tmpfile, err := os.Create(tmpdir)

		if err != nil {
			log.Println(err)
		}
		_, err = tmpfile.Write(tmplbuf.Bytes())
		if err != nil {
			log.Println(err)
		}

		runcmd := exec.Command("go", "run", tmpdir)
		runcmd.Env = append(os.Environ(), "GOOS=windows", "GOARCH=amd64", "GO111MODULE=off")
		runcmd.Stdout = os.Stdout
		runcmd.Stderr = os.Stderr
		err = runcmd.Start()
		if err != nil {
			log.Println(err)
		}
		runcmd.Wait()
		tmpfile.Close()

		err = os.Remove(tmpdir)
		if err != nil {
			log.Println(err)
		}
		fmt.Println("script completed")

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
			GoPath: os.Getenv("GOPATH"),
		})
		i.Use(stdlib.Symbols)
		i.Use(unrestricted.Symbols)
		i.Use(syscall.Symbols)
		i.Use(unsafe.Symbols)

		_, err = i.Eval(string(tmplbuf.Bytes()))
		if err != nil {
			log.Println(err)
		}
	}
}

//ScriptSend populates a script with a given struct and sends it to a node for processing
func ScriptSend(scriptname, NodeID string, scropts interface{}) {
	scriptbytes, err := os.ReadFile("Scripts/" + scriptname)
	if err != nil {
		log.Println(err)
		return
	}
	tmplbuf := new(bytes.Buffer)
	tmpl, err := template.New("").Parse(string(scriptbytes))
	if err != nil {
		log.Println(err)
	}
	tmpl.Execute(tmplbuf, scropts)

	encodedscript := base64.RawStdEncoding.EncodeToString(tmplbuf.Bytes())
	for i := range nodes.Nodes {
		if nodes.Nodes[i].NodeOpts.ID == NodeID {
			fmt.Printf("running %s on Node: %s with options %+v\n", scriptname, NodeID, scropts)
			_, err = nodes.Nodes[i].Conn.Write([]byte("runscript\n"))
			if err != nil {
				log.Println(err)
			}
			_, err = nodes.Nodes[i].Conn.Write([]byte(encodedscript + "\n"))
			if err != nil {
				log.Println(err)
			}
		}
	}
}

//ScriptSendandReturn populates a script with a given struct and sends it to a node for processing
func ScriptSendandReturn(scriptname, NodeID string, scropts interface{}) {
	scriptbytes, err := os.ReadFile("Scripts/" + scriptname)
	if err != nil {
		log.Println(err)
		return
	}
	tmplbuf := new(bytes.Buffer)
	tmpl, err := template.New("").Parse(string(scriptbytes))
	if err != nil {
		log.Println(err)
	}
	tmpl.Execute(tmplbuf, scropts)

	encodedscript := base64.RawStdEncoding.EncodeToString(tmplbuf.Bytes())
	for i := range nodes.Nodes {
		if nodes.Nodes[i].NodeOpts.ID == NodeID {
			fmt.Printf("running %s on Node: %s with options %+v\n", scriptname, NodeID, scropts)
			_, err = nodes.Nodes[i].Conn.Write([]byte("runwithoutput\n"))
			if err != nil {
				log.Println(err)
			}
			_, err = nodes.Nodes[i].Conn.Write([]byte(encodedscript + "\n"))
			if err != nil {
				log.Println(err)
			}
		}
	}
}

//ScriptOps will hold all the options for our scripts
type ScriptOps struct {
	LHOST, RHOST, CMD, LPORT, RPORT string
}

//GetJsonStruct extracts options from kdz script, and matches script type with appropriate struct. This was tricky without generics
func GetJsonStruct(input string) (winlocal WinLocal, winremote WinRemote, web Web, linlocal LinuxLocal, linremote LinuxRemote, compileandrun CompileAndRun) {
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
				fmt.Println("Script Type:", scrtype)
			}
			if strings.HasPrefix(strings.TrimSpace(i), "Description") == true {
				fmt.Println(i)
			}
		}
	}
	begin := bytes.Index(scriptbytes[intro+3:outtro], []byte("{"))
	end := bytes.Index(scriptbytes[intro+3:outtro], []byte("}"))
	fmt.Println("Script Options:", string(scriptbytes[begin+5:end+1]))
	switch scrtype {
	case "WinLocal":
		test := WinLocal{}
		err := json.Unmarshal(scriptbytes[begin:end+1], &test)
		if err != nil {
			if err.Error() != "invalid character ':' looking for beginning of value" {
				log.Println(err)
			}
		}
		test.Use = true

		return test, WinRemote{Use: false}, Web{Use: false}, LinuxLocal{Use: false}, LinuxRemote{Use: false}, CompileAndRun{Use: false}
	case "WinRemote":
		test := WinRemote{}
		err := json.Unmarshal(scriptbytes[begin:end+1], &test)
		if err != nil {
			if err.Error() != "invalid character ':' looking for beginning of value" {
				log.Println(err)
			}
		}
		test.Use = true

		return WinLocal{Use: false}, test, Web{Use: false}, LinuxLocal{Use: false}, LinuxRemote{Use: false}, CompileAndRun{Use: false}
	case "Web":
		test := Web{}
		err := json.Unmarshal(scriptbytes[begin:end+1], &test)
		if err != nil {
			if err.Error() != "invalid character ':' looking for beginning of value" {
				log.Println(err)
			}
		}
		test.Use = true

		return WinLocal{Use: false}, WinRemote{Use: false}, test, LinuxLocal{Use: false}, LinuxRemote{Use: false}, CompileAndRun{Use: false}
	case "LinuxLocal":
		test := LinuxLocal{}
		err := json.Unmarshal(scriptbytes[begin:end+1], &test)
		if err != nil {
			if err.Error() != "invalid character ':' looking for beginning of value" {
				log.Println(err)
			}
		}
		test.Use = true

		return WinLocal{Use: false}, WinRemote{Use: false}, Web{Use: false}, test, LinuxRemote{Use: false}, CompileAndRun{Use: false}
	case "LinuxRemote":
		test := LinuxRemote{}
		err := json.Unmarshal(scriptbytes[begin:end+1], &test)
		if err != nil {
			if err.Error() != "invalid character ':' looking for beginning of value" {
				log.Println(err)
			}
		}
		test.Use = true

		return WinLocal{Use: false}, WinRemote{Use: false}, Web{Use: false}, LinuxLocal{Use: false}, test, CompileAndRun{Use: false}
	case "CompileAndRun":
		test := CompileAndRun{}
		err := json.Unmarshal(scriptbytes[begin:end+1], &test)
		if err != nil {
			if err.Error() != "invalid character ':' looking for beginning of value" {
				log.Println(err)
			}
		}
		test.Use = true

		return WinLocal{Use: false}, WinRemote{Use: false}, Web{Use: false}, LinuxLocal{Use: false}, LinuxRemote{Use: false}, test

	}
	return
}

//WinLocal holds options for local windows scripts
type WinLocal struct {
	Use       bool   `json:"Use"`
	NodeID    string `json:"NodeID"`
	Username  string `json:"Username"`
	Password  string `json:"Password"`
	Domain    string `json:"Domain"`
	Cmd       string `json:"Cmd"`
	Lhost     string `json:"Lhost"`
	Lport     string `json:"Lport"`
	Rhost     string `json:"Rhost"`
	Rport     string `json:"Rport"`
	Hostname  string `json:"Hostname"`
	Directory string `json:"Directory"`
	Filename  string `json:"Filename"`
}

//WinRemote holds options for remote windows scripts
type WinRemote struct {
	Use       bool   `json:"Use"`
	NodeID    string `json:"NodeID"`
	Username  string `json:"Username"`
	Password  string `json:"Password"`
	Domain    string `json:"Domain"`
	Cmd       string `json:"Cmd"`
	Lhost     string `json:"Lhost"`
	Lport     string `json:"Lport"`
	Rhost     string `json:"Rhost"`
	Rport     string `json:"Rport"`
	Hostname  string `json:"Hostname"`
	Directory string `json:"Directory"`
	Filename  string `json:"Filename"`
}

//Web holds options for web scripts
type Web struct {
	Use       bool   `json:"Use"`
	NodeID    string `json:"NodeID"`
	Username  string `json:"Username"`
	Password  string `json:"Password"`
	URL       string `json:"Domain"`
	Cmd       string `json:"Cmd"`
	Lhost     string `json:"Lhost"`
	Lport     string `json:"Lport"`
	Rhost     string `json:"Rhost"`
	Rport     string `json:"Rport"`
	Hostname  string `json:"Hostname"`
	Directory string `json:"Directory"`
}

//LinuxLocal holds options for local linux scripts
type LinuxLocal struct {
	Use       bool   `json:"Use"`
	NodeID    string `json:"NodeID"`
	Username  string `json:"Username"`
	Password  string `json:"Password"`
	Cmd       string `json:"Cmd"`
	Lhost     string `json:"Lhost"`
	Lport     string `json:"Lport"`
	Rhost     string `json:"Rhost"`
	Rport     string `json:"Rport"`
	Hostname  string `json:"Hostname"`
	Directory string `json:"Directory"`
	Filename  string `json:"Filename"`
}

//LinuxRemote holds options for remote linux scripts
type LinuxRemote struct {
	Use       bool   `json:"Use"`
	NodeID    string `json:"NodeID"`
	Username  string `json:"Username"`
	Password  string `json:"Password"`
	Cmd       string `json:"Cmd"`
	Lhost     string `json:"Lhost"`
	Lport     string `json:"Lport"`
	Rhost     string `json:"Rhost"`
	Rport     string `json:"Rport"`
	Hostname  string `json:"Hostname"`
	Directory string `json:"Directory"`
	Filename  string `json:"Filename"`
}

//CompileAndRun holds options for scripts that dont play well with yaegi
type CompileAndRun struct {
	Use       bool   `json:"Use"`
	NodeID    string `json:"NodeID"`
	Username  string `json:"Username"`
	Password  string `json:"Password"`
	Domain    string `json:"Domain"`
	Cmd       string `json:"Cmd"`
	Lhost     string `json:"Lhost"`
	Lport     string `json:"Lport"`
	Rhost     string `json:"Rhost"`
	Rport     string `json:"Rport"`
	Hostname  string `json:"Hostname"`
	Directory string `json:"Directory"`
	Filename  string `json:"Filename"`
}

func GenUID() (uid string) {
	t := time.Now()
	str := t.String()
	h := sha1.New()
	h.Write([]byte(str))
	uid = fmt.Sprintf("%x", h.Sum(nil))
	return string([]byte(uid[:10]))
}
