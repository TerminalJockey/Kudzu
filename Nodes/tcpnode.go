package nodes

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
)

//SetupTCPNode initializes a listener and sends connections for handling
func SetupTCPNode(opts NodeOpts) {
	tcpln, err := net.Listen("tcp", opts.Addr+":"+opts.Port)
	if err != nil {
		log.Println(err)
	}
	newlistener := Listener{
		ID:           GenUID(),
		ListenerType: "tcp",
		Listener:     tcpln,
	}
	Listeners = append(Listeners, newlistener)
	for {
		tcpconn, err := tcpln.Accept()
		if err != nil {
			log.Println(err)
			break
		}
		go TCPNode(tcpconn, opts)
	}
}

//SetupTCPEncNode takes options and sets up tls listener
func SetupTCPEncNode(opts NodeOpts) {
	newid := GenUID()
	keyname, certname := GenCerts(newid, opts.Addr)
	if keyname == "" || certname == "" {
		return
	}
	keybytes, _ := os.ReadFile("certs/" + keyname)
	certbytes, _ := os.ReadFile("certs/" + certname)
	cert, err := tls.X509KeyPair(certbytes, keybytes)
	if err != nil {
		log.Println(err)
	}
	csuites := []uint16{0x0005, 0x000a, 0x002f, 0x0035, 0x003c, 0x009c, 0x009d, 0xc007, 0xc009, 0xc00a, 0xc011, 0xc012, 0xc014, 0xc013, 0xc023, 0xc027, 0xc02f, 0xc02b, 0xc030, 0xc02c, 0xcca8, 0xcca9}

	tcpencln, err := tls.Listen("tcp", opts.Addr+":"+opts.Port, &tls.Config{Certificates: []tls.Certificate{cert},
		MinVersion: tls.VersionTLS11, MaxVersion: tls.VersionTLS11, InsecureSkipVerify: true, CipherSuites: csuites})
	if err != nil {
		log.Println(err)
	}
	newlistener := Listener{
		ID:           newid,
		ListenerType: "tcp_tls",
		Listener:     tcpencln,
	}
	Listeners = append(Listeners, newlistener)
	for {
		tcpconn, err := tcpencln.Accept()
		if err != nil {
			log.Println(err)
			break
		}

		go TCPEncNode(tcpconn, opts)
	}
}

//TCPEncNode manages tls connections and adds to array
func TCPEncNode(tcpconn net.Conn, opts NodeOpts) {
	err := tcpconn.(*tls.Conn).Handshake()
	if err != nil {
		log.Println(err)
	}
	opts.ID = GenUID()
	fmt.Printf("Got Connection ID: %s %s:%s", opts.ID, opts.Addr, opts.Port)
	innode := Node{
		Conn:     tcpconn,
		NodeOpts: opts,
	}
	Nodes = append(Nodes, innode)
}

//TCPNode assigns ID to sessions, appends them to session list
func TCPNode(tcpconn net.Conn, opts NodeOpts) {
	opts.ID = GenUID()
	fmt.Printf("Got Connection ID: %s %s:%s", opts.ID, opts.Addr, opts.Port)
	innode := Node{
		Conn:     tcpconn,
		NodeOpts: opts,
	}
	Nodes = append(Nodes, innode)
}

//SelectNode allows node selection by id with differentiation from console without specifying node type
func SelectNode(ID string) {
	for i := range Nodes {
		if Nodes[i].NodeOpts.ID == ID {
			switch Nodes[i].NodeOpts.NodeType {
			case "tcp":
				InteractNode(ID)
			case "tcp_tls":
				InteractEncNode(ID)
			}
		}
	}
}

func GetListener(ID string) Listener {
	for i := range Listeners {
		if Listeners[i].ID == ID {
			return Listeners[i]
		}
	}
	return Listener{}
}

//InteractEncNode gets tls encrypted node by id for interaction
func InteractEncNode(ID string) {
	for i := range Nodes {
		if Nodes[i].NodeOpts.ID == ID {
			nodein := bufio.NewReader(os.Stdin)
		out:
			for {

				go func() {
					_, err := io.Copy(os.Stdout, Nodes[i].Conn)
					if err != nil {
						return
					}

				}()
				clinput, err := nodein.ReadString('\n')
				if err != nil {
					log.Println("cli read err:", err)
				}
				clinput = strings.TrimSpace(clinput)
				sep := strings.Split(strings.TrimSpace(clinput), " ")
				if len(sep) == 1 && sep[0] == "kdz_bg" {
					break out
				}
				if len(sep) == 1 && sep[0] == "kdz_exit" {
					Nodes[i].Conn.Write([]byte(clinput))
					CloseNode(ID)
					break out
				}
				if sep[0] == "runscript" {
					if len(sep) == 2 && sep[1] != "" && (strings.HasSuffix(sep[1], ".kzs") == true) {
						scriptbytes, err := ioutil.ReadFile("Scripts/" + sep[1])
						if err != nil {
							log.Println(err)
							fmt.Println("usage: runscript *.kzs")
							continue
						}
						encodedscript := base64.RawStdEncoding.EncodeToString(scriptbytes)
						_, err = Nodes[i].Conn.Write([]byte(clinput + "\n"))
						_, err = Nodes[i].Conn.Write([]byte(encodedscript + "\n"))
						continue
					}
					fmt.Println("usage: runscript *.kzs")
					continue
				}
				_, err = Nodes[i].Conn.Write([]byte(clinput + "\n"))
				if err != nil {
					break out
				}
			}
			fmt.Println("Done")
			return
		}
	}
	fmt.Println("Node not found")
	return
}

