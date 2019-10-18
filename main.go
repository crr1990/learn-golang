package main

import (
"learn-golang/Router"
"learn-golang/Databases"
)

func main() {
	defer Mysql.DB.Close()
	Router.InitRouter()
}
