package fmsdata

import (
	"fmt"
)

func (d *DataBase) GetAllFrom(layout string) ([]byte, error) {
	url := fmt.Sprintf("%s/fmi/data/v1/databases/%s/layouts/%s/records?_limit=20000", d.url, d.database, layout)
	return d.makeCall(url, "GET", nil)
}

func (d *DataBase) GetAllFromWithMax(layout string, max int) ([]byte, error) {
	url := fmt.Sprintf("%s/fmi/data/v1/databases/%s/layouts/%s/records?_limit=%d", d.url, d.database, layout, max)
	return d.makeCall(url, "GET", nil)
}
