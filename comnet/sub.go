package comnet

import (
	"context"
	"fmt"
	"github.com/go-zeromq/zmq4"
	"task-offloading/util"
	"time"
)

/* constant parameters */
const (
	// server address along with its port number
	PubAddressDebug = "localhost:5555"
	ZmqSubTopic     = "comnetTO"
	ZmqSubTestTopic = "comnetBackend"
	Method          = "tcp"
)

var sub zmq4.Socket

var pubAddressRun string

func StartZmqSubscriber(pubAddress string) {
	if len(pubAddress) != 0 {
		pubAddressRun = pubAddress
	} else {
		pubAddressRun = PubAddressDebug
	}
	pubAddressRun = Method + "://" + pubAddressRun
	//  Prepare our subscriber
	sub = zmq4.NewSub(context.Background())
	err := sub.Dial(pubAddressRun)
	if err != nil {
		fmt.Println("[ZMQ Sub] could not dial: ", err)
	} else {
		subscribe()
	}
}

func subscribe() {
	err := sub.SetOption(zmq4.OptionSubscribe, ZmqSubTopic)
	if err != nil {
		fmt.Println("[ZMQ Sub] could not subscribe: ", err)
	}
	fmt.Println("[ZMQ Sub] started successfully")
	// start receiving data
	for {
		msg, err := sub.Recv()
		if err != nil {
			// try to reconnect to the publisher
			fmt.Println("[ZMQ Sub] could not receive message: ", err)
			reconnect()
		} else {
			// store data
			StoreComNet(util.Bytes2String(msg.Frames[1]))
		}
	}
}

func reconnect() {
	for {
		fmt.Println("[ZMQ Sub] reconnecting... ")
		err := sub.Dial(pubAddressRun)
		if err != nil {
			fmt.Println("[ZMQ Sub] could not dial: ", err)
		} else {
			break
		}
	}
	subscribe()
}

func StartZmqPublisher() {
	pub := zmq4.NewPub(context.Background())
	defer pub.Close()
	err := pub.Listen(Method + "://" + PubAddressDebug)
	if err != nil {
		fmt.Printf("[ZMQ Pub] could not listen: %v\n", err)
	}
	// debug for TO
	str := `{
		"ID": "16024697834867",
		"Nodes": [
			{
				"IP": "10.128.254.5",
				"CGs": [
					{"Type": "Face Recognition", "Utility": [2, 1, 0]},
					{"Type": "Gesture Recognition", "Utility": [1, 1, 0]}
				],
				"RestComputing": [4, 2, 0.9],
				"RestStorage": [1000, 0.728]}, 
			{
				"IP": "10.128.248.253",
				"CGs": [
					{"Type": "Face Recognition", "Utility": [1, 0, 0]},
					{"Type": "Gesture Recognition", "Utility": [1, 0, 0]}
				],
				"RestComputing": [1, 2.5, 0.7],
				"RestStorage": [400, 0.114]
			}
		],
		"Links": [
			{
				"NodeFrom": "10.128.254.5",
				"NodeTo": "10.128.248.253",
				"Rate": 3.44,
				"EsDelay": 1
			}, {
				"NodeFrom": "10.128.248.253",
				"NodeTo": "10.128.254.5",
				"Rate": 4.82,
				"EsDelay": 1
			}
		]}`
	msgA := zmq4.NewMsgFrom(
		[]byte(ZmqSubTestTopic),
		[]byte("We don't want to see this"),
	)
	msgB := zmq4.NewMsgFrom(
		[]byte(ZmqSubTopic),
		[]byte(str),
	)
	for {
		fmt.Println("[ZMQ Pub] sending message", util.Bytes2String(msgA.Bytes()))
		err = pub.Send(msgA)
		if err != nil {
			fmt.Println("[ZMQ Pub] error:", err)
		}
		fmt.Println("[ZMQ Pub] sending message", util.Bytes2String(msgB.Bytes()))
		err = pub.Send(msgB)
		if err != nil {
			fmt.Println("[ZMQ Pub] error:", err)
		}
		time.Sleep(time.Second)
	}
}
