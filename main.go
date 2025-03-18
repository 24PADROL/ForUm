package main

import (
    e "engine/server"
)

func main() {
    e.InitDB()
    var forum e.User
	e.Run(&forum)
}   