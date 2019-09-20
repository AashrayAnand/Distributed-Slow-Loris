package shared

import (
  "net"
)

type Headers struct {
      Base          []byte // base HTTP request string
      //General       []byte // typical HTTP headers
      Loris         []byte // header which will be repeatadly sent
}

type State struct {
      Endpoint      string
      Delay         int
      NumAttackers  int
      Timeout       int

}

// struct representing a slow loris attacker
type Attacker struct {
      Id            int
      Conn          net.Conn // network connection
      Active        int // 0 if connection is lost, 1 o.w.
      Errors        int // # of times connection has been lost
      Writes        int // # of writes to endpoint
}
