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
	//GoRestfulWebServices []*restful.WebService
	SwaggerConfig *swagger.Config
	isSwaggerServiceRegistered bool
}

//Not sure if this should have defaults or not, but it feels like there should be recommended swagger uri accross implementations.
func (hrs *HttpRouterSwagger) Init(webServiceAddr string) {
	hrs.GoRestfulContainer = restful.NewContainer()
	hrs.SwaggerConfig = &swagger.Config{
		WebServicesUrl: webServiceAddr,
		ApiPath:        "/doc/apidocs.json",
		// Optionally, specifiy where the UI is located
		SwaggerPath:     "/doc/apidocs/",
		SwaggerFilePath: "../swagger-ui/dist"}

}
func NewHRS(webServiceAddr string, httpRouter *httprouter.Router) *HttpRouterSwagger {
	hrs := &HttpRouterSwagger{}
	hrs.Init(webServiceAddr)
	hrs.HttpRouter = httpRouter
	return hrs
}

//wrapper for restful.WebService
type WebService struct {
	*restful.WebService
	Hrs *HttpRouterSwagger
	Path string
}

func (hrs *HttpRouterSwagger) AddService(path string, serviceDoc ServiceDoc) *WebService {
	ws := new(restful.WebService)
	mws := &WebService{ws, hrs, path}
	ws.Path(path).
		Doc(serviceDoc.Doc).
		Consumes(serviceDoc.Consumes...).
		Produces(serviceDoc.Produces...)
	hrs.GoRestfulContainer.Add(ws)
	//update
	hrs.SwaggerConfig.WebServices =  hrs.GoRestfulContainer.RegisteredWebServices()

	return mws
}

func (hrs *HttpRouterSwagger) GetSwaggerHandler() http.Handler {
	if !hrs.isSwaggerServiceRegistered {
		hrs.RegisterSwaggerService()
		hrs.isSwaggerServiceRegistered = true
	}
	return hrs.GoRestfulContainer
}

func (hrs *HttpRouterSwagger) RegisterSwaggerService() {
	swagger.RegisterSwaggerService(*hrs.SwaggerConfig, hrs.GoRestfulContainer)
}

func (hrs *HttpRouterSwagger) AddRoute(ws *WebService, method string, path string, function httprouter.Handle, methodDoc MethodDoc) *restful.RouteBuilder {
	return ws.AddRoute(method, path, function, methodDoc)
}

func (ws *WebService) AddRoute(method string, path string, function httprouter.Handle, methodDoc MethodDoc) *restful.RouteBuilder {
	fullPath := ws.Path + path
	ws.Hrs.HttpRouter.Handle(method, fullPath, function)
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
	if methodDoc.Operation != "" {
		rb.Operation(methodDoc.Operation)
	}
	if methodDoc.Doc != "" {
		rb.Doc(methodDoc.Doc)
	}
	if methodDoc.Writes != nil {
		rb.Writes(methodDoc.Writes)
	}
	if methodDoc.Reads != nil {
		rb.Reads(methodDoc.Reads)
	}
	if methodDoc.Params != nil {
		for _, p := range methodDoc.Params {
			switch p.Type {
			case "header":
				param := ws.HeaderParameter(p.Name, p.Doc).DataType(p.DataType).Required(true)
				rb.Param(param)
			case "query":
				param := ws.QueryParameter(p.Name, p.Doc).DataType(p.DataType)
				rb.Param(param)
			case "form":
				param := ws.FormParameter(p.Name, p.Doc).DataType(p.DataType)
				rb.Param(param)
			default:
				param := ws.PathParameter(p.Name, p.Doc).DataType(p.DataType)
				rb.Param(param)
			}
		}
	}
	ws.Route(rb)
	return rb
}

