package fmsdata

import (
	"fmt"
)

func (d *DataBase) Post(layout string, data interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s/fmi/data/v1/databases/%s/layouts/%s/records", d.url, d.database, layout)
	return d.makeCall(url, "POST", data)
}

func (d *DataBase) PostWithScript(layout string, data interface{}, script string) ([]byte, error) {
	url := fmt.Sprintf("%s/fmi/data/v1/databases/%s/layouts/%s/records?script=%s", d.url, d.database, layout, script)
	return d.makeCall(url, "POST", data)
}
