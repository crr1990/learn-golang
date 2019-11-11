package Models

import (
	"learn-golang/Databases"
)

type Soul struct {
	Id int
	Testcol string `gorm:"column:name"`
}

// 设置Test的表名为`test`
// func (Test) TableName() string {
//     return "test"
// }

func (this *Soul) Insert() (id int, err error) {
	result := Mysql.DB.Create(&this)
	id = this.Id
	if result.Error != nil {
		err = result.Error
		return
	}
	return
}