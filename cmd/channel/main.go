package main

import (
	"fmt"
	"github.com/pusher/pusher-http-go/v5"
)

//
//app_id = "1716290"
//key = "25cfd36d24966103d735"
//secret = "e34846163d875ae66d14"
//cluster = "mt1"
func main() {
	pusherClient := pusher.Client{
		AppID:   "1716290",
		Key:     "25cfd36d24966103d735",
		Secret:  "e34846163d875ae66d14",
		Cluster: "mt1",
	}

	data := map[string]string{"message": "hello world"}
	pusherClient.Trigger("test-channel-example", "my-event", data)

	fmt.Println("done")
}
