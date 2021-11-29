package vysagota_libs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime"
	"mime/multipart"
)

type Action struct {
	Name    string   `json:"name"`
	Methods []string `json:"methods"`
}

type Service struct {
	Name      string             `json:"name"`
	DNS       string             `json:"dns"`
	Host      string             `json:"host"`
	Port      string             `json:"port"`
	Protocols []string           `json:"protocols"`
	Status    string             `json:"status"`
	Actions   map[string]*Action `json:"actions"`
}

func (service *Service) ToJson() []byte {
	json, err := json.Marshal(service)
	if err != nil {
		return nil
	}
	return json
}

func (service *Service) ToMultipart() (*bytes.Buffer, string) {
	var buffer *bytes.Buffer = new(bytes.Buffer)
	var writer *multipart.Writer = multipart.NewWriter(buffer)
	writer.WriteField("name", service.Name)
	writer.WriteField("dns", service.DNS)
	writer.WriteField("port", service.Port)
	writer.WriteField("host", service.Host)
	writer.WriteField("status", service.Status)

	var protocols string = "["
	for h, protocol := range service.Protocols {
		protocols += "\"" + protocol + "\""
		if h < len(service.Protocols)-1 {
			protocols += ","
		}
	}
	protocols += "]"
	writer.WriteField("protocols", protocols)

	var actions []byte
	actions, err := json.Marshal(service.Actions)
	if err != nil {
		fmt.Println(err.Error())
		return nil, ""
	}
	writer.WriteField("actions", string(actions))
	writer.Close()

	return buffer, writer.FormDataContentType()
}

func PacientFromMultipart(data []byte, content_type string) (*Service, error) {
	/*
		reads a multipart form and returns a pacient
	*/
	_, params, _ := mime.ParseMediaType(content_type)
	var reader *multipart.Reader = multipart.NewReader(bytes.NewReader(data), params["boundary"])
	form, err := reader.ReadForm(32 << 20)
	if err != nil {
		return nil, err
	}
	var new_service *Service = new(Service)
	new_service.Name = form.Value["name"][0]
	new_service.DNS = form.Value["dns"][0]
	new_service.Port = form.Value["port"][0]
	new_service.Host = form.Value["host"][0]
	new_service.Status = form.Value["status"][0]

	err = json.Unmarshal([]byte(form.Value["protocols"][0]), &new_service.Protocols)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(form.Value["actions"][0]), &new_service.Actions)

	return new_service, nil
}
