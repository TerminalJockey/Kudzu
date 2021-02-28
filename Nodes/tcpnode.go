package nodes

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

//SetupTCPNode initializes a listener and sends connections for handling
func SetupTCPNode(opts NodeOpts) {
	tcpln, err := net.Listen("tcp", opts.Addr+":"+opts.Port)
	if err != nil {
		log.Println(err)
	}
	newlistener := Listener{
		ID:       GenUID(),
		Listener: tcpln,
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

//InteractNode provides management for remote sessions
func InteractNode(ID string) {
	for i := range Nodes {
		fmt.Println(Nodes[i].NodeOpts.ID)
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
				_, err = Nodes[i].Conn.Write([]byte(clinput))
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
