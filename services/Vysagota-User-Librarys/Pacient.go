package vysagota_libs

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"time"
)

type Pacient struct {
	Uuid              int64     `json:"uuid"`
	Name              string    `json:"name"`
	Birthday          time.Time `json:"birthday"`
	Username          string    `json:"username"`
	Password          string    `json:"secret"`
	Email             string    `json:"email"`
	TutorEmail        string    `json:"tutor_email"`
	Phone             string    `json:"phone"`
	AppointmentsCount int       `json:"appointments_count"`
	CreatedAt         time.Time `json:"created_at"`
	Gender            string    `json:"gender"`
	LastSeen          time.Time `json:"last_seen"`
}

func (pacient *Pacient) ToJson() ([]byte, error) {
	var json_formatted []byte
	var err error
	json_formatted, err = json.Marshal(pacient)
	return json_formatted, err
}

func (Pacient *Pacient) ToMultipart() (*bytes.Buffer, string) {
	var buffer *bytes.Buffer = new(bytes.Buffer)
	var writer *multipart.Writer = multipart.NewWriter(buffer)
	defer writer.Close()

	writer.WriteField("uuid", string(Pacient.Uuid))
	writer.WriteField("name", Pacient.Name)
	writer.WriteField("birthday", Pacient.Birthday.Format("2006-01-02"))
	writer.WriteField("username", Pacient.Username)
	writer.WriteField("password", Pacient.Password)
	writer.WriteField("email", Pacient.Email)
	writer.WriteField("tutor_email", Pacient.TutorEmail)
	writer.WriteField("phone", Pacient.Phone)

	var appointments_count_str string = string(Pacient.AppointmentsCount)
	writer.WriteField("appointments_count", appointments_count_str)

	writer.WriteField("last_seen", Pacient.LastSeen.Format("2006-01-02"))
	writer.WriteField("created_at", Pacient.CreatedAt.Format("2006-01-02"))
	writer.WriteField("gender", Pacient.Gender)

	return buffer, writer.FormDataContentType()

}
