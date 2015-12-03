package helpers

import (
	"github.com/carbonsrv/carbon/ctest"
	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestString(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.GET("/string", func(c *gin.Context) {
		String(c, 200, "Hello world!")
	})
	Convey("Given the simple route /string", t, func() {
		Convey("When a request hits", func() {
			w := ctest.Request(r, "GET", "/string")
			Convey("The Response Code should be 200", func() {
				So(w.Code, ShouldEqual, 200)
			})
			Convey("The Body should be \"Hello world!\"", func() {
				So(w.Body.String(), ShouldEqual, "Hello world!")
			})
			Convey("The Content Type should be text/plain", func() {
				So(w.HeaderMap.Get("Content-Type"), ShouldEqual, "text/plain")
			})
		})
	})
}

func TestHTMLString(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.GET("/htmlstring", func(c *gin.Context) {
		HTMLString(c, 200, "<p1>Hello world!</p1>")
	})
	Convey("Given the simple route /htmlstring", t, func() {
		Convey("When a request hits", func() {
			w := ctest.Request(r, "GET", "/htmlstring")
			Convey("The Response Code should be 200", func() {
				So(w.Code, ShouldEqual, 200)
			})
			Convey("The Body should be \"<p1>Hello world!</p1>\"", func() {
				So(w.Body.String(), ShouldEqual, "<p1>Hello world!</p1>")
			})
			Convey("The Content Type should be text/html", func() {
				So(w.HeaderMap.Get("Content-Type"), ShouldEqual, "text/html")
			})
		})
	})
}
