package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"../go-common-master/im/channel"
	"../go-common-master/im/channel/message"
	"../go-common-master/user"
	"../go-common-master/user/password"
)

var WebPrefix = "/api"

func main() {
	fileInit()
	httpHandle()
	server := &http.Server{Addr: ":1235", Handler: nil}
	go listenerCloseHttp(server)
	logServeError(server.ListenAndServe())
}

func fileInit() {
	fmt.Println("//************************************************")
	SetFileServerRunState(true)
}

func httpHandle() {
	http.Handle("/", http.FileServer(http.Dir("../html/")))
	http.HandleFunc(WebPrefix+user.ApiPath, user.Start)
	http.HandleFunc(WebPrefix+password.ApiPath, password.Start)
	http.HandleFunc(WebPrefix+channel.ApiPath, channel.Start)
	http.HandleFunc(WebPrefix+message.ApiPath, message.Start)
}

func logServeError(err error) {
	if err != nil {
		fmt.Println("server close,ListenAndServe : ", err)
	} else {
		fmt.Println("server close,but not found error!")
	}
	fmt.Println("//************************************************")
}

func listenerCloseHttp(srv *http.Server) {
	var chanCloseServer = make(chan bool)
	go listenerFile(chanCloseServer)
	go listenerCommand(chanCloseServer)
	<-chanCloseServer
	fmt.Println("chan : active close server")
	SetFileServerRunState(false)
	srv.Shutdown(nil)
	srv.Close()
}

func listenerFile(chanCloseServer chan bool) {
	for true {
		fileServerRunState, err := GetFileServerRunState()
		if err == nil && !fileServerRunState {
			fmt.Println("file : server run state parm : ", fileServerRunState)
			break
		}
		time.Sleep(5 * time.Second)
	}
	SafeSendMessage(chanCloseServer, true)
}

func listenerCommand(chanCloseServer chan bool) {
	chanCloseCommand := make(chan os.Signal)
	signal.Notify(chanCloseCommand)
	s := <-chanCloseCommand
	fmt.Println("close command : ", s)
	SafeSendMessage(chanCloseServer, true)
}

type T = bool

func SafeSendMessage(ch chan T, value T) {
	select {
	case ch <- value:
	default:
		//no deal
	}
}
