package server

import "gopkg.in/macaron.v1"
import "fmt"

func ServerUse(m *macaron.Macaron) {
	fmt.Println("ServerUse")
	m.Use(macaron.Renderer(macaron.RenderOptions{
		// 模板文件目录，默认为 "client"
		Directory: "client",
		// 模板文件后缀，默认为 [".tmpl", ".html"]
		Extensions: []string{".tmpl", ".html"},
		// 模板语法分隔符，默认为 ["{{", "}}"]
		Delims: macaron.Delims{"<%=", "%>"},
		//<%=contextUrl%>
		// 追加的 Content-Type 头信息，默认为 "UTF-8"
		Charset: "UTF-8",
		// 渲染具有缩进格式的 JSON，默认为不缩进
		IndentJSON: true,
		// 渲染具有缩进格式的 XML，默认为不缩进
		IndentXML: true,
		// 渲染具有前缀的 JSON，默认为无前缀
		PrefixJSON: []byte("macaron"),
		// 渲染具有前缀的 XML，默认为无前缀
		PrefixXML: []byte("macaron"),
		// 允许输出格式为 XHTML 而不是 HTML，默认为 "text/html"
		HTMLContentType: "text/html",
	}))
	m.Use(macaron.Static("client"))
	m.Use(macaron.Static("client/custom",
		macaron.StaticOptions{
			Prefix:      "custom",
			SkipLogging: true,
		}))
	m.Use(macaron.Static("client/custom/previews",
		macaron.StaticOptions{
			Prefix:      "/",
			SkipLogging: true,
		}))
}
