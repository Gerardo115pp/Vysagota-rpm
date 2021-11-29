package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	prouter "github.com/Gerardo115pp/PatriotRouter"
	echo "github.com/gwen/putils/echo"
	vlibs "github.com/vysagota/libs"
)

type AccountsServer struct {
	router   *prouter.Router
	pacients map[int64]vlibs.Pacient
	doctors  map[int64]vlibs.Doctor
	host     string
	port     string
	dns      string
	services map[string]*vlibs.Service
}

func (accounts *AccountsServer) authenticate(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(200)
}

func (accounts *AccountsServer) boot() {
	/*

	   Registers itself and its actions to JD, creates a service objects and all the necessary
	   Actions to be used by the other services.

	*/
	echo.Echo(echo.GreenFG, "Booting Accounts Server")

	var service *vlibs.Service = new(vlibs.Service)

	service.Name = "accounts"
	service.Host = accounts.host
	service.Port = accounts.port
	service.DNS = "accounts"
	service.Actions = make(map[string]*vlibs.Action)

	service.Actions["authenticate"] = &vlibs.Action{
		Name:    "authenticate",
		Methods: []string{"POST"},
	}
	service.Actions["register"] = &vlibs.Action{
		Name:    "register",
		Methods: []string{"POST"},
	}
	service.Actions["account_exists"] = &vlibs.Action{
		Name:    "account-exists",
		Methods: []string{"GET"},
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
		echo.EchoFatal(err)
	}

	echo.Echo(echo.GreenFG, "JD responded with a 2xx")

	// check if Storage service is up, if so load the pacients and doctors, doctors will remain as a TODO

	echo.Echo(echo.CyanFG, "Checking if Storage service is up")
	storage_service, err := accounts.requestServiceData("storage")
	if err != nil {
		echo.Echo(echo.RedFG, "Storage service is not up, accounts data will have to be loaded later")
		echo.EchoErr(err)
	} else {
		accounts.services["storage"] = storage_service
		err = accounts.loadAccounts()
		if err != nil {
			echo.Echo(echo.RedFG, "Error loading accounts data")
			echo.EchoErr(err)
		} else {
			echo.Echo(echo.GreenFG, "Success: accounts data loaded")
		}
	}

	echo.Echo(echo.GreenFG, "Accounts Server is up")
}

func (accounts *AccountsServer) getPacients(response http.ResponseWriter, request *http.Request) {
	var http_method string = request.Method
	echo.Echo(echo.CyanFG, "%s request from %s", http_method, request.RemoteAddr)

	switch http_method {
	case "GET":
		request.ParseMultipartForm(32 << 20)
		var username string = request.FormValue("username")
		// TODO: also support email, uuid.
		if username == "" {
			// return everything
			var pacients_json []byte
			var err error
			echo.Echo(echo.GreenFG, fmt.Sprintf("Serving %d pacients", len(accounts.pacients)))

			var pacients_list []vlibs.Pacient
			for _, pacient := range accounts.pacients {
				pacients_list = append(pacients_list, pacient)
			}

			pacients_json, err = json.Marshal(pacients_list)
			if err != nil {
				echo.Echo(echo.RedFG, "Error marshalling pacients")
				echo.EchoErr(err)
				response.WriteHeader(500)
			} else {
				echo.Echo(echo.GreenFG, "Success marshalling pacients")
				response.Header().Set("Content-Type", "application/json")
				response.Write(pacients_json)
			}
		} else {

		}

	case "OPTIONS":
		response.Header().Set("Access-Control-Allow-Headers", "Content-Type,Cache-Control,X-Requested-With,Authorization, Pragma, Accept, Accept-Encoding, Accept-Language, Origin, User-Agent, Content-Length, X-CSRF-Token, X-Requested-With, X-XSRF-Token")
		response.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		response.WriteHeader(200)

	default:
		echo.Echo(echo.RedFG, "Unsupported HTTP method")
		response.WriteHeader(405)
	}
}

func (accounts *AccountsServer) getPacientByUsername(username string) (vlibs.Pacient, error) {
	var pacient vlibs.Pacient
	var err error
	for _, pacient = range accounts.pacients {
		if pacient.Username == username {
			break
		}
	}
	if pacient.Username != username {
		err = fmt.Errorf("Pacient %s not found", username)
	}
	return pacient, err
}

