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
var scropts scripts.ScriptOps
var winlocalopts scripts.WinLocal
var winremoteopts scripts.WinRemote
var webopts scripts.Web
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
			fmt.Println("info: displays kudzu script info")
			fmt.Println("	usage: info hello.kzs")
			fmt.Println("ls/list: lists scripts available")
			fmt.Println("	usage: list")
		case "<kudzu implants> ":
			fmt.Println("setop: sets option for implant")
			fmt.Println("	options: implanttype (cmd/psh)")
			fmt.Println("		    listener (ID)")
			fmt.Println("		    filename (name of generated implant")
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
			fmt.Println("	options: nodetype (tcp/udp/ntp)")
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
				fmt.Println("Addr:", m.Listener.Addr().String()+"\n")
			}
		default:
			fmt.Println("use from agents, nodes, or scripts menu")
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
			scripts.ScriptRun(scropts, sep...)
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
				}
			}
		//generate implant with given options
		case "<kudzu implants> ":
			fmt.Println(implantops)
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
				winlocalopts, winremoteopts, webopts = scripts.GetJsonStruct(sep[1])
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
				}

			case "RHOST", "Rhost", "rhost":
				if winlocalopts.Use == true {
					winlocalopts.Rhost = sep[2]
				} else if winremoteopts.Use == true {
					winremoteopts.Rhost = sep[2]
				}
			case "CMD", "cmd":
				cmdarr := sep[2:]
				scropts.CMD = strings.Join(cmdarr, " ")
				if winlocalopts.Use == true {
					winlocalopts.Cmd = strings.Join(sep[2:], " ")
				} else if winremoteopts.Use == true {
					winremoteopts.Cmd = strings.Join(sep[2:], " ")
				}
			case "LPORT", "lport", "Lport":
				if winlocalopts.Use == true {
					winlocalopts.Lport = sep[2]
				} else if winremoteopts.Use == true {
					winremoteopts.Lport = sep[2]
				}
			case "RPORT", "rport", "Rport":
				if winlocalopts.Use == true {
					winlocalopts.Rport = sep[2]
				} else if winremoteopts.Use == true {
					winremoteopts.Rport = sep[2]
				}
			case "Domain", "domain", "DOMAIN":
				if winlocalopts.Use == true {
					winlocalopts.Domain = sep[2]
				} else if winremoteopts.Use == true {
					winremoteopts.Domain = sep[2]
				}
			case "Username", "username", "USERNAME":
				if winlocalopts.Use == true {
					winlocalopts.Username = sep[2]
				} else if winremoteopts.Use == true {
					winremoteopts.Username = sep[2]
				}
			case "Password", "password", "PASSWORD":
				if winlocalopts.Use == true {
					winlocalopts.Password = sep[2]
				} else if winremoteopts.Use == true {
					winremoteopts.Password = sep[2]
				}
			case "NodeID", "Nodeid", "nodeid", "NODEID":
				if winlocalopts.Use == true {
					winlocalopts.NodeID = sep[2]
				} else if winremoteopts.Use == true {
					winremoteopts.NodeID = sep[2]
				}
			case "Hostname", "HOSTNAME", "HostName":
				if winlocalopts.Use == true {
					winlocalopts.Hostname = sep[2]
				} else if winremoteopts.Use == true {
					winremoteopts.Hostname = sep[2]
				}
			}
		//manage node options struct
		case "<kudzu nodes> ":
			switch sep[1] {
			case "nodetype", "NodeType":
				if (sep[2] == "tcp") || (sep[2] == "udp") || (sep[2] == "ntp") {
					nodeopts.NodeType = sep[2]
				} else {
					fmt.Println("available nodetypes: tcp, udp, ntp")
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
				nodes.InteractNode(sep[1])
				return
			}
		default:
			fmt.Println("use this in the nodes menu")
		}
	//print options in asset menu
	case "showops":
		switch cltag {
		case "<kudzu nodes> ":
			fmt.Println("Nodetype:", nodeopts.NodeType)
			fmt.Println("Address:", nodeopts.Addr)
			fmt.Println("Port:", nodeopts.Port)
		case "<kudzu scripts> ":
			fmt.Println("lhost:", scropts.LHOST)
			fmt.Println("lport:", scropts.LPORT)
			fmt.Println("rhost:", scropts.RHOST)
			fmt.Println("rport:", scropts.RPORT)
			fmt.Println("cmd:", scropts.CMD)
			fmt.Println("testing--")
			fmt.Printf("winlocal: %+v\n", winlocalopts)
			fmt.Printf("winremote: %+v\n", winremoteopts)
		case "<kudzu implants> ":
			fmt.Println("Filename", implantops.FileName)
			fmt.Println("ImplantType:", implantops.ImplantType)
			fmt.Println("Listener ID:", implantops.Listener)
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
