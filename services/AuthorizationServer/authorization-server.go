package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	prouter "github.com/Gerardo115pp/PatriotRouter"
	echo "github.com/gwen/putils/echo"
	vlibs "github.com/vysagota/libs"
)

type AuthServer struct {
	router       *prouter.Router
	host         string
	port         string
	internal_dns string
	services     map[string]*vlibs.Service
}

func (auth *AuthServer) boot() {

	/*

	   Registers itself and its actions to JD, creates a service objects and all the necessary
	   Actions to be used by the other services.

	*/
	echo.Echo(echo.GreenFG, "Booting Authorization Server")

	var service *vlibs.Service = new(vlibs.Service)

	service.Name = "authorization"
	service.Host = auth.host
	service.Port = auth.port
	service.DNS = "authorization"
	service.Actions = make(map[string]*vlibs.Action)

	service.Actions["isJWTvalid"] = &vlibs.Action{
		Name:    "isJWTvalid",
		Methods: []string{"GET"},
	}
	service.Actions["login"] = &vlibs.Action{
		Name:    "login",
		Methods: []string{"POST"},
	}

	service_form, content_type := service.ToMultipart()
	var client *http.Client = new(http.Client)

	request, _ := http.NewRequest("POST", fmt.Sprintf("%s/new-service", vlibs.JD_ADDRESS), service_form)
	request.Header.Set("Content-Type", content_type)
	response, err := client.Do(request)
	if err != nil {
		echo.Echo(echo.OrangeFG, "Cant connect to JD node, no reason to stay up")
		echo.EchoFatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		echo.Echo(echo.OrangeFG, "JD responded a non-2xx, boot process failed, no reason to stay up")
		echo.EchoFatal(err) // kills the server
	}

	echo.Echo(echo.GreenFG, "JD responded with a 2xx")

	err = auth.loadAccountsServiceData()
	if err != nil {
		echo.Echo(echo.OrangeFG, "Accounts service is not up.")
	} else {
		echo.Echo(echo.GreenFG, "Accounts service data loaded")
	}

	echo.Echo(echo.GreenFG, "Authorization Server is up")
}

func (auth *AuthServer) handleLogin(response http.ResponseWriter, request *http.Request) {
	echo.Echo(echo.CyanFG, "Handling login request")

	request.ParseMultipartForm(32 << 20)

	var username string = request.FormValue("username")
	var password string = request.FormValue("password")
	echo.Echo(echo.CyanFG, "Username: ", username)
	echo.Echo(echo.CyanFG, "Password: ", password)

	var account_data []byte
	account_data, err := auth.requestAccountData(username)
	if err != nil {
		echo.EchoErr(err)
		response.WriteHeader(500)
		return
	}

	account_type, err := auth.getAccountType(account_data)
	if err != nil {
		echo.EchoErr(err)
		response.WriteHeader(500)
		return
	}

	// account_type can either be "patient" or "doctor" nothign else.
	echo.Echo(echo.CyanFG, "Account type: ", account_type)

	if account_type == "pacient" {
		//PACIENT

		var pacient *vlibs.Pacient

		account_holder := &struct {
			UserData *vlibs.Pacient `json:"user_data"`
			UserType string         `json:"user_type"`
		}{}

		err = json.Unmarshal(account_data, &account_holder)
		if err != nil {
			echo.EchoErr(err)
			response.WriteHeader(500)
			return
		}

		pacient = account_holder.UserData

		password_bytes := sha256.Sum256([]byte(password))
		password_hash := hex.EncodeToString(password_bytes[:])

		response.WriteHeader(200)
		if pacient.Password == password_hash {
			response.Write([]byte(fmt.Sprintf("{\"status\": true, \"type\": \"%s\"}", account_type)))
		} else {
			response.Write([]byte("{\"status\": false}"))
		}
	} else {
		//DOCTOR
		var doctor *vlibs.Doctor

		account_holder := &struct {
			UserData *vlibs.Doctor `json:"user_data"`
			UserType string        `json:"user_type"`
		}{}

		err = json.Unmarshal(account_data, &account_holder)
		if err != nil {
			echo.EchoErr(err)
			response.WriteHeader(500)
			return
		}

		doctor = account_holder.UserData

		password_bytes := sha256.Sum256([]byte(password))
		password_hash := hex.EncodeToString(password_bytes[:])

		response.WriteHeader(200)
		if doctor.Password == password_hash {
			response.Write([]byte(fmt.Sprintf("{\"status\": true, \"type\": \"%s\"}", account_type)))
		} else {
			response.Write([]byte("{\"status\": false}"))
		}
	}

}

