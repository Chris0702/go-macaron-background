package server

import "gopkg.in/macaron.v1"
import "fmt"

func WebsiteOn(m *macaron.Macaron){
	fmt.Println("WebsiteOn")
	m.Get("/", func(ctx *macaron.Context) {
		ctx.HTML(200, "index") // 200 为响应码
	})
}