
//                       _ooOoo_
//                      o8888888o
//                      88" . "88
//                      (| -_- |)
//                      O\  =  /O
//                    ___/`___'\___
//                  .' \\|     |// '.
//                 / \\|||  :  |||// \
//                / _||||| _:_ |||||- \
//               |   | \\\  -  /// |   |
//               | \_|  ''\---/''  |_/ |
//               \  .-\__  '-'  ___/-. /
//            ___'. .'   /--.--\   `. .'___
//          .""  '< `.___\_<|>_/___.' >'  "".
//         | | :   `_ \`.:`\_/`:.`/ - ` : | |
//         \  \ `_.    \_ __\/__ _/  .-` /  /
//     =====`-.____`.___ \_____/___.-`___.-'=====
//                       `=---='
//
//  ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
//
//			Buddha shines     Bugs disappear ever
//

/**
* 这种做发更像是队列
*/

package main

import (
	"fmt"
	"flag"
	"github.com/anachronistic/apns"
	"runtime"
)

var WORKERS = runtime.NumCPU()

func main () {
	runtime.GOMAXPROCS(runtime.NumCPU())
	begin()
}

func begin() {

	var message *string 	= flag.String("m", "hello world", "take your message with -m or --m")
	var pemFilePath *string = flag.String("x", "/Users/crosstime1986/Sites/__GO__/test/apns-dev.pem", "take your pem file with -x or --x")
	flag.Parse()

	tokens := []string {
		"004e22c06cf1438f753dca5daf85869840c80f7c0c2c0f376466f1270f9cedfa",
	}

	job := make(chan Job, WORKERS);
	result := make(chan JobResult, 1000);
	done :=  make(chan bool, WORKERS);

	go addJob(job, tokens, result)
	for i := 0; i < WORKERS; i++ {
		go doSendJob(done, job, fmt.Sprintf("[线程%d] %s", i, *message), *pemFilePath)
	}
	endJob(done, result)
}


func addJob (job chan<- Job, tokes []string, result chan<- JobResult)  {
	for _, toke := range tokes {
		job <- Job{toke, result, nil}
	}
	close(job)
}

func doSendJob (done chan<-bool, jobs <- chan Job, message string, xpath string) {

	client := apns.NewClient("gateway.sandbox.push.apple.com:2195", xpath, xpath)

	for job := range jobs {
		job.client = client
		job.send(message)
	}
	done <- true;
}

func endJob(done <-chan bool, result <-chan JobResult) {

	for work := WORKERS ; work > 0; {
		select {

		case <- done:
			fmt.Printf("---\n")
			work--
		case re := <- result:
			fmt.Println(re)
		}
	}

	select {
	case re := <- result:
		fmt.Println(re)
	default:
		return;
	}
}



type Job struct {
	token		string
	result		chan<- JobResult
	client 		*apns.Client
}

func (job Job) send (message string) {

	payload := apns.NewPayload()
	payload.Alert = message
	payload.Badge = 1
	payload.Sound = "fa9977e71e1f2e84cfc57a2ba1197c5b.mp3" 					//"fa9977e71e1f2e84cfc57a2ba1197c5b.mp3"

	pn := apns.NewPushNotification()
	pn.DeviceToken = job.token
	pn.AddPayload(payload)

	resp := job.client.Send(pn)
	alert, _ := pn.PayloadString()

	job.result <- JobResult{to: resp.Success, alert: alert}
}

type JobResult struct {
	alert	string
	to		bool
}
