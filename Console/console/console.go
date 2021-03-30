package console

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/TerminalJockey/Kudzu/Console/scripts"
	implants "github.com/TerminalJockey/Kudzu/Implants"
	nodes "github.com/TerminalJockey/Kudzu/Nodes"
)

var cltag = "<kudzu> "
var curscript = ""

var winlocalopts scripts.WinLocal
var winremoteopts scripts.WinRemote
var linlocalopts scripts.LinuxLocal
var linremoteopts scripts.LinuxRemote
var webopts scripts.Web
var compileandrun scripts.CompileAndRun
var nodeopts nodes.NodeOpts = nodes.NodeOpts{
	NodeType: "tcp",
	Port:     "7896",
	Addr:     "127.0.0.1",
}
var implantops implants.ImplantOps = implants.ImplantOps{
	ImplantType: "cmd",
}

//InitConsole starts the kudzu console
func InitConsole() {
	in := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf(cltag)
		clinput, err := in.ReadString('\n')
		if err != nil {
			log.Println(err)
		}
		ParseCLI(clinput)
	}
}

//ParseCLI parses command line input and manages the Kudzu console
func ParseCLI(input string) {
	sep := strings.Split(input, " ")
	for i := range sep {
		sep[i] = strings.TrimSpace(sep[i])
	}
	switch sep[0] {
	//display help for each menu
	case "help", "?":
		switch cltag {
		case "<kudzu> ":
			fmt.Println("kudzu? zu kud!")
			fmt.Println("scripts: enter script menu.")
			fmt.Println("implants: enter implants menu. Used to generate new implants")
			fmt.Println("nodes: enter nodes menu. listener and node management")
		case "<kudzu scripts> ":
			fmt.Println("run: executes kudzu script")
			fmt.Println("	usage: run hello.kzs")
			fmt.Println("load: loads kudzu script for use and prints options")
			fmt.Println("	usage: load hello.kzs")
			fmt.Println("ls/list: lists scripts available")
			fmt.Println("	usage: list")
			fmt.Println("runscript: sends script to provided node for execution")
			fmt.Println("	usage: runscript Windows/win_calc.kzs")
			fmt.Println("runwithoutput: sends script to provided node, executes, and enters node shell")
			fmt.Println("	usage: runwithoutput Windows/winenum.kzs")
			fmt.Println("	^ only available on kdzshells ^")
		case "<kudzu implants> ":
			fmt.Println("setop: sets option for implant")
			fmt.Println("	options: implanttype (cmd/psh/kdzshell_win/sh/cmd_tls/psh_tls/kdzshell_win_tls/sh_tls)")
			fmt.Println("		    listener (ID)")
			fmt.Println("		    filename (name of generated implant)")
			fmt.Println("	usage: setop <option> <val>")
			fmt.Println("run: generates implant with given options")
			fmt.Println("	usage: run")
			fmt.Println("showops: prints currently set options")
			fmt.Println("	usage: showops")
		case "<kudzu nodes> ":
			fmt.Println("ls/list: lists listeners, nodes or both")
			fmt.Println("	usage: ls nodes/listeners")
			fmt.Println("run: starts listener with given options")
			fmt.Println("	usage: run")
			fmt.Println("setop: set options for listener")
			fmt.Println("	options: nodetype (tcp/tcp_tls/udp/ntp)")
			fmt.Println("		    address/addr (ip of listener)")
			fmt.Println("			port (port of listener)")
			fmt.Println("	usage: setop <option> <val>")
			fmt.Println("close: closes node or listener")
			fmt.Println("	usage: close node/listener <ID>")
			fmt.Println("interact: interact with node")
			fmt.Println("	usage: interact <ID>")
			fmt.Println("showops: shows options for listener")
			fmt.Println("	usage: showops")
		}
	//list assets in category
	case "list", "ls":
		switch cltag {
		case "<kudzu scripts> ":
			scripts.ScriptList(sep...)
		case "<kudzu nodes> ":
			if len(sep) > 2 {
				fmt.Println("usage: ls/list")
				fmt.Println("		ls listeners/nodes")
				return
			}
			if len(sep) == 2 && (sep[1] == "nodes" || sep[1] == "listeners") {
				switch sep[1] {
				case "nodes":
					fmt.Println("Nodes:")
					for _, n := range nodes.Nodes {
						fmt.Println("ID:", n.NodeOpts.ID)
						fmt.Println("NodeType:", n.NodeOpts.NodeType)
						fmt.Println("Addr:", n.NodeOpts.Addr+"\n")
					}
					return
				case "listeners":
					fmt.Println("Listeners:")
					for _, n := range nodes.Listeners {
						fmt.Println("ID:", n.ID)
						fmt.Println("ListenerType:", n.ListenerType)
						fmt.Println("Addr:", n.Listener.Addr().String()+"\n")
					}
					return
				}
			}
			fmt.Println("Nodes:")
			for _, n := range nodes.Nodes {
				fmt.Println("ID:", n.NodeOpts.ID)
				fmt.Println("NodeType:", n.NodeOpts.NodeType)
				fmt.Println("Addr:", n.NodeOpts.Addr+"\n")
			}
			fmt.Println("Listeners:")
			for _, m := range nodes.Listeners {
				fmt.Println("ID:", m.ID)
				fmt.Println("ListenerType:", m.ListenerType)
				fmt.Println("Addr:", m.Listener.Addr().String()+"\n")
			}
		default:
			fmt.Println("use from agents, nodes, or scripts menu")
		}
	case "runremote":
		switch cltag {
		case "<kudzu scripts> ":
			if len(sep) != 2 {
				fmt.Println("load a script using load <scriptname>, set your options, then fire away!")
				fmt.Println("usage: runremote <scriptname>")
				return
			}
			if winlocalopts.Use == true {
				scripts.ScriptSend(sep[1], winlocalopts.NodeID, winlocalopts)
			}
			if winremoteopts.Use == true {
				scripts.ScriptSend(sep[1], winremoteopts.NodeID, winremoteopts)
			}
			if webopts.Use == true {
				scripts.ScriptSend(sep[1], webopts.NodeID, webopts)
			}
			if linlocalopts.Use == true {
				scripts.ScriptSend(sep[1], linlocalopts.NodeID, linlocalopts)
			}
			if linremoteopts.Use == true {
				scripts.ScriptSend(sep[1], linremoteopts.NodeID, linremoteopts)
			}
		default:
			fmt.Println("use from scripts menu with a kudzu implant")
		}
	case "runwithoutput":
		switch cltag {
		case "<kudzu scripts> ":
			if len(sep) != 2 {
				fmt.Println("load a script using load <scriptname>, set your options, then fire away!")
				fmt.Println("usage: runremote <scriptname>")
				return
			}
			if winlocalopts.Use == true {
				scripts.ScriptSendandReturn(sep[1], winlocalopts.NodeID, winlocalopts)
				fmt.Println("interacting...")
				nodes.SelectNode(winlocalopts.NodeID)
				return
			}
			if winremoteopts.Use == true {
				scripts.ScriptSendandReturn(sep[1], winremoteopts.NodeID, winremoteopts)
				fmt.Println("interacting...")
				nodes.SelectNode(winremoteopts.NodeID)
				return
			}
			if webopts.Use == true {
				scripts.ScriptSendandReturn(sep[1], webopts.NodeID, webopts)
				fmt.Println("interacting...")
				nodes.SelectNode(webopts.NodeID)
				return

			}
			if linlocalopts.Use == true {
				scripts.ScriptSendandReturn(sep[1], linlocalopts.NodeID, linlocalopts)
				fmt.Println("interacting...")
				nodes.SelectNode(linlocalopts.NodeID)
				return
			}
			if linremoteopts.Use == true {
				scripts.ScriptSendandReturn(sep[1], linremoteopts.NodeID, linremoteopts)
				fmt.Println("interacting...")
				nodes.SelectNode(linremoteopts.NodeID)
				return
			}
		default:
			fmt.Println("use from scripts menu with a kudzu implant")
		}
	//execute element
	case "run", "execute":
		switch cltag {
		//execute script with given options
		case "<kudzu scripts> ":
			if len(sep) != 2 {
				fmt.Println("Run what? usage: run *.kzs")
				return
			}
			if winlocalopts.Use == true {
				scripts.ScriptRun(winlocalopts, sep...)
			}
			if winremoteopts.Use == true {
				scripts.ScriptRun(winremoteopts, sep...)
			}
			if webopts.Use == true {
				scripts.ScriptRun(webopts, sep...)
			}
			if linlocalopts.Use == true {
				scripts.ScriptRun(linlocalopts, sep...)
			}
			if linremoteopts.Use == true {
				scripts.ScriptRun(linremoteopts, sep...)
			}
			if compileandrun.Use == true {
				scripts.ScriptCompileAndRun(compileandrun, sep...)
			}
		//start listener with given options
		case "<kudzu nodes> ":
			if len(sep) != 1 {
				fmt.Println("run does not take arguments in this menu")
			} else {
				switch nodeopts.NodeType {
				case "tcp":
					for _, x := range nodes.Listeners {
						if nodeopts.Addr+":"+nodeopts.Port == x.Listener.Addr().String() {
							fmt.Println("Listener already started!")
							return
						}
					}
					go nodes.SetupTCPNode(nodeopts)
				case "tcp_tls":
					for _, x := range nodes.Listeners {
						if nodeopts.Addr+":"+nodeopts.Port == x.Listener.Addr().String() {
							fmt.Println("Listener already started!")
							return
						}
					}
					go nodes.SetupTCPEncNode(nodeopts)
				}
			}
		//generate implant with given options
		case "<kudzu implants> ":
			fmt.Printf("%+v\n", implantops)
			fmt.Printf("proceed? Y/N > ")
			checker := bufio.NewReader(os.Stdin)
			proceed, err := checker.ReadString('\n')
			if err != nil {
				log.Println(err)
			}
			switch strings.TrimSpace(proceed) {
			case "Y", "y", "yes":
				implants.GenerateImplant(implantops)
			default:
				fmt.Println("check options and try again")
			}

		}
	//display asset info
	case "load":
		switch cltag {
		case "<kudzu scripts> ":
			if len(sep) == 2 && sep[1] != "" {
				curscript = sep[1]
				winlocalopts, winremoteopts, webopts, linlocalopts, linremoteopts, compileandrun = scripts.GetJsonStruct(sep[1])
			}
		}
	//set options for element
	case "setop":
		if len(sep) < 2 {
			fmt.Println("use in nodes, scripts, or implants menu")
			return
		}
		switch cltag {
		//manage implant options struct
		case "<kudzu implants> ":
			if len(sep) == 3 && sep[2] != "" {
				switch sep[1] {
				case "ImplantType", "implanttype", "Implanttype":
					implantops.ImplantType = sep[2]
				case "FileName", "filename", "Filename":
					implantops.FileName = sep[2]
				case "Listener", "listener":
					for i := range nodes.Listeners {
						if nodes.Listeners[i].ID == sep[2] {
							implantops.Listener = nodes.Listeners[i]
						}
					}
				}
			}
		//manage script options struct
		case "<kudzu scripts> ":
			switch sep[1] {
			case "LHOST", "Lhost", "lhost":
				if winlocalopts.Use == true {
					winlocalopts.Lhost = sep[2]
				} else if winremoteopts.Use == true {
					winremoteopts.Lhost = sep[2]
				} else if webopts.Use == true {
					webopts.Lhost = sep[2]
				} else if linlocalopts.Use == true {
					linlocalopts.Lhost = sep[2]
				} else if linremoteopts.Use == true {
					linremoteopts.Lhost = sep[2]
				} else if compileandrun.Use == true {
					compileandrun.Lhost = sep[2]
				}

			case "RHOST", "Rhost", "rhost":
				if winlocalopts.Use == true {
					winlocalopts.Rhost = sep[2]
				} else if winremoteopts.Use == true {
					winremoteopts.Rhost = sep[2]
				} else if webopts.Use == true {
					webopts.Rhost = sep[2]
				} else if linlocalopts.Use == true {
					linlocalopts.Rhost = sep[2]
				} else if linremoteopts.Use == true {
					linremoteopts.Rhost = sep[2]
				} else if compileandrun.Use == true {
					compileandrun.Rhost = sep[2]
				}

			case "CMD", "cmd", "Cmd":
				if winlocalopts.Use == true {
					winlocalopts.Cmd = strings.Join(sep[2:], " ")
				} else if winremoteopts.Use == true {
					winremoteopts.Cmd = strings.Join(sep[2:], " ")
				} else if webopts.Use == true {
					webopts.Cmd = strings.Join(sep[2:], " ")
				} else if linlocalopts.Use == true {
					linlocalopts.Cmd = strings.Join(sep[2:], " ")
				} else if linremoteopts.Use == true {
					linremoteopts.Cmd = strings.Join(sep[2:], " ")
				} else if compileandrun.Use == true {
					compileandrun.Cmd = sep[2]
				}

			case "LPORT", "lport", "Lport":
				if winlocalopts.Use == true {
					winlocalopts.Lport = sep[2]
				} else if winremoteopts.Use == true {
					winremoteopts.Lport = sep[2]
				} else if webopts.Use == true {
					webopts.Lport = sep[2]
				} else if linlocalopts.Use == true {
					linlocalopts.Lport = sep[2]
				} else if linremoteopts.Use == true {
					linremoteopts.Lport = sep[2]
				} else if compileandrun.Use == true {
					compileandrun.Lport = sep[2]
				}
			case "RPORT", "rport", "Rport":
				if winlocalopts.Use == true {
					winlocalopts.Rport = sep[2]
				} else if winremoteopts.Use == true {
					winremoteopts.Rport = sep[2]
				} else if webopts.Use == true {
					webopts.Rport = sep[2]
				} else if linlocalopts.Use == true {
					linlocalopts.Rport = sep[2]
				} else if linremoteopts.Use == true {
					linremoteopts.Rport = sep[2]
				} else if compileandrun.Use == true {
					compileandrun.Rport = sep[2]
				}
			case "Domain", "domain", "DOMAIN":
				if winlocalopts.Use == true {
					winlocalopts.Domain = sep[2]
				} else if winremoteopts.Use == true {
					winremoteopts.Domain = sep[2]
				} else if compileandrun.Use == true {
					compileandrun.Domain = sep[2]
				}
			case "Username", "username", "USERNAME":
				if winlocalopts.Use == true {
					winlocalopts.Username = sep[2]
				} else if winremoteopts.Use == true {
					winremoteopts.Username = sep[2]
				} else if webopts.Use == true {
					webopts.Username = sep[2]
				} else if linlocalopts.Use == true {
					linlocalopts.Username = sep[2]
				} else if linremoteopts.Use == true {
					linremoteopts.Username = sep[2]
				} else if compileandrun.Use == true {
					compileandrun.Username = sep[2]
				}
			case "Password", "password", "PASSWORD":
				if winlocalopts.Use == true {
					winlocalopts.Password = sep[2]
				} else if winremoteopts.Use == true {
					winremoteopts.Password = sep[2]
				} else if webopts.Use == true {
					webopts.Password = sep[2]
				} else if linlocalopts.Use == true {
					linlocalopts.Password = sep[2]
				} else if linremoteopts.Use == true {
					linremoteopts.Password = sep[2]
				} else if compileandrun.Use == true {
					compileandrun.Password = sep[2]
				}
			case "NodeID", "Nodeid", "nodeid", "NODEID":
				if winlocalopts.Use == true {
					winlocalopts.NodeID = sep[2]
				} else if winremoteopts.Use == true {
					winremoteopts.NodeID = sep[2]
				} else if webopts.Use == true {
					webopts.NodeID = sep[2]
				} else if linlocalopts.Use == true {
					linlocalopts.NodeID = sep[2]
				} else if linremoteopts.Use == true {
					linremoteopts.NodeID = sep[2]
				} else if compileandrun.Use == true {
					compileandrun.NodeID = sep[2]
				}
			case "Hostname", "HOSTNAME", "HostName":
				if winlocalopts.Use == true {
					winlocalopts.Hostname = sep[2]
				} else if winremoteopts.Use == true {
					winremoteopts.Hostname = sep[2]
				} else if webopts.Use == true {
					webopts.Hostname = sep[2]
				} else if linlocalopts.Use == true {
					linlocalopts.Hostname = sep[2]
				} else if linremoteopts.Use == true {
					linremoteopts.Hostname = sep[2]
				} else if compileandrun.Use == true {
					compileandrun.Hostname = sep[2]
				}
			case "Dir", "directory", "Directory", "dir":
				if winlocalopts.Use == true {
					winlocalopts.Directory = sep[2]
				} else if winremoteopts.Use == true {
					winremoteopts.Directory = sep[2]
				} else if webopts.Use == true {
					webopts.Directory = sep[2]
				} else if linlocalopts.Use == true {
					linlocalopts.Directory = sep[2]
				} else if linremoteopts.Use == true {
					linremoteopts.Directory = sep[2]
				} else if compileandrun.Use == true {
					compileandrun.Directory = sep[2]
				}
			case "Filename", "filename", "FileName":
				if winlocalopts.Use == true {
					winlocalopts.Filename = sep[2]
				} else if winremoteopts.Use == true {
					winremoteopts.Filename = sep[2]
				} else if linlocalopts.Use == true {
					linlocalopts.Filename = sep[2]
				} else if linremoteopts.Use == true {
					linremoteopts.Filename = sep[2]
				} else if compileandrun.Use == true {
					compileandrun.Filename = sep[2]
				} else if webopts.Use == true {
					webopts.Filename = sep[2]
				}
			}
		//manage node options struct
		case "<kudzu nodes> ":
			switch sep[1] {
			case "nodetype", "NodeType":
				if (sep[2] == "tcp") || (sep[2] == "udp") || (sep[2] == "ntp") || (sep[2] == "tcp_tls") {
					nodeopts.NodeType = sep[2]
				} else {
					fmt.Println("available nodetypes: tcp, udp, ntp, tcp_tls")
				}
			case "address", "addr":
				if (len(sep) != 3) == true || (strings.Count(sep[2], ".") != 3) {
					fmt.Println("usage: setop Address 127.0.0.1")
				} else {
					nodeopts.Addr = sep[2]
				}
			case "port":
				if len(sep) != 3 {
					fmt.Println("usage: setop port 8080")
				} else {
					nodeopts.Port = sep[2]
				}
			}
		}
	case "close":
		switch cltag {
		//closes nodes/listeners
		case "<kudzu nodes> ":
			if len(sep) != 3 || sep[2] == "" {
				fmt.Println("usage: close node/listener <nodeID>")
			} else {
				switch sep[1] {
				case "node":
					nodes.CloseNode(sep[2])
				case "listener":
					nodes.CloseListener(sep[2])
				}
			}
		}
	case "interact":
		switch cltag {
		//interact with given node
		case "<kudzu nodes> ":
			if len(sep) == 2 && sep[1] != "" {
				fmt.Println("interacting...")
				nodes.SelectNode(sep[1])
				return
			}
		default:
			fmt.Println("use this in the nodes menu")
		}
	//print options in asset menu
	case "showops":
		switch cltag {
		case "<kudzu nodes> ":
			fmt.Printf("%+v\n", nodeopts)
		case "<kudzu scripts> ":
			switch true {
			case linlocalopts.Use:
				fmt.Printf("%+v\n", linlocalopts)
			case linremoteopts.Use:
				fmt.Printf("%+v\n", linremoteopts)
			case winlocalopts.Use:
				fmt.Printf("%+v\n", winlocalopts)
			case winremoteopts.Use:
				fmt.Printf("%+v\n", winremoteopts)
			case webopts.Use:
				fmt.Printf("%+v\n", webopts)
			case compileandrun.Use:
				fmt.Printf("%+v\n", compileandrun)
			}
		case "<kudzu implants> ":
			fmt.Printf("%+v\n", implantops)
		}
	//implant management panel
	case "implants":
		cltag = "<kudzu implants> "
	//script management panel
	case "scripts":
		cltag = "<kudzu scripts> "
	//node/listener management panel
	case "nodes":
		cltag = "<kudzu nodes> "
	//return to main menu
	case "back", "home":
		cltag = "<kudzu> "
	//leave kudzu?!
	case "exit":
		os.Exit(0)
	}
}
