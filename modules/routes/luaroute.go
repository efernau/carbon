package routes

import (
	"../glue"
	"../scheduler"
	"bufio"
	"errors"
	"fmt"
	"github.com/DeedleFake/Go-PhysicsFS/physfs"
	"github.com/gin-gonic/gin"
	"github.com/pmylund/go-cache"
	"github.com/vifino/golua/lua"
	"github.com/vifino/luar"
	"net/http"
	"time"
)

// Cache
var cbc *cache.Cache

func cacheDump(L *lua.State, file string) (string, error, bool) {
	data_tmp, found := cbc.Get(file)
	if found == false {
		data, err := fileRead(file)
		if err != nil {
			return "", err, false
		}
		if L.LoadString(data) != 0 {
			return "", errors.New(L.ToString(-1)), true
		}
		res := L.FDump()
		L.Pop(1)
		cbc.Set(file, res, cache.DefaultExpiration)
		return res, nil, false
	} else {
		//debug("Using Bytecode-cache for " + file)
		return data_tmp.(string), nil, false
	}
}

// FS
var filesystem http.FileSystem

func fileRead(file string) (string, error) {
	f, err := filesystem.Open(file)
	defer f.Close()
	if err != nil {
		return "", err
	}
	fi, err := f.Stat()
	if err != nil {
		return "", err
	}
	r := bufio.NewReader(f)
	buf := make([]byte, fi.Size())
	_, err = r.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf), err
}

// Preloader/Starter
var jobs *int
var preloaded chan *lua.State

func preloader() {
	preloaded = make(chan *lua.State, *jobs)
	for {
		state := luar.Init()
		err := state.DoString(glue.Glue())
		if err != nil {
			fmt.Println(err)
		}
		preloaded <- state
	}
}
func getInstance() *lua.State {
	return <-preloaded
}

// Init
func Init(j *int) {
	jobs = j
	filesystem = physfs.FileSystem()
	cbc = cache.New(5*time.Minute, 30*time.Second) // Initialize cache with 5 minute lifetime and purge every 30 seconds
	go preloader()                                 // Run the instance starter.
}

func Lua(dir string) func(*gin.Context) {
	LDumper := luar.Init()
	return func(context *gin.Context) {
		L := getInstance()
		defer scheduler.Add(func() {
			L.Close()
		})
		file := dir + context.Request.URL.Path
		luar.Register(L, "", luar.Map{
			"context": context,
			"req":     context.Request,
		})
		code, err, lerr := cacheDump(LDumper, file)
		if err != nil {
			if lerr == false {
				context.Next()
				return
			} else {
				context.HTMLString(http.StatusInternalServerError, `<html>
				<head><title>Syntax Error in `+context.Request.URL.Path+`</title>
				<body>
					<h1>Syntax Error in file `+context.Request.URL.Path+`</h1>
					<code>`+string(err.Error())+`</code>
				</body>
				</html>`)
				context.Next()
				return
			}
		}
		L.LoadBuffer(code, len(code), file) // This shouldn't error, was checked earlier.
		if L.Pcall(0, 0, 0) != 0 {          // != 0 means error in execution
			context.HTMLString(http.StatusInternalServerError, `<html>
			<head><title>Runtime Error in `+context.Request.URL.Path+`</title>
			<body>
				<h1>Runtime Error in file `+context.Request.URL.Path+`</h1>
				<code>`+L.ToString(-1)+`</code>
			</body>
			</html>`)
			context.Abort()
		}
		/*L.DoString("return CONTENT_TO_RETURN")
		v := luar.CopyTableToMap(L, nil, -1)
		m := v.(map[string]interface{})
		i := int(m["code"].(float64))
		if err != nil {
			i = http.StatusOK
		}*/
		//context.HTMLString(i, m["content"].(string))
	}
}
