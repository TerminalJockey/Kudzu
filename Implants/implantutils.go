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
		//setup template buffer
		tmplbuf := new(bytes.Buffer)
		tmpl, err := template.New("").Parse(string(scriptbytes))
		type Filler struct {
			Laddr string
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
	case "psh":
		scriptbytes, err := ioutil.ReadFile("Scripts/Implants/simpleagent_psh.kzs")
		if err != nil {
			log.Println(err)
		}
		//setup template buffer
		tmplbuf := new(bytes.Buffer)
		tmpl, err := template.New("").Parse(string(scriptbytes))
		if err != nil {
			log.Println(err)
		}
		type Filler struct {
			Laddr string
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
	case "kdzshell_win":
		scriptbytes, err := ioutil.ReadFile("Scripts/Implants/kdzshell_win.kzs")
		if err != nil {
			log.Println(err)
		}
		//setup template buffer
		tmplbuf := new(bytes.Buffer)
		tmpl, err := template.New("").Parse(string(scriptbytes))
		type Filler struct {
			Laddr string
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
	case "sh":
		scriptbytes, err := ioutil.ReadFile("Scripts/Implants/simpleagent_sh.kzs")
		if err != nil {
			log.Println(err)
		}
		//setup template buffer
		tmplbuf := new(bytes.Buffer)
		tmpl, err := template.New("").Parse(string(scriptbytes))
		if err != nil {
			log.Println(err)
		}
		type Filler struct {
			Laddr string
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
	case "cmd_tls":
		scriptbytes, err := ioutil.ReadFile("Scripts/Implants/simpleagent_cmd_tls.kzs")
		if err != nil {
			log.Println(err)
		}
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
		type Filler struct {
			Laddr string
			Cert  string
			Key   string
		}

		filler := Filler{
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
	case "psh_tls":
		scriptbytes, err := ioutil.ReadFile("Scripts/Implants/simpleagent_psh_tls.kzs")
		if err != nil {
			log.Println(err)
		}
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
		type Filler struct {
			Laddr string
			Cert  string
			Key   string
		}

		filler := Filler{
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
	case "sh_tls":
		scriptbytes, err := ioutil.ReadFile("Scripts/Implants/simpleagent_sh_tls.kzs")
		if err != nil {
			log.Println(err)
		}
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
		type Filler struct {
			Laddr string
			Cert  string
			Key   string
		}

		filler := Filler{
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
	case "kdzshell_win_tls":
		scriptbytes, err := ioutil.ReadFile("Scripts/Implants/kdzshell_win_tls.kzs")
		if err != nil {
			log.Println(err)
		}
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
		type Filler struct {
			Laddr string
			Cert  string
			Key   string
		}

		filler := Filler{
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
	default:
		fmt.Println("implant types: cmd, cmd_tls, psh, psh_tls, kdzshell_win, kdzshell_win_tls, sh, sh_tls")
	}
}
