package main

import (
	"fmt"
	"github.com/laher/swugger"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"mime"

)

func generalGreeting(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func personalGreeting(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

//Can't seem to use string as a type. Use a struct for the time being
type greeting struct {}


func main() {
	httpRouter := httprouter.New()
	hrs := swugger.NewHRS("http://localhost:8080", httpRouter)
	//I've amended the api to use 'Doc' structs. It feels more 'go'-like
	ws := hrs.AddService("/", swugger.ServiceDoc{"Hello API", []string{mime.TypeByExtension(".xml"), mime.TypeByExtension(".json")}, []string{mime.TypeByExtension(".xml"), mime.TypeByExtension(".json")}})
	//this would preferably be a call on 'ws'
	hrs.AddRoute(ws, "GET", "/hello", generalGreeting, swugger.MethodDoc {
		Operation: "generalGreeting",
		Doc: "General greeting",
		Writes: greeting{}})
	hrs.AddRoute(ws, "GET", "/hello/:name", personalGreeting, swugger.MethodDoc {
		Operation: "personalGreeting",
		Doc: "personal greeting, with a name",
		Params: []swugger.ParamDoc{ swugger.ParamDoc{Name:"name", Doc:"identifier of the user",
			DataType: "string"}},
			Writes: greeting{} })
	http.Handle("/doc/", hrs.GetSwaggerHandler())

	http.Handle("/", hrs.HttpRouter)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

