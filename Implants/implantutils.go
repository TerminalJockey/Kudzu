package implants

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"text/template"

	nodes "github.com/TerminalJockey/Kudzu/Nodes"
)

type TLS_Filler struct {
	Laddr string
	Cert  string
	Key   string
}

type Filler struct {
	Laddr string
}

//ImplantOps holds values for implant options
type ImplantOps struct {
	ImplantType, FileName, Arch string
	Listener                    nodes.Listener
}

//GenerateImplant takes an ImplantOps struct and generates a new implant
func GenerateImplant(Ops ImplantOps) {
	switch Ops.ImplantType {
	case "cmd":
		scriptbytes, err := ioutil.ReadFile("Scripts/Implants/simpleagent_cmd.kzs")
		if err != nil {
			log.Println(err)
		}
		genwindowsbasic(Ops, scriptbytes)

	case "psh":
		scriptbytes, err := ioutil.ReadFile("Scripts/Implants/simpleagent_psh.kzs")
		if err != nil {
			log.Println(err)
		}
		genwindowsbasic(Ops, scriptbytes)

	case "kdzshell_win":
		scriptbytes, err := ioutil.ReadFile("Scripts/Implants/kdzshell_win.kzs")
		if err != nil {
			log.Println(err)
		}
		genwindowsbasic(Ops, scriptbytes)

	case "sh":
		scriptbytes, err := ioutil.ReadFile("Scripts/Implants/simpleagent_sh.kzs")
		if err != nil {
			log.Println(err)
		}
		genlinuxbasic(Ops, scriptbytes)

	case "kdzshell_lin":
		scriptbytes, err := ioutil.ReadFile("Scripts/Implants/kdzshell_lin.kzs")
		if err != nil {
			log.Println(err)
		}
		genlinuxbasic(Ops, scriptbytes)

	case "cmd_tls":
		scriptbytes, err := ioutil.ReadFile("Scripts/Implants/simpleagent_cmd_tls.kzs")
		if err != nil {
			log.Println(err)
		}
		genwindowstls(Ops, scriptbytes)

	case "psh_tls":
		scriptbytes, err := ioutil.ReadFile("Scripts/Implants/simpleagent_psh_tls.kzs")
		if err != nil {
			log.Println(err)
		}
		genwindowstls(Ops, scriptbytes)

	case "sh_tls":
		scriptbytes, err := ioutil.ReadFile("Scripts/Implants/simpleagent_sh_tls.kzs")
		if err != nil {
			log.Println(err)
		}
		genlinuxtls(Ops, scriptbytes)

	case "kdzshell_win_tls":
		scriptbytes, err := ioutil.ReadFile("Scripts/Implants/kdzshell_win_tls.kzs")
		if err != nil {
			log.Println(err)
		}
		genwindowstls(Ops, scriptbytes)

	default:
		fmt.Println("implant types: cmd, cmd_tls, psh, psh_tls, kdzshell_win, kdzshell_win_tls, kdzshell_lin, sh, sh_tls")
	}
}

func genwindowsbasic(Ops ImplantOps, scriptbytes []byte) {
	tmplbuf := new(bytes.Buffer)
	tmpl, err := template.New("").Parse(string(scriptbytes))

	filler := Filler{
		Laddr: Ops.Listener.Listener.Addr().String(),
	}
	tmpl.Execute(tmplbuf, filler)
	tmpname := nodes.GenUID() + ".go"
	tmpfile, err := os.Create("tmp/" + tmpname)
	_, err = tmpfile.Write(tmplbuf.Bytes())
	if err != nil {
		log.Println(err)
	}
	buildcmd := exec.Command("go", "build", "-o", "tmp/"+Ops.FileName, "tmp/"+tmpname)
	buildcmd.Env = append(os.Environ(), "GOOS=windows", "GOARCH=amd64")
	err = buildcmd.Run()
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("generated implant! check tmp/%s\n", Ops.FileName)
	tmpfile.Close()
	err = os.Remove("tmp/" + tmpname)
	if err != nil {
		log.Println(err)
	}
}