func (auth *AuthServer) getAccountType(account_data []byte) (string, error) {
	// this parses the response from  accounts service which should be a json object with the
	// format {"user_data": {...}, "user_type": "pacient"|"doctor"}. this function
	// returns the user-type

	var account_type string
	var err error

	var account_data_map map[string]interface{}
	err = json.Unmarshal(account_data, &account_data_map)
	if err != nil {
		return "", err
	}

	account_type = account_data_map["user_type"].(string)
	return account_type, nil
}

func (auth *AuthServer) loadAccountsServiceData() (err error) {
	echo.Echo(echo.CyanFG, "Checking if Accounts service is up")
	authorization_service, err := auth.requestServiceData("accounts")
	if err == nil {
		auth.services["accounts"] = authorization_service // accounts_serv
	}
	return err
}

func (auth *AuthServer) run() {
	echo.Echo(echo.GreenFG, fmt.Sprintf("Starting %s server on -> %s:%s", auth.internal_dns, auth.host, auth.port))

	auth.router.RegisterRoute(prouter.NewRoute("/login", true), auth.handleLogin)

	auth.router.SetCorsHandler(prouter.CorsAllowAll)

	auth.boot()

	echo.Echo(echo.BlueFG, "Auth Server listening on - ", auth.host+":"+auth.port)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", auth.host, auth.port), auth.router); err != nil {
		echo.Echo(echo.RedFG, "Error starting auth server: ", err)
		os.Exit(1)
	}
}

func (auth *AuthServer) requestServiceData(service_name string) (*vlibs.Service, error) {
	var service *vlibs.Service = new(vlibs.Service)
	var client *http.Client = new(http.Client)

	request, _ := http.NewRequest("GET", fmt.Sprintf("%s/services/%s", vlibs.JD_ADDRESS, service_name), nil)
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	} else if response.StatusCode >= 300 {
		return nil, fmt.Errorf("JD responded with a %d for service %s", response.StatusCode, service_name)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return nil, err
	}

	var multipart_data []byte
	multipart_data, err = ioutil.ReadAll(response.Body)

	service, err = vlibs.PacientFromMultipart(multipart_data, response.Header.Get("Content-Type"))

	return service, err
}

func (auth *AuthServer) requestAccountData(username string) ([]byte, error) {
	var client *http.Client = new(http.Client)

	if _, exist := auth.services["accounts"]; !exist {
		if err := auth.loadAccountsServiceData(); err != nil {
			return nil, err
		}
	}

	var accounts_service *vlibs.Service = auth.services["accounts"]

	request, _ := http.NewRequest("GET", fmt.Sprintf("http://%s:%s/get-user?username=%s", accounts_service.Host, accounts_service.Port, username), nil)
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	} else if response.StatusCode >= 300 {
		return nil, fmt.Errorf("Accounts service responded with a %d for user %s", response.StatusCode, username)
	}

	defer response.Body.Close()
	var response_body []byte
	response_body, err = ioutil.ReadAll(response.Body)

	return response_body, err
}

func createAuthServer() *AuthServer {
	var auth_server *AuthServer = new(AuthServer)
	auth_server.router = prouter.CreateRouter()
	auth_server.host = "localhost"
	if hostname := os.Getenv("AUTH_HOST"); hostname != "" {
		auth_server.host = hostname
	}

	auth_server.port = "3030"
	if port := os.Getenv("AUTH_PORT"); port != "" {
		auth_server.port = port
	}

	auth_server.internal_dns = "authorization"
	if internal_dns := os.Getenv("AUTH_DNS"); internal_dns != "" {
		auth_server.internal_dns = internal_dns
	}

	auth_server.services = make(map[string]*vlibs.Service)

	return auth_server
}

func main() {
	var auth_server *AuthServer = createAuthServer()
	auth_server.run()
}
