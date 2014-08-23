package main

import (
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"github.com/laher/swugger"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"

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
	hrs := swugger.NewHRS("http://localhost:8080")
	hrs.GoRestfulWebService.
		Path("/").
		Doc("Hello API").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)
	rb := hrs.AddRoute("GET", "/hello", generalGreeting).
		Operation("generalGreeting").
		Doc("General greeting").
		Writes(greeting{})
	hrs.GoRestfulWebService.Route(rb)
	rb = hrs.AddRoute("GET", "/hello/:name", personalGreeting).
		Operation("personalGreeting").
		Doc("personal greeting, with a name").
		Param(hrs.GoRestfulWebService.PathParameter("name", "identifier of the user").DataType("string")).
		Writes(greeting{})
	hrs.GoRestfulWebService.Route(rb)
	swagger.RegisterSwaggerService(*hrs.SwaggerConfig, hrs.GoRestfulContainer)
	http.Handle("/doc/", hrs.GoRestfulContainer)
	http.Handle("/", hrs.HttpRouter)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