func (accounts *AccountsServer) getDoctorByUsername(username string) (vlibs.Doctor, error) {
	var doctor vlibs.Doctor
	var err error
	for _, doctor = range accounts.doctors {
		if doctor.Name == username {
			break
		} else {
			echo.Echo(echo.RedFG, "%s !== %s", doctor.Name, username)
		}
	}
	echo.Echo(echo.RedFG, fmt.Sprint(len(accounts.doctors)))
	if doctor.Name != username {
		err = fmt.Errorf("Doctor %s not found", username)
	}
	return doctor, err
}

func (accounts *AccountsServer) getUser(response http.ResponseWriter, request *http.Request) {
	// gets a username, uuid or email, verifies if the user exists, and if it does,
	// returns the user data + the usertype (doctor or pacient).
	// Only supports GET requests.
	// TODO: add support for email and uuid
	if request.Method == "GET" {
		var username string = request.URL.Query().Get("username")

		echo.Echo(echo.CyanFG, fmt.Sprintf("GET request for user %s", username))

		if username != "" {
			pacient_data, err := accounts.getPacientByUsername(username)
			if err == nil {
				echo.Echo(echo.GreenFG, "Serving pacient data")
				pacient_json, err := pacient_data.ToJson()
				if err != nil {
					echo.Echo(echo.RedFG, "Error marshalling pacient")
					response.WriteHeader(500)
					return
				}
				var response_data string = fmt.Sprintf("{\"user_data\": %s, \"user_type\": \"pacient\"}", string(pacient_json))

				response.Header().Set("Content-Type", "application/json")
				response.Write([]byte(response_data))
				return
			}

			doctor_data, err := accounts.getDoctorByUsername(username)
			echo.Echo(echo.GreenFG, fmt.Sprintf("Doctor: %v+", doctor_data))
			if err == nil {
				echo.Echo(echo.GreenFG, "Serving doctor data")
				doctor_json, err := doctor_data.ToJson()
				if err != nil {
					echo.Echo(echo.RedFG, "Error marshalling doctor")
					response.WriteHeader(500)
					return
				}
				var response_data string = fmt.Sprintf("{\"user_data\": %s, \"user_type\": \"doctor\"}", string(doctor_json))

				response.Header().Set("Content-Type", "application/json")
				response.Write([]byte(response_data))
				return
			} else {
				echo.Echo(echo.RedFG, "User not found")
				response.WriteHeader(404)
				return
			}

		} else {
			// TODO: support for email and uuid
			response.WriteHeader(http.StatusNotImplemented)
		}
		response.WriteHeader(200)
	} else {
		response.WriteHeader(405) // method not allowed
	}
}

func (accounts *AccountsServer) loadAccounts() (err error) {
	if _, exists := accounts.services["storage"]; !exists {
		return fmt.Errorf("Storage service is not up")
	}

	var storage_service *vlibs.Service = accounts.services["storage"]
	var client *http.Client = new(http.Client)

	echo.Echo(echo.CyanFG, "Requesting pacients data from Storage service")
	request, _ := http.NewRequest("GET", fmt.Sprintf("http://%s:%s/accounts", storage_service.Host, storage_service.Port), nil)
	response, err := client.Do(request)
	if err != nil {
		return err
	} else if response.StatusCode >= 300 {
		return fmt.Errorf("Storage service responded with %d", response.StatusCode)
	}

	defer response.Body.Close()

	data := &struct {
		Pacients []vlibs.Pacient `json:"pacients"`
		Doctors  []vlibs.Doctor  `json:"doctors"`
	}{}

	var body_data []byte
	body_data, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body_data, data)
	if err != nil {
		echo.EchoErr(err)
		return err
	}

	for _, pacient := range data.Pacients {
		accounts.pacients[pacient.Uuid] = pacient
	}
	echo.Echo(echo.GreenFG, fmt.Sprintf("Success: %d pacients data loaded", len(accounts.pacients)))

	for _, doctor := range data.Doctors {
		accounts.doctors[doctor.Uuid] = doctor
	}
	echo.Echo(echo.GreenFG, fmt.Sprintf("Success: %d doctors data loaded", len(accounts.doctors)))

	return nil
}

func (accounts AccountsServer) newPacientFromRequest(request *http.Request) (*vlibs.Pacient, error) {
	var new_pacient *vlibs.Pacient = new(vlibs.Pacient)

	new_pacient.Name = request.FormValue("name")
	new_pacient.Username = request.FormValue("username")
	new_pacient.Password = request.FormValue("password")
	new_pacient.Email = request.FormValue("email")
	new_pacient.TutorEmail = request.FormValue("tutor_email")
	new_pacient.Phone = request.FormValue("phone")
	new_pacient.AppointmentsCount = 0

	new_pacient.Gender = request.FormValue("gender")

	var datetime_iso_string string = request.FormValue("birthday")
	new_pacient.LastSeen = time.Now()
	var birthday time.Time
	birthday, err := time.Parse("2006-01-02T15:04:05.000Z", datetime_iso_string)

	if err != nil {
		return nil, err
	}

	new_pacient.Birthday = birthday

	return new_pacient, nil
}

