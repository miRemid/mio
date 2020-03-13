package mio

import (
	"log"
	"net/http"
	"testing"
)

func TestMio(t *testing.T) {
	mil := Default()
	v1 := mil.NewGroup("/hello")
	{
		v2 := v1.NewGroup("/haha")
		{
			v2.GET("/", func(c *Context) {
				c.String(http.StatusOK, "Hello haha")
			})
		}
	}
	log.Fatal(mil.Server(":8800"))
}
