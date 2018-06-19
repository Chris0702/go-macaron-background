package server

import "gopkg.in/macaron.v1"
import "fmt"

func RestApiOn(m *macaron.Macaron){
	fmt.Println("RestApiOn")
	m.Get("/explore/*", explore)

	m.Post("/upload", upload)

	m.Post("/mkdir", mkdir)

	m.Post("/rename", rename)

	m.Post("/remove", remove)

	m.Get("/displays/*", displays)

	m.Get("/symbols/*", symbols)

	m.Get("/components/*", components)

	m.Get("/assets/*", assets)

}
