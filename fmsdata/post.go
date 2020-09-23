package fmsdata

import (
	"fmt"
)

func (d *DataBase) Post(layout string, data interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s/fmi/data/v1/databases/%s/layouts/%s/records", d.url, d.database, layout)
	return d.makeCall(url, "POST", data)
}
