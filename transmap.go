package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net"
	"os"
	"regexp"
)

func main() {
	config := make(map[string]string)

	conf, err := os.Open("transmap.cfg")
	defer conf.Close()
	if err != nil {
		panic(err)
	}
	cr, _ := regexp.Compile("(\\w+)\\s*=\\s*(.*\\w)")
	scan := bufio.NewScanner(conf)
	for scan.Scan() {
		cm := cr.FindStringSubmatch(scan.Text())
		config[cm[1]] = cm[2]
	}

	db, _ := sql.Open("mysql", config["connect_string"])

	res, err := db.Query("select id,name from transmap.test order by id asc;")
	defer res.Close()
	if err != nil {
		panic(err)
	}

	for res.Next() {
		var id int
		var name string
		res.Scan(&id, &name)
		fmt.Println(id, name)
	}

	fmt.Println("Started")

	ln, _ := net.Listen("tcp", ":12345")
	for {
		c, _ := ln.Accept()
		go handleConn(c)
	}
}

func handleConn(l net.Conn) {
	msg, _ := bufio.NewReader(l).ReadString('\n')

	fmt.Println(string(msg))
	l.Close()
}
