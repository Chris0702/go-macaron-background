package server

import (
	"encoding/json"
	"fmt"
	"gopkg.in/macaron.v1"
	"net/http"
	"os"
	// "reflect"
)
import db "../database"

func ServerRun() {
	fmt.Println("ServerRun")
	result, serverPort, host, port, user, password, dbname := cloudfoundryInit()
	if result == false {
		serverPort = "4000"
		host = "localhost"
		port = "5432"
		user = "postgres"
		password = "79317931"
		dbname = "postgres"
	}
	m := macaron.Classic()
	db.InitDB("postgres", host, port, user, password, dbname, result)
	db.CreateTable()
	ServerUse(m)
	RestApiOn(m)
	WebsiteOn(m)
	fmt.Println("server listening  " + serverPort)
	http.ListenAndServe(":"+serverPort, m)
}

func cloudfoundryInit() (result bool, serverPort string, host string, port string, user string, password string, dbname string) {
	result = false
	serverPort = os.Getenv("PORT")
	if serverPort == "" {
		return
	}
	VCAP_SERVICES := os.Getenv("VCAP_SERVICES")
	if VCAP_SERVICES != "" {
		u := map[string]interface{}{}
		err := json.Unmarshal([]byte(VCAP_SERVICES), &u)
		if err != nil {
			return
		}

		if postgresqlRec, ok := u["postgresql"].([]interface{}); ok {
			for _, val := range postgresqlRec {
				if postgresqlRecInner, ok := val.(map[string]interface{}); ok {
					for key, val := range postgresqlRecInner {
						// fmt.Printf(" [========>] %s = %s", key, val)
						if key == "credentials" {
							if credentials, ok := val.(map[string]interface{}); ok {
								// fmt.Println(" [======password=======%s",credentials["password"])
								// fmt.Println(" [======database=======%s",credentials["database"])
								// fmt.Println(" [======port=======%s",credentials["port"])
								// fmt.Println(" [======host=======%s",credentials["host"])
								// fmt.Println(" [======username=======%s",credentials["username"])
								// fmt.Println(" [======password=======%s",credentials["password"])
								host = fmt.Sprintf("%v", credentials["host"])
								port = fmt.Sprintf("%v", credentials["port"])
								user = fmt.Sprintf("%v", credentials["username"])
								password = fmt.Sprintf("%v", credentials["password"])
								dbname = fmt.Sprintf("%v", credentials["database"])
								result = true
								return
							}
						}
					}
				}
			}
			return
		} else {
			return
		}
	}
	return
}
