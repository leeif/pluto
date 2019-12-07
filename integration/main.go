package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/leeif/pluto/integration/test"
)

func main() {
	test.Integration()
}
