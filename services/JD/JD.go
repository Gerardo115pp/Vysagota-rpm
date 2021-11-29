package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	prouter "github.com/Gerardo115pp/PatriotRouter"
	echo "github.com/gwen/putils/echo"
	vlibs "github.com/vysagota/libs"
)

type JD struct {
	router   *prouter.Router
	port     string
	host     string
	services map[string]*vlibs.Service
}

func (jd *JD) publish(response http.ResponseWriter, request *http.Request) {
	/*
	   Publishes domain events to all services that are registered to the event.
	   Gets an event object and publishes it.


	*/

	response.WriteHeader(200)
}

func (jd *JD) registerService(response http.ResponseWriter, request *http.Request) {
	/*

	   Any service that wants to register with JD must send a POST request to /new-service, and then
	   it will be  avaliable to be used by other services.
	   ------------------------------------------------------------------------------------------------
	   The request form must be a multipart/form-data with the following fields:
	       - name: The name of the service
	       - DNS: The DNS of the service
	       - port: The port of the service
	       - Protocol: The protocol of the service, can be more than one, separated by commas without spaces
	       - Actions: The actions of the service, is a JSON stringified object, with the format {action: {method: "GET", path: "/action", description: "This is an action"}}
	*/
	echo.Echo(echo.CyanFG, fmt.Sprintf("Registering service request from %s", request.RemoteAddr))

	err := request.ParseMultipartForm(32 << 20)
	if err != nil {
		echo.EchoErr(err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	var service *vlibs.Service = new(vlibs.Service)
	service.Name = request.FormValue("name")
	service.DNS = request.FormValue("dns")
	service.Port = request.FormValue("port")
	service.Host = request.FormValue("host")
	service.Status = "online"
	service.Protocols = strings.Split(request.FormValue("protocol"), ",")

	var actions map[string]*vlibs.Action

	// Debug
	echo.EchoDebug(fmt.Sprintf("actions = %s", request.FormValue("actions")))
	echo.EchoDebug(fmt.Sprintf("name = %s", request.FormValue("actions")))

	err = json.Unmarshal([]byte(request.FormValue("actions")), &actions)
	if err != nil {
		echo.EchoErr(err)
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	echo.Echo(echo.CyanFG, fmt.Sprintf("Service %s registered with actions:\n", service.Name))
	for action, actionData := range actions {
		echo.Echo(echo.CyanFG, fmt.Sprintf("\t%s: %s\n", action, actionData.Name))
	}

	service.Actions = actions
	jd.services[service.Name] = service

}

func (jd *JD) run() {
	echo.Echo(echo.CyanFG, fmt.Sprintf("JD is running on %s:%s", jd.host, jd.port))
	jd.router = prouter.CreateRouter()

	jd.router.RegisterRoute(prouter.NewRoute("/new-service", true), jd.registerService)
	jd.router.RegisterRoute(prouter.NewRoute("/publish", true), jd.publish)
	jd.router.RegisterRoute(prouter.NewRoute("/endpoints", true), jd.serveEndpoints)
	jd.router.RegisterRoute(prouter.NewRoute("/services/.+", false), jd.serveServices)

	jd.router.SetCorsHandler(prouter.CorsAllowAll)

	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", jd.host, jd.port), jd.router); err != nil {
		echo.EchoErr(err)
	}
}

func (jd *JD) serveEndpoints(response http.ResponseWriter, request *http.Request) {
	/*

		NOTE: for now all services are endpoints, but this will change in the future.
		this function will return a list of endpoints for a given service.

	*/
	echo.Echo(echo.PinkFG, fmt.Sprintf("Sending endpoints request to %s", request.RemoteAddr))
	response.Header().Set("Content-Type", "application/json")

	var response_data []byte
	var err error

	response_data, err = json.Marshal(jd.services)
	if err != nil {
		response.WriteHeader(500) // Internal Server Error
		return
	}

	response.Write(response_data)
	response.WriteHeader(200)
}

func (jd *JD) serveServices(response http.ResponseWriter, request *http.Request) {
	echo.Echo(echo.PinkFG, fmt.Sprintf("Sending services request to %s", request.RequestURI))
	var service_name string = strings.Replace(request.RequestURI, "/services/", "", 1)

	if service, exists := jd.services[service_name]; exists {
		service_data, content_type := service.ToMultipart()

		response.Header().Set("Content-Type", content_type)
		response.Write(service_data.Bytes())
		response.WriteHeader(200)
	} else {
		response.WriteHeader(404)
	}

}

func main() {
	var server *JD = new(JD)
	server.port = os.Getenv("JD_PORT")
	if server.port == "" {
		echo.EchoErr(fmt.Errorf("JD failed to inicialize, JD_PORT not set"))
		os.Exit(1)
	}

	server.host = os.Getenv("JD_HOST")
	if server.host == "" {
		echo.EchoErr(fmt.Errorf("JD failed to inicialize, JD_HOST not set"))
		os.Exit(1)
	}

	server.services = make(map[string]*vlibs.Service)

	server.run()
}
