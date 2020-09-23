package fmsdata

import (
	"fmt"
)

func (d *DataBase) Script(layout, script, param string) ([]byte, error) {
	url := fmt.Sprintf("%s/fmi/data/v1/databases/%s/layouts/%s/records?script=%s&script.param=%s&_limit=1", d.url, d.database, layout, script, param)
	//fmt.Println(url)

	return d.makeCall(url, "GET", nil)
}