//GenCerts is os aware and will either use wsl or native openssl to generate self signed certs output to certs dir.
func GenCerts(ID, addr string) (keyname, certname string) {
	gpath := os.Getenv("GOPATH")
	var rebuild string
	var kdzpath string
	if strings.Contains(gpath, "\\") == true {
		gsep := strings.Split(gpath, "\\")
		rebuild = "/" + strings.Join(gsep[1:], "/")
		kdzpath = "/mnt/c" + rebuild + "/src/github.com/TerminalJockey/Kudzu/certs/"

		srvkeycmd := exec.Command("wsl", "openssl", "req", "-new", "-newkey", "rsa:2048", "-nodes",
			"-days", "365", "-x509", "-addext", "subjectAltName = IP:"+addr, "-subj", "/CN=test",
			"-keyout", kdzpath+ID+".key", "-out", kdzpath+ID+".cert")
		err := srvkeycmd.Run()
		if err != nil {
			if strings.HasSuffix(err.Error(), `executable file not found in %PATH%`) == true {
				fmt.Println("install wls for tls support on Windows (or implement your own!)")
				return
			}
			log.Println(err)
			return
		}
		return ID + ".key", ID + ".cert"
	} else {
		rebuild = "/" + gpath
		kdzpath = rebuild + "/github.com/TerminalJockey/Kudzu/certs/"
		srvkeycmd := exec.Command("openssl", "req", "-new", "-newkey", "rsa:2048", "-nodes",
			"-days", "365", "-x509", "-addext", "subjectAltName = IP:"+addr, "-subj", "/CN=test",
			"-keyout", kdzpath+ID+".key", "-out", kdzpath+ID+".cert")
		err := srvkeycmd.Run()
		if err != nil {
			if strings.HasSuffix(err.Error(), `executable file not found in %PATH%`) == true {
				fmt.Println("install wls for tls support (or implement your own!)")
				return
			}
			log.Println(err)
			return
		}
		return ID + ".key", ID + ".cert"
	}
}

//InteractNode provides management for remote sessions
func InteractNode(ID string) {
	for i := range Nodes {
		if Nodes[i].NodeOpts.ID == ID {
			fmt.Println("found node for interaction")
			nodein := bufio.NewReader(os.Stdin)

		out:
			for {

				go func() {
					_, err := io.Copy(os.Stdout, Nodes[i].Conn)
					if err != nil {
						return
					}

				}()

				clinput, err := nodein.ReadString('\n')
				if err != nil {
					log.Println("cli read err:", err)
				}
				clinput = strings.TrimSpace(clinput)
				sep := strings.Split(strings.TrimSpace(clinput), " ")
				if len(sep) == 1 && sep[0] == "kdz_bg" {
					break out
				}
				if len(sep) == 1 && sep[0] == "kdz_exit" {
					Nodes[i].Conn.Write([]byte(clinput))
					CloseNode(ID)
					break out
				}
				if sep[0] == "runscript" {
					if len(sep) == 2 && sep[1] != "" && (strings.HasSuffix(sep[1], ".kzs") == true) {
						scriptbytes, err := ioutil.ReadFile("Scripts/" + sep[1])
						if err != nil {
							log.Println(err)
							fmt.Println("usage: runscript *.kzs")
							continue
						}
						encodedscript := base64.RawStdEncoding.EncodeToString(scriptbytes)
						_, err = Nodes[i].Conn.Write([]byte(clinput))
						_, err = Nodes[i].Conn.Write([]byte(encodedscript + "\n"))
						continue
					}
					fmt.Println("usage: runscript *.kzs")
					continue
				}
				_, err = Nodes[i].Conn.Write([]byte(clinput + "\n"))
				if err != nil {
					break out
				}
			}
			fmt.Println("Done")
			return
		}

	}
	fmt.Println("Node not found")
	return
}

//CloseNode closes given node
func CloseNode(ID string) {
	for i := range Nodes {
		if Nodes[i].NodeOpts.ID == ID {
			Nodes[i].Conn.Close()
			fmt.Println("Closed node:", Nodes[i])
			copy(Nodes[i:], Nodes[i+1:])
			if len(Nodes) > 1 {
				Nodes = Nodes[:len(Nodes)-1]
			} else {
				Nodes = []Node{}
			}
			return
		}
	}
	fmt.Println("Node not found, check ID")
}

//CloseListener closes the listener provided by ID, and cleans array of listeners
func CloseListener(ID string) {
	for i := range Listeners {
		if Listeners[i].ID == ID {
			Listeners[i].Listener.Close()
			fmt.Println("Closed node:", Listeners[i])
			copy(Listeners[i:], Listeners[i+1:])
			if len(Listeners) > 1 {
				Listeners = Listeners[:len(Listeners)-1]
			} else {
				Listeners = []Listener{}
			}
			return
		}
	}
	fmt.Println("Listener not found, check ID")
}
