package fmsdata

import (
	"fmt"
)

func (d *DataBase) Patch(data interface{}, layout, id string) ([]byte, error) {
	url := fmt.Sprintf("%s/fmi/data/v1/databases/%s/layouts/%s/records/%s", d.url, d.database, layout, id)
	return d.makeCall(url, "PATCH", data)
}
