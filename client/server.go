package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"reflect"
	"sync"
	"task-offloading/craft"
	"task-offloading/structs"
	"task-offloading/to"
	"task-offloading/util"
	"time"
)

type server struct{}

/* constant parameters */
const (
	// server address along with its port number
	ServerAddressDebug = "localhost:50051"
	Method             = "tcp"
	ColdDownPeriod     = time.Millisecond * 997
)

var serverAddressRun string

var ctxM context.Context
var cancel context.CancelFunc
var mutex sync.Mutex

var taskIns []structs.TaskIn
var taskOuts []structs.TaskOut

var tors []structs.Tor

/* get request and reply when processed */
func (s *server) GetChannel(_ context.Context, a *craft.Request) (*craft.Reply, error) {
	var taskIn structs.TaskIn
	req := a.Req
	err := json.Unmarshal(util.String2Bytes(req), &taskIn)
	if err != nil {
		return nil, err
	}
	taskOuts := concurrentTO(ctxM, taskIn)
	if taskOuts != nil {
		for _, taskOut := range taskOuts {
			if taskOut.ID == taskIn.ID {
				// set reply
				rep, err := json.Marshal(taskOut)
				if err != nil {
					return nil, err
				}
				res := craft.Reply{
					Rep: util.Bytes2String(rep),
				}
				fmt.Println("[gRPC server] send message:", res.String())
				return &res, nil
			}
		}
	}
	err = errors.New("[gRPC server] task seems lost")
	util.PrintError(err.Error())
	return nil, nil
}

func concurrentTO(ctx context.Context, taskIn structs.TaskIn) []structs.TaskOut {
	taskIns = append(taskIns, taskIn)
	for {
		select {
		case <-ctx.Done():
			mutex.Lock()
			if len(taskIns) != 0 {
				fmt.Println("[gRPC server] concurrency number: ", len(taskIns))
				taskOuts = to.Execute(taskIns)
				taskIns = taskIns[0:0]
				tors = genTor()
			}
			mutex.Unlock()
			return taskOuts
		}
	}
}

func genTor() []structs.Tor {
	var tors []structs.Tor
	taskMap := make(map[string][]structs.Task)
	for _, taskOut := range taskOuts {
		offloadTo := taskOut.OffloadTo
		if reflect.ValueOf(taskMap).MapIndex(reflect.ValueOf(offloadTo)).IsValid() {
			tasks := taskMap[offloadTo]
			flag := false
			for i, task := range tasks {
				if task.Type == taskOut.Type {
					task.Number = task.Number + 1
					tasks[i] = task
					flag = true
					break
				}
			}
			if !flag {
				tasks = append(tasks, structs.Task{
					Type:   taskOut.Type,
					Number: 1,
				})
			}
			taskMap[offloadTo] = tasks
		} else {
			tasks := []structs.Task{
				{
					Type:   taskOut.Type,
					Number: 1,
				},
			}
			taskMap[offloadTo] = tasks
		}
	}
	// taskMap -> Tor
	for ip, tasks := range taskMap {
		tors = append(tors, structs.Tor{
			IP:    ip,
			Tasks: tasks,
		})
	}
	fmt.Println("[Tor] TO Report:", tors)
	return tors
}

func StartGrpcServer(serverAddress string) {
	if len(serverAddress) != 0 {
		serverAddressRun = serverAddress
	} else {
		serverAddressRun = ServerAddressDebug
	}
	// listening network port
	listener, err := net.Listen(Method, serverAddressRun)
	if err != nil {
		return
	}
	// initialize public var
	taskIns = make([]structs.TaskIn, 0)
	taskOuts = make([]structs.TaskOut, 0)
	// initialize StartSlot
	con := context.Background()
	go StartSlot(con)
	// open gRPC service
	s := grpc.NewServer()
	craft.RegisterGetChannelServer(s, &server{})
	reflection.Register(s)
	fmt.Println("[gRPC server] started successfully")
	err = s.Serve(listener)
	if err != nil {
		return
	}
}

func StartSlot(con context.Context) {
	ctx, can := context.WithCancel(con)
	ctxM = ctx
	cancel = can
	for {
		time.Sleep(ColdDownPeriod)
		fmt.Println("[gRPC server] new time slot")
		cancel()
		ctx, can := context.WithCancel(context.Background())
		ctxM = ctx
		cancel = can
	}
}
