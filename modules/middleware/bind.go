package middleware

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/DeedleFake/Go-PhysicsFS/physfs"
	"github.com/carbonsrv/carbon/modules/glue"
	"github.com/carbonsrv/carbon/modules/helpers"
	"github.com/carbonsrv/carbon/modules/scheduler"
	"github.com/carbonsrv/carbon/modules/static"
	"github.com/fzzy/radix/redis"
	"github.com/gin-gonic/gin"
	"github.com/pmylund/go-cache"
	"github.com/shurcooL/github_flavored_markdown"
	"github.com/vifino/contrib/gzip"
	"github.com/vifino/golua/lua"
	"github.com/vifino/luar"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

func Bind(L *lua.State) {
	BindCarbon(L)
	BindMiddleware(L)
	BindRedis(L)
	BindKVStore(L)
	BindPhysFS(L)
	BindIOEnhancements(L)
	BindOSEnhancements(L)
	BindThread(L)
	BindNet(L)
	BindConversions(L)
	BindComs(L)
	BindMarkdown(L)
	BindOther(L)
}

func BindCarbon(L *lua.State) {
	luar.Register(L, "carbon", luar.Map{ // Carbon specific API
		"glue": glue.GetGlue,
	})
}

func BindEngine(L *lua.State) {
	luar.Register(L, "carbon", luar.Map{
		"_gin_new": gin.New,
	})
}

func BindMiddleware(L *lua.State) {
	luar.Register(L, "mw", luar.Map{
		// Essentials
		"Logger":   gin.Logger,
		"Recovery": gin.Recovery,

		// Lua related stuff
		"Lua":       Lua,
		"DLR_NS":    DLR_NS,
		"DLR_RUS":   DLR_RUS,
		"DLRWS_RUS": DLRWS_RUS,

		// Custom sub-routers.
		"ExtRoute": (func(plan map[string]interface{}) func(*gin.Context) {
			newplan := make(Plan, len(plan))
			for k, v := range plan {
				newplan[k] = v.(func(*gin.Context))
			}
			return ExtRoute(newplan)
		}),
		"VHOST": (func(plan map[string]interface{}) func(*gin.Context) {
			newplan := make(Plan, len(plan))
			for k, v := range plan {
				newplan[k] = v.(func(*gin.Context))
			}
			return VHOST(newplan)
		}),
		"VHOST_Middleware": (func(plan map[string]interface{}) gin.HandlerFunc {
			newplan := make(Plan, len(plan))
			for k, v := range plan {
				newplan[k] = v.(gin.HandlerFunc)
			}
			return VHOST_Middleware(newplan)
		}),

		// To run or not to run, that is the question!
		"if_regex":       If_Regexp,
		"if_written":     If_Written,
		"if_status":      If_Status,
		"if_not_regex":   If_Not_Regexp,
		"if_not_written": If_Not_Written,
		"if_not_status":  If_Not_Status,

		// Modification stuff.
		"GZip": func() func(*gin.Context) {
			return gzip.Gzip(gzip.DefaultCompression)
		},

		// Basic
		"Echo":     EchoHTML,
		"EchoText": Echo,
	})
	luar.Register(L, "carbon", luar.Map{
		"_mw_CGI":         CGI,         // Run an CGI App!
		"_mw_CGI_Dynamic": CGI_Dynamic, // Run CGI Apps based on path!
		"_mw_combine": (func(middlewares []interface{}) func(*gin.Context) { // Combine routes, doesn't properly route like middleware or anything.
			newmiddlewares := make([]func(*gin.Context), len(middlewares))
			for k, v := range middlewares {
				newmiddlewares[k] = v.(func(*gin.Context))
			}
			return Combine(newmiddlewares)
		}),
	})
	L.DoString(glue.RouteGlue())
}

func BindPhysFS(L *lua.State) {
	luar.Register(L, "fs", luar.Map{ // PhysFS
		"mount":       physfs.Mount,
		"exits":       physfs.Exists,
		"getFS":       physfs.FileSystem,
		"mkdir":       physfs.Mkdir,
		"umount":      physfs.RemoveFromSearchPath,
		"delete":      physfs.Delete,
		"setWriteDir": physfs.SetWriteDir,
		"getWriteDir": physfs.GetWriteDir,
	})
}

func BindIOEnhancements(L *lua.State) {
	luar.Register(L, "carbon", luar.Map{ // Small enhancements to the io stuff.
		"_io_list": (func(path string) ([]string, error) {
			files, err := ioutil.ReadDir(path)
			if err != nil {
				return make([]string, 1), err
			} else {
				list := make([]string, len(files))
				for i := range files {
					list[i] = files[i].Name()
				}
				return list, nil
			}
		}),
		"_io_glob": filepath.Glob,
		"_io_modtime": (func(path string) (int, error) {
			info, err := os.Stat(path)
			if err != nil {
				return -1, err
			} else {
				return int(info.ModTime().UTC().Unix()), nil
			}
		}),
	})
}

func BindOSEnhancements(L *lua.State) {
	luar.Register(L, "carbon", luar.Map{ // Small enhancements to the io stuff.
		"_os_exists": (func(path string) bool {
			if _, err := os.Stat(path); err == nil {
				return true
			} else {
				return false
			}
		}),
		"_os_sleep": (func(secs int64) {
			time.Sleep(time.Duration(secs) * time.Second)
		}),
		"_os_chdir":   os.Chdir,
		"_os_abspath": filepath.Abs,
	})
}

func BindRedis(L *lua.State) {
	luar.Register(L, "redis", luar.Map{
		"connectTimeout": (func(host string, timeout int) (*redis.Client, error) {
			return redis.DialTimeout("tcp", host, time.Duration(timeout)*time.Second)
		}),
		"connect": (func(host string) (*redis.Client, error) {
			return redis.Dial("tcp", host)
		}),
	})
}

