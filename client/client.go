package client

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"strconv"
	"task-offloading/craft"
	"task-offloading/structs"
	"task-offloading/util"
	"time"
)

/* constant parameters */
const (
	AddressClient = "localhost:50051"
)

var ArriveAts = [2]string{"10.128.254.5", "10.128.254.5"}
var Types = [2]string{"Face Recognition", "Gesture Recognition"}

/* comNet function for client */
func StartGrpcClient(period int) {
	// set a 3-second timer
	t := time.NewTicker(2 * time.Second)
	for range t.C {
		tag := "[gRPC client" + strconv.Itoa(period) + "]"
		// connect to gRPC service
		conn, err := grpc.Dial(AddressClient, grpc.WithInsecure())
		if err != nil {
			fmt.Println(tag, " error: ", err)
			return
		}
		defer conn.Close()
		c := craft.NewGetChannelClient(conn)
		// set the request message
		bytes, err := json.Marshal(structs.TaskIn{
			ID:       strconv.Itoa(period),
			ArriveAt: ArriveAts[period%2],
			Type:     Types[period%2],
		})
		fmt.Println(tag, " send message: ", util.Bytes2String(bytes))
		req := craft.Request{
			Req: util.Bytes2String(bytes),
		}
		r, err := c.GetChannel(context.Background(), &req)
		if err != nil {
			fmt.Println(tag, " error: ", err)
			return
		}
		// pint the reply message
		fmt.Println(tag, " recv message: ", r)
	}
}