func (accounts *AccountsServer) register(response http.ResponseWriter, request *http.Request) {
	var new_pacient *vlibs.Pacient
	var err error

	echo.Echo(echo.CyanFG, "Parsing new pacient data")
	new_pacient, err = accounts.newPacientFromRequest(request)
	if err != nil {
		echo.Echo(echo.RedFG, "Error parsing pacient data")
		echo.EchoErr(err)
		response.WriteHeader(400)
		return
	}

	echo.Echo(echo.CyanFG, fmt.Sprintf("Creating new pacient '%s'", new_pacient.Name))

	err = accounts.saveNewPacient(new_pacient)
	if err != nil {
		echo.EchoErr(err)
		response.WriteHeader(500)
		return
	}

	response.WriteHeader(200)
}

func (accounts *AccountsServer) requestServiceData(service_name string) (*vlibs.Service, error) {
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

func (accounts *AccountsServer) run() {
	accounts.router = prouter.CreateRouter()

	accounts.router.RegisterRoute(prouter.NewRoute("/authenticate", true), accounts.authenticate)
	accounts.router.RegisterRoute(prouter.NewRoute("/register", true), accounts.register)
	accounts.router.RegisterRoute(prouter.NewRoute("/pacients", true), accounts.getPacients)
	accounts.router.RegisterRoute(prouter.NewRoute("/get-user", true), accounts.getUser)

	accounts.router.SetCorsHandler(prouter.CorsAllowAll)

	accounts.boot()

	echo.Echo(echo.BlueFG, "Accounts Server listening on -", accounts.host+":"+accounts.port)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", accounts.host, accounts.port), accounts.router); err != nil {
		panic(err)
	}
}

func (accounts *AccountsServer) saveNewPacient(new_pacient *vlibs.Pacient) error {
	/*

	   Saves a new pacient to the database.

	*/
	// check if Storage service is on services list
	var storage_service *vlibs.Service

	if _, exists := accounts.services["storage"]; !exists {
		echo.Echo(echo.OrangeFG, "Storage service is not on services cache, requesting it")
		storage_service, err := accounts.requestServiceData("storage")
		if err != nil {
			echo.Echo(echo.RedFG, "Error requesting Storage service")
			echo.EchoErr(err)
			return err
		}
		echo.Echo(echo.GreenFG, fmt.Sprintf("Storage service found on '%s:%s'", storage_service.Host, storage_service.Port))
		accounts.services["storage"] = storage_service
	}

	storage_service = accounts.services["storage"]

	var request_url string = fmt.Sprintf("http://%s:%s/pacients/add", storage_service.Host, storage_service.Port)

	var client *http.Client = new(http.Client)
	var multipart_data *bytes.Buffer
	var content_type string
	multipart_data, content_type = new_pacient.ToMultipart()

	request, _ := http.NewRequest("POST", request_url, multipart_data)
	request.Header.Set("Content-Type", content_type)

	response, err := client.Do(request)
	if err != nil {
		echo.Echo(echo.RedFG, "Error connecting to Storage service")
		echo.EchoErr(err)
		return err
	} else if response.StatusCode >= 300 {
		err = fmt.Errorf("Storage service responded with a %d", response.StatusCode)
		return err
	}
	defer response.Body.Close()
	echo.Echo(echo.GreenFG, "Storage service responded with a 2xx, pacient saved")
	go accounts.loadAccounts()

	return nil
}

func createAccountsServer() *AccountsServer {
	var accounts *AccountsServer = new(AccountsServer)
	accounts.host = os.Getenv("ACC_HOST")
	accounts.port = os.Getenv("ACC_PORT")
	accounts.dns = os.Getenv("ACC_DNS")
	accounts.pacients = make(map[int64]vlibs.Pacient)
	accounts.services = make(map[string]*vlibs.Service)
	accounts.doctors = make(map[int64]vlibs.Doctor)
	// accounts.router will be set in run()
	return accounts
}

func main() {
	var accounts *AccountsServer = createAccountsServer()
	accounts.run()
}
