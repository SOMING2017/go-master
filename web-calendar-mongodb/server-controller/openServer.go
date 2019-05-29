// openServer
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"../go-common-master/calendar"
)

var chanCloseServer = make(chan bool)

func main() {
	SetFileCloseServer(false)
	httpHandle()
	server := &http.Server{Addr: ":1234", Handler: nil}
	go listenerCloseHttp(server)
	fmt.Println("http server open...")
	logServeError(server.ListenAndServe())
	fmt.Println("http server close...")
}

func httpHandle() {
	http.Handle("/", http.FileServer(http.Dir("../html/")))
	http.HandleFunc("/api/calendar/controller", calendar.Start)
}

func logServeError(err error) {
	fmt.Println("server close...")
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	} else {
		fmt.Println("but not found error!")
	}
}

func listenerCloseHttp(srv *http.Server) {
	go listenerCloseFile()
	go listenerCloseCommand()
	b := <-chanCloseServer
	fmt.Println("close server state is : ", b)
	srv.Shutdown(nil)
	srv.Close()
}

func listenerCloseFile() {
	for true {
		fileCloseServer, err := GetFileCloseServer()
		if err == nil && fileCloseServer {
			fmt.Println("file close parm : ", fileCloseServer)
			break
		}
		time.Sleep(5 * time.Second)
	}
	chanCloseServer <- true
}

func listenerCloseCommand() {
	chanCloseCommand := make(chan os.Signal)
	signal.Notify(chanCloseCommand)
	s := <-chanCloseCommand
	fmt.Println("close command : ", s)
	chanCloseServer <- true
}
