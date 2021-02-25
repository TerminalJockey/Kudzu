package nodes

import (
	"crypto/sha1"
	"fmt"
	"net"
	"time"
)

//NodeOpts contains options for nodes
type NodeOpts struct {
	NodeType, Addr, Port, ID string
}

//Node holds a node object
type Node struct {
	Conn     net.Conn
	NodeOpts NodeOpts
}

//Nodes is our node array
var Nodes []Node

//Listener holds our listener info
type Listener struct {
	ID       string
	Listener net.Listener
}

//Listeners is our listener array
var Listeners []Listener

//GenUID creates pseudorandom 8 char ids
func GenUID() (uid string) {
	t := time.Now()
	str := t.String()
	h := sha1.New()
	h.Write([]byte(str))
	uid = fmt.Sprintf("%x", h.Sum(nil))
	return string([]byte(uid[:10]))
}
