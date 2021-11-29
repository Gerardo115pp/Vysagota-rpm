package vysagota_libs

import (
	"bytes"
	"fmt"
	"mime/multipart"
)

type EventType int

const (
	EVE_USER_CREATED EventType = iota
	EVE_USER_UPDATED
	EVE_USER_DELETED
)

type DomainEvent struct {
	Name string
	Type EventType
	Data map[string]string
}

func (domain_event *DomainEvent) toJson() string {
	var json_string string = "{"
	json_string += fmt.Sprintf("\"name\":\"%s\",", domain_event.Name)
	json_string += fmt.Sprintf("\"type\":%d,", domain_event.Type)

	var data_string string = "\"data\":{"
	for key, value := range domain_event.Data {
		data_string += fmt.Sprintf("\"%s\":\"%s\",", key, value)
	}
	data_string = data_string[:len(data_string)-1] // remove last comma
	data_string += "}"
	json_string += data_string
	return json_string + "}"
}

func (domain_event *DomainEvent) toMultipart() []byte {
	var buffer *bytes.Buffer = &bytes.Buffer{}
	var writer *multipart.Writer = multipart.NewWriter(buffer)

	writer.WriteField("name", domain_event.Name)
	writer.WriteField("type", fmt.Sprintf("%d", domain_event.Type))
	for key, value := range domain_event.Data {
		writer.WriteField(key, value)
	}
	writer.Close()
	return buffer.Bytes()
}