func BindKVStore(L *lua.State) { // Thread safe Key Value Store that doesn't persist.
	luar.Register(L, "kvstore", luar.Map{
		"_set": (func(k string, v interface{}) {
			kvstore.Set(k, v, -1)
		}),
		"_del": (func(k string) {
			kvstore.Delete(k)
		}),
		"_get": (func(k string) interface{} {
			res, found := kvstore.Get(k)
			if found {
				return res
			} else {
				return nil
			}
		}),
		"_inc": (func(k string, n int64) error {
			return kvstore.Increment(k, n)
		}),
		"_dec": (func(k string, n int64) error {
			return kvstore.Decrement(k, n)
		}),
	})
}

func BindThread(L *lua.State) {
	luar.Register(L, "thread", luar.Map{
		"_spawn": (func(bcode string, dobind bool, vals map[string]interface{}, buffer int) (chan interface{}, error) {
			var ch chan interface{}
			if buffer == -1 {
				ch = make(chan interface{})
			} else {
				ch = make(chan interface{}, buffer)
			}

			L := luar.Init()
			Bind(L)
			err := L.DoString(glue.MainGlue())
			if err != nil {
				panic(err)
			}

			luar.Register(L, "", luar.Map{
				"threadcom": ch,
			})

			if dobind {
				luar.Register(L, "", vals)
			}

			if L.LoadBuffer(bcode, len(bcode), "thread") != 0 {
				return make(chan interface{}), errors.New(L.ToString(-1))
			}

			scheduler.Add(func() {
				if L.Pcall(0, 0, 0) != 0 { // != 0 means error in execution
					fmt.Println("thread error: " + L.ToString(-1))
				}
			})
			return ch, nil
		}),
	})
}

func BindComs(L *lua.State) {
	luar.Register(L, "com", luar.Map{
		"create": (func() chan interface{} {
			return make(chan interface{})
		}),
		"createBuffered": (func(buffer int) chan interface{} {
			return make(chan interface{}, buffer)
		}),
		"receive": (func(c chan interface{}) interface{} {
			return <-c
		}),
		"try_receive": (func(c chan interface{}) interface{} {
			select {
			case msg := <-c:
				return msg
			default:
				return nil
			}
		}),
		"send": (func(c chan interface{}, val interface{}) bool {
			c <- val
			return true
		}),
		"try_send": (func(c chan interface{}, val interface{}) bool {
			select {
			case c <- val:
				return true
			default:
				return false
			}
		}),
		"size": (func(c chan interface{}) int {
			return len(c)
		}),
		"cap": (func(c chan interface{}) int {
			return cap(c)
		}),
		"pipe": (func(a, b chan interface{}) {
			for {
				b <- <-a
			}
		}),
		"pipe_background": (func(a, b chan interface{}) {
			scheduler.Add(func() {
				for {
					b <- <-a
				}
			})
		}),
	})
}

func BindNet(L *lua.State) {
	luar.Register(L, "net", luar.Map{
		"dial": net.Dial,
		"write": (func(con net.Conn, str string) {
			fmt.Fprintf(con, str)
		}),
		"readline": (func(con net.Conn) (string, error) {
			return bufio.NewReader(con).ReadString('\n')
		}),
		"pipe_com": (func(con net.Conn, input, output chan interface{}) {
			go func() {
				reader := bufio.NewReader(con)
				for {
					line, _ := reader.ReadString('\n')
					output <- line
				}
			}()
			for {
				line := <-input
				fmt.Fprintf(con, line.(string))
			}
		}),
		"pipe_com_background": (func(con net.Conn, input, output chan interface{}) {
			scheduler.Add(func() {
				reader := bufio.NewReader(con)
				for {
					line, _ := reader.ReadString('\n')
					output <- line
				}
			})
			scheduler.Add(func() {
				for {
					line := <-input
					fmt.Fprintf(con, line.(string))
				}
			})
		}),
	})
}

func BindConversions(L *lua.State) {
	luar.Register(L, "convert", luar.Map{
		"stringtocharslice": (func(x string) []byte {
			return []byte(x)
		}),
		"charslicetostring": (func(x []byte) string {
			return string(x)
		}),
	})
}

func BindContext(L *lua.State, context *gin.Context) {
	luar.Register(L, "", luar.Map{
		"context": context,
		"req":     context.Request,

		"host":   context.Request.URL.Host,
		"path":   context.Request.URL.Path,
		"scheme": context.Request.URL.Scheme,
	})
	luar.Register(L, "carbon", luar.Map{
		"_header_set": context.Header,
		"_header_get": context.Request.Header.Get,
		"_paramfunc":  context.Param,
		"_formfunc":   context.PostForm,
		"_queryfunc":  context.Query,
	})
}
func BindStatic(L *lua.State, cfe *cache.Cache) {
	luar.Register(L, "carbon", luar.Map{
		"_staticserve": (func(path, prefix string) func(*gin.Context) {
			return staticServe.ServeCached(prefix, staticServe.PhysFS(path, prefix, true, true), cfe)
		}),
	})
}

func BindMarkdown(L *lua.State) {
	luar.Register(L, "markdown", luar.Map{
		"github": (func(source string) string {
			return string(github_flavored_markdown.Markdown([]byte(source)))
		}),
	})
}

func BindOther(L *lua.State) {
	luar.Register(L, "", luar.Map{
		"unixtime": (func() int {
			return int(time.Now().UTC().Unix())
		}),
		"regexp": regexp.Compile,
	})
	luar.Register(L, "carbon", luar.Map{
		"_syntaxhl": helpers.SyntaxHL,
	})
}
