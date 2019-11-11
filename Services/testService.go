package Services

import (
	"learn-golang/Models"
)

type Test struct {
	Id int `json:"id"`
	Testcol string `json:"testcol"`
}

func (this *Test) Insert() (id int, err error) {
	//var b = Controllers.NewSSEHandler()
	//b.SendString("hello")
	var testModel Models.Test
	testModel.Id = this.Id
	testModel.Testcol = this.Testcol
	id, err = testModel.Insert()
	return
}
