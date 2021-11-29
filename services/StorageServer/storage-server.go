package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	prouter "github.com/Gerardo115pp/PatriotRouter"
	_ "github.com/go-sql-driver/mysql"
	echo "github.com/gwen/putils/echo"
	vlibs "github.com/vysagota/libs"
)

type MysqlCredentials struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type StorageServer struct {
	router       *prouter.Router
	host         string
	port         string
	internal_dns string
	services     map[string]*vlibs.Service
	databases    map[string]MysqlCredentials
}

func (storage *StorageServer) boot() {
	/*

	   Registers itself and its actions to JD, creates a service objects and all the necessary
	   Actions to be used by the other services.

	*/
	echo.Echo(echo.GreenFG, "Booting Storage Server")

	var service *vlibs.Service = new(vlibs.Service)

	service.Name = "storage"
	service.Host = storage.host
	service.Port = storage.port
	service.DNS = storage.internal_dns
	service.Actions = make(map[string]*vlibs.Action)

	service.Actions["pacients/add-new"] = &vlibs.Action{
		Name:    "new-pacient",
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
		echo.EchoFatal(err)
	}

	echo.Echo(echo.GreenFG, "JD responded with a 2xx, boot process completed")
}

func (storage *StorageServer) createDSN(database string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", storage.databases[database].User, storage.databases[database].Password, storage.databases[database].Host, storage.databases[database].Database)
}

func (storage *StorageServer) handlePacients(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "POST":
		echo.Echo(echo.BlueFG, "Storage request: ", request.URL.Path, " POST")
		conn, err := sql.Open("mysql", storage.createDSN("rpm"))
		if err != nil {
			echo.Echo(echo.RedFG, "Cant connect to RPM database")
			echo.EchoFatal(err)
		}

		defer conn.Close()
		query := "INSERT INTO pacients (`name`, `username`, `password`, `email`, `phone`, `birthday`, `appointments_count`, `gender`, `tutor_email`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
		stmt, err := conn.Prepare(query)
		if err != nil {
			echo.Echo(echo.RedFG, "Cant prepare query")
			echo.EchoFatal(err)
		}

		defer stmt.Close()

		var name, username, password, email, phone, birthday, appointments_count, genes string

		err = request.ParseMultipartForm(32 << 20)
		if err != nil {
			echo.Echo(echo.RedFG, "Cant parse form")
			echo.EchoFatal(err)
		}
		// TODO: change vars for a pacient object
		name = request.FormValue("name")
		username = request.FormValue("username")
		password = request.FormValue("password")
		email = request.FormValue("email")
		phone = request.FormValue("phone")
		birthday = request.FormValue("birthday")
		appointments_count = "0"
		genes = request.FormValue("gender")

		_, err = stmt.Exec(name, username, password, email, phone, birthday, appointments_count, genes, request.FormValue("tutor_email"))

		if err != nil {
			echo.Echo(echo.RedFG, "Cant execute query")
			echo.EchoErr(err)
			response.WriteHeader(500)
		} else {
			response.WriteHeader(200)
		}
	case "GET":
		echo.Echo(echo.BlueFG, "Pacients ", echo.PinkFG, "GET", echo.BlueFG, " request: ", request.URL.Path)
		var pacients []vlibs.Pacient
		pacients, err := storage.queryPacients("1") // gets all pacients where 1(that means all)
		if err != nil {
			echo.Echo(echo.RedFG, "Cant query pacients")
			echo.EchoErr(err)
			response.WriteHeader(500)
		}

		response.Header().Set("Content-Type", "application/json")

		echo.Echo(echo.GreenFG, "Retrived pacients: ", len(pacients))

		response.WriteHeader(200)
		response.Header().Set("Content-Type", "application/json")

		json.NewEncoder(response).Encode(pacients)
	default:
		echo.Echo(echo.RedFG, fmt.Sprintf("Storage request: %s %s not supported", request.URL.Path, request.Method))
		response.WriteHeader(405) // Method not allowed
	}
}

func (storage *StorageServer) handleAccounts(response http.ResponseWriter, request *http.Request) {
	// Like pacients, but for both patients and doctors and only for GET

	switch request.Method {
	case "GET":
		echo.Echo(echo.BlueFG, "Accounts ", echo.PinkFG, "GET", echo.BlueFG, " request: ", request.URL.Path)

		var pacients []vlibs.Pacient
		pacients, err := storage.queryPacients("1") // gets all pacients where 1(that means all)

		if err != nil {
			echo.Echo(echo.RedFG, "Cant query pacients")
			echo.EchoErr(err)
			response.WriteHeader(500)
		}

		var doctors []vlibs.Doctor
		doctors, err = storage.queryDoctors("1") // gets all doctors where 1(that means all)
		if err != nil {
			echo.Echo(echo.RedFG, "Cant query doctors")
			echo.EchoErr(err)
			response.WriteHeader(500)
		}

		pacients_bytes, err := json.Marshal(pacients)
		if err != nil {
			echo.Echo(echo.RedFG, "Cant marshal pacients")
			echo.EchoErr(err)
			response.WriteHeader(500)
		}

		doctors_bytes, err := json.Marshal(doctors)
		if err != nil {
			echo.Echo(echo.RedFG, "Cant marshal doctors")
			echo.EchoErr(err)
			response.WriteHeader(500)
		}

		response.Header().Set("Content-Type", "application/json")
		var response_data string = fmt.Sprintf("{\"pacients\": %s, \"doctors\": %s}", string(pacients_bytes), string(doctors_bytes))

		response.WriteHeader(200)
		response.Write([]byte(response_data))

	}

	response.WriteHeader(200)
}

