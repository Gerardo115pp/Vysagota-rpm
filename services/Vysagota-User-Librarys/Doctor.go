package vysagota_libs

import "encoding/json"

type Doctor struct {
	Uuid             int64
	Name             string
	Password         string
	Email            string
	RegistrationDate string
}

func (doctor *Doctor) ToJson() ([]byte, error) {
	var json_formatted []byte
	var err error
	json_formatted, err = json.Marshal(doctor)
	return json_formatted, err
}