func genlinuxbasic(Ops ImplantOps, scriptbytes []byte) {
	tmplbuf := new(bytes.Buffer)
	tmpl, err := template.New("").Parse(string(scriptbytes))
	if err != nil {
		log.Println(err)
	}

	filler := Filler{
		Laddr: Ops.Listener.Listener.Addr().String(),
	}
	tmpl.Execute(tmplbuf, filler)
	tmpname := nodes.GenUID() + ".go"
	tmpfile, err := os.Create("tmp/" + tmpname)
	_, err = tmpfile.Write(tmplbuf.Bytes())
	if err != nil {
		log.Println(err)
	}
	buildcmd := exec.Command("go", "build", "-o", "tmp/"+Ops.FileName, "tmp/"+tmpname)
	buildcmd.Env = append(os.Environ(), "GOOS=linux", "GOARCH=amd64")
	err = buildcmd.Run()
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("generated implant! check tmp/%s\n", Ops.FileName)
	tmpfile.Close()
	err = os.Remove("tmp/" + tmpname)
	if err != nil {
		log.Println(err)
	}
}

func genwindowstls(Ops ImplantOps, scriptbytes []byte) {
	//get tls cert
	cbytes, err := os.ReadFile("certs/" + Ops.Listener.ID + ".cert")
	if err != nil {
		log.Println(err)
		return
	}
	kbytes, err := os.ReadFile("certs/" + Ops.Listener.ID + ".key")
	if err != nil {
		log.Println(err)
		return
	}
	//setup template buffer
	tmplbuf := new(bytes.Buffer)
	tmpl, err := template.New("").Parse(string(scriptbytes))

	filler := TLS_Filler{
		Laddr: Ops.Listener.Listener.Addr().String(),
		Cert:  string(cbytes),
		Key:   string(kbytes),
	}
	tmpl.Execute(tmplbuf, filler)
	tmpname := nodes.GenUID() + ".go"
	tmpfile, err := os.Create("tmp/" + tmpname)
	_, err = tmpfile.Write(tmplbuf.Bytes())
	if err != nil {
		log.Println(err)
	}
	buildcmd := exec.Command("go", "build", "-o", "tmp/"+Ops.FileName, "tmp/"+tmpname)
	buildcmd.Env = append(os.Environ(), "GOOS=windows", "GOARCH=amd64", "GO111MODULE=off")
	err = buildcmd.Run()
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("generated implant! check tmp/%s\n", Ops.FileName)
	tmpfile.Close()
	err = os.Remove("tmp/" + tmpname)
	if err != nil {
		log.Println(err)
	}
}

func genlinuxtls(Ops ImplantOps, scriptbytes []byte) {
	//get tls cert
	cbytes, err := os.ReadFile("certs/" + Ops.Listener.ID + ".cert")
	if err != nil {
		log.Println(err)
		return
	}
	kbytes, err := os.ReadFile("certs/" + Ops.Listener.ID + ".key")
	if err != nil {
		log.Println(err)
		return
	}
	//setup template buffer
	tmplbuf := new(bytes.Buffer)
	tmpl, err := template.New("").Parse(string(scriptbytes))

	filler := TLS_Filler{
		Laddr: Ops.Listener.Listener.Addr().String(),
		Cert:  string(cbytes),
		Key:   string(kbytes),
	}
	tmpl.Execute(tmplbuf, filler)
	tmpname := nodes.GenUID() + ".go"
	tmpfile, err := os.Create("tmp/" + tmpname)
	_, err = tmpfile.Write(tmplbuf.Bytes())
	if err != nil {
		log.Println(err)
	}
	buildcmd := exec.Command("go", "build", "-o", "tmp/"+Ops.FileName, "tmp/"+tmpname)
	buildcmd.Env = append(os.Environ(), "GOOS=linux", "GOARCH=amd64", "GO111MODULE=off")
	err = buildcmd.Run()
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("generated implant! check tmp/%s\n", Ops.FileName)
	tmpfile.Close()
	err = os.Remove("tmp/" + tmpname)
	if err != nil {
		log.Println(err)
	}
}