func (storage *StorageServer) run() {
	storage.router = prouter.CreateRouter()

	storage.router.RegisterRoute(prouter.NewRoute("/pacients/.+", false), storage.handlePacients)
	storage.router.RegisterRoute(prouter.NewRoute("/accounts", true), storage.handleAccounts)

	storage.router.SetCorsHandler(prouter.CorsAllowAll)

	storage.boot()

	echo.Echo(echo.BlueFG, "Storage Server listening on - ", storage.host+":"+storage.port)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", storage.host, storage.port), storage.router); err != nil {
		panic(err)
	}
}

func (storage *StorageServer) queryPacients(where string) ([]vlibs.Pacient, error) {

	conn, err := sql.Open("mysql", storage.createDSN("rpm"))
	if err != nil {
		echo.EchoErr(err)
		return nil, err
	}

	if where == "" {
		where = "1"
	} else if strings.Contains(where, ";") {
		// super secure security jaja
		echo.Echo(echo.RedFG, "Invalid query with where = ", where)
		return nil, fmt.Errorf("Invalid query")
	}

	defer conn.Close()
	query := fmt.Sprintf("SELECT * FROM pacients WHERE %s", where)

	rows, err := conn.Query(query)
	if err != nil {
		echo.Echo(echo.RedFG, "Cant execute query")
		echo.EchoErr(err)
	}

	defer rows.Close()
	var pacients []vlibs.Pacient = make([]vlibs.Pacient, 0)

	for rows.Next() {
		var pacient vlibs.Pacient
		var last_seen_nil_handler sql.NullTime

		if err := rows.Scan(&pacient.Uuid, &pacient.Name, &pacient.Username, &pacient.Password,
			&pacient.Email, &pacient.Phone, &pacient.TutorEmail, &pacient.Birthday, &pacient.AppointmentsCount,
			&last_seen_nil_handler, &pacient.CreatedAt, &pacient.Gender); err != nil {

			return nil, fmt.Errorf("Error while reading pacient %d: %s", len(pacients)+1, err.Error())
		}
		if last_seen_nil_handler.Valid {
			pacient.LastSeen = last_seen_nil_handler.Time
		}

		echo.EchoDebug(fmt.Sprintf("%#v\n", pacient))
		pacients = append(pacients, pacient)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return pacients, nil
}

func (storage *StorageServer) queryDoctors(where string) ([]vlibs.Doctor, error) {

	conn, err := sql.Open("mysql", storage.createDSN("rpm"))
	if err != nil {
		return nil, err
	}

	if where == "" {
		where = "1"
	} else if strings.Contains(where, ";") {
		// super secure security jaja
		echo.Echo(echo.RedFG, "Invalid query with where = ", where)
		return nil, fmt.Errorf("Invalid query")
	}

	defer conn.Close()
	query := fmt.Sprintf("SELECT * FROM doctors WHERE %s", where)

	rows, err := conn.Query(query)
	if err != nil {
		echo.Echo(echo.RedFG, "Cant execute query")
		echo.EchoErr(err)
	}

	defer rows.Close()
	var doctors []vlibs.Doctor = make([]vlibs.Doctor, 0)

	for rows.Next() {
		var doctor vlibs.Doctor

		if err := rows.Scan(&doctor.Uuid, &doctor.Name, &doctor.Password, &doctor.Email, &doctor.RegistrationDate); err != nil {
			return nil, fmt.Errorf("Error while reading doctor %d: %s", len(doctors)+1, err.Error())
		}
		echo.EchoDebug(fmt.Sprintf("%#v\n", doctor))
		doctors = append(doctors, doctor)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return doctors, nil
}

func createStorageServer() *StorageServer {
	storage := new(StorageServer)

	storage.router = prouter.CreateRouter()
	storage.host = os.Getenv("STORAGE_HOST")
	storage.port = os.Getenv("STORAGE_PORT")
	storage.internal_dns = os.Getenv("STORAGE_DNS")

	storage.services = make(map[string]*vlibs.Service)
	storage.databases = make(map[string]MysqlCredentials)
	storage.databases["rpm"] = MysqlCredentials{
		Host:     os.Getenv("RPM_HOST"),
		User:     os.Getenv("RPM_USER"),
		Password: os.Getenv("RPM_PASSWORD"),
		Database: os.Getenv("RPM_DATABASE"),
	}

	return storage
}

func main() {
	var storage *StorageServer = createStorageServer()
	storage.run()
}
