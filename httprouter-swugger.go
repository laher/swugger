package swugger

import (
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
)

type HttpRouterSwagger struct {
	HttpRouter *httprouter.Router
	GoRestfulContainer *restful.Container
	GoRestfulWebService *restful.WebService
	SwaggerConfig *swagger.Config
}
func (hrs *HttpRouterSwagger) Init(webServiceAddr string) {
	hrs.HttpRouter = httprouter.New()
	hrs.GoRestfulContainer = restful.NewContainer()
	hrs.GoRestfulWebService = new(restful.WebService)
	hrs.GoRestfulContainer.Add(hrs.GoRestfulWebService)
	hrs.SwaggerConfig = &swagger.Config{
		WebServices:    hrs.GoRestfulContainer.RegisteredWebServices(),
		WebServicesUrl: webServiceAddr,
		ApiPath:        "/doc/apidocs.json",
		// Optionally, specifiy where the UI is located
		SwaggerPath:     "/doc/apidocs/",
		SwaggerFilePath: "../swagger-ui/dist"}

}
func NewHRS(webServiceAddr string) *HttpRouterSwagger {
	hrs := &HttpRouterSwagger{}
	hrs.Init(webServiceAddr)
	return hrs
}


func (hrs *HttpRouterSwagger) AddRoute(method string, path string, function httprouter.Handle) *restful.RouteBuilder {
	hrs.HttpRouter.Handle(method, path, function)
	ws := hrs.GoRestfulWebService

	pathGoRestful := path
	pathParts := strings.Split(path, ":")
	if len(pathParts) > 1 {
		pathGoRestful = pathParts[0]
		for i, part := range pathParts {
			if i > 0 {
				pathPartParts := strings.SplitN(part,"/", 2)
				pathGoRestful = pathGoRestful + "{" + pathPartParts[0] + "}"
				if len(pathPartParts) > 1 {
					pathGoRestful = pathGoRestful + "/" + pathPartParts[1]
				}
			}
		}
	}
	//This is the ugly part. As it stands, go-restful requires that you need to actually have a route in order for its documentation to work.
	rb := ws.Method(method).Path(pathGoRestful).To(func (req *restful.Request, resp *restful.Response) {
		// just reject it.
		resp.WriteErrorString(http.StatusNotFound, "Page not found")
		})
	return rb
}

