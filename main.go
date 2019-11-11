package main

import (
	Mysql "learn-golang/Databases"
	"learn-golang/Router"
)

func main() {
	defer Mysql.DB.Close()
	Router.InitRouter()

}
