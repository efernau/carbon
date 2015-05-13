package glue

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
	"os"
	"time"
	"io/ioutil"
	"path"
	"path/filepath"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindata_file_info struct {
	name string
	size int64
	mode os.FileMode
	modTime time.Time
}

func (fi bindata_file_info) Name() string {
	return fi.name
}
func (fi bindata_file_info) Size() int64 {
	return fi.size
}
func (fi bindata_file_info) Mode() os.FileMode {
	return fi.mode
}
func (fi bindata_file_info) ModTime() time.Time {
	return fi.modTime
}
func (fi bindata_file_info) IsDir() bool {
	return false
}
func (fi bindata_file_info) Sys() interface{} {
	return nil
}

var _mainglue_lua = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xdc\x57\x4d\x73\xdb\x36\x13\x3e\x5b\xbf\x02\x81\x5f\xda\xe0\x5b\x8a\xb1\x73\xe8\x4c\x65\xd3\x9d\x4e\xdb\x99\x1e\xda\x69\xa7\xf1\xa5\x95\x14\x15\xa6\x20\x89\x23\x0a\x54\x49\x50\xa3\x54\x95\x7f\x7b\x17\x5c\x80\x04\x69\x51\x71\x0f\xb9\xf4\x10\xc5\xd8\x8f\x07\xbb\x0f\x16\x8b\x65\x9a\xc5\x3c\x25\x2b\xb5\x49\x67\xa2\x88\xf9\x56\x44\x87\x31\xbd\xa7\xd3\x88\x5e\xa5\xea\x8e\x06\x63\xfa\x50\x2d\x96\xb8\xb8\xaa\x16\x7c\xb3\xbd\xa3\xc7\x01\xfa\x96\x79\x62\x5d\x17\xa5\x8c\x55\x92\x49\xc6\xfd\xc1\x45\x2e\x54\x99\x4b\xc2\xa8\xe7\x79\x37\xef\xf6\xd4\x1f\x2d\xb2\x7c\xc3\x15\xe3\xa3\xa7\x8f\x4a\x30\xdf\x1f\x08\x39\x77\x40\x4a\x79\x06\xa6\x50\x79\x22\x97\x61\xbc\xe2\x39\x53\x99\x2c\x37\x4f\x22\x67\x3c\xb8\xfd\xd2\xc0\xd8\xe0\x07\x17\x3a\x97\x06\x02\xfc\x00\xa4\x0e\x06\x96\xa3\x65\x51\x3e\x31\x3a\xbe\x7f\xb8\x9a\xd2\xc0\xc9\x1c\x90\x2e\x00\x2a\x18\x5c\x94\xf9\xab\x10\x3e\xf0\xe1\x5f\xdf\x0c\x7f\xbf\x19\x7e\x35\x0b\x9f\x87\x00\xd6\x50\xd1\x60\x15\x2b\x91\xbe\x0a\xcd\x2b\xfe\x78\x7e\x73\xf9\xbf\xab\xff\x33\xff\xef\xc9\xe4\x7a\x42\xef\xee\x1f\xbe\x3e\x1c\xc7\xde\xf4\x03\x40\xd3\xc9\xc4\xbb\xa5\x06\x76\x70\x1c\xd4\x64\x1d\x5e\x1b\xee\x17\x00\x42\xe0\x14\x70\xe5\x79\xcc\xdb\x7b\x7b\x1f\x83\xb6\x68\x0d\xfe\x70\x48\x14\x5f\x92\x8d\x50\x5c\xf1\xa7\x54\x98\x73\x02\x19\x88\x56\x7a\xd7\x1c\x0c\x45\xee\x6c\x2c\xd2\x85\xde\x19\x0d\xe3\x4c\x2a\x21\x15\xac\x93\x05\xd1\xaa\xd0\x48\x88\x5a\x09\x09\xe2\x0b\xb3\x8e\x0e\x47\xbd\x82\xda\x20\x49\xb0\x23\x89\x24\xc9\x96\x27\x79\xc1\x5c\x27\x9f\xcc\x33\x6d\xa5\xc1\xd4\xc7\xad\x60\x3b\x3f\x8a\x28\xd6\x04\xad\x11\x6b\xcc\x71\x32\x8d\x76\x95\x44\xa4\x85\x68\xf9\x60\xe5\xf4\xf8\xa8\x0c\x21\xc1\xb4\xf6\x7e\x09\x3c\xc2\xcc\x99\xb1\x01\xba\xea\xff\xf0\x17\x19\xc8\xb6\x9a\x96\xc2\x61\xc0\x48\x08\x97\x73\x22\xc5\x5e\x31\x57\xea\xd7\x21\x19\x81\x43\xcc\x1a\x89\x71\x78\xa9\x9d\x3a\xbc\xac\x4f\xe7\x58\x1d\x61\x98\xc8\x42\xe4\x8a\x19\xdf\xc0\xc9\xb6\x9b\xee\x49\xfb\x75\x18\xd2\x68\x42\x69\x18\x3a\x9e\xa6\x9c\x40\x1c\xd0\xab\x3f\xcb\x0c\xda\x84\x0f\x76\xb0\xee\xa5\xc7\xd2\xb1\x88\xd3\xac\x10\x75\x94\x20\xee\x96\x88\xad\x62\x0c\x07\xb4\x31\xb4\x0f\x5b\x12\xb0\xcb\xfd\x5b\x08\xa6\xc2\x92\x7c\x23\x40\xf0\x40\x9b\xcd\xac\x73\x8f\x95\xa9\x0c\x8c\x24\xdb\x0a\x59\xef\x8a\xc7\x97\x8b\xa2\x4c\x95\x09\xcc\x9e\x9c\x13\x98\xd6\x46\xd0\x29\xdb\xd0\x44\xb3\xe3\x46\x6b\xc9\xd3\x57\xaf\x09\xd0\x12\xdd\x07\xd3\xca\xe3\x0c\x31\xe8\xdf\xd9\xd2\x12\xf4\x92\x8a\x3a\x27\x1b\xc0\x67\x4b\xf5\xd3\x59\xfe\xbb\x04\x81\x91\x9e\x2c\x5f\x5b\x06\x35\x12\x79\x8b\xe7\xaf\xb5\xd8\xa2\xf9\x7c\x3e\xb3\xdd\xa8\xd5\xce\x82\x30\x0c\x7d\xac\x58\x99\xa9\xd3\x6d\xcc\x15\xe2\x95\xc5\x7d\xab\x86\x16\xdd\x06\xa0\x17\xb1\x62\xd7\x97\xd7\x15\x9a\xb9\xb0\x48\xfc\x8e\xa7\xa5\x88\x8c\x45\x62\x77\x6b\x9a\x9c\x56\x9f\x6c\x74\x60\x71\x59\x69\x9f\xa3\x9b\xbe\x6b\xee\x06\x16\x60\x7b\x0f\xf5\x6b\x67\x60\xdb\xd7\xd3\x1e\x55\x3f\x02\x7a\x75\xef\xb2\x7d\x9d\xc1\xd2\xb2\x09\x89\xc7\xe2\xbf\xce\x67\xc3\xc6\xe7\xa0\x50\x17\xa4\x7d\x05\xda\x04\xaa\xa7\xb4\x4b\x60\xf7\xc2\xba\xc2\x2e\x81\xad\x87\x44\x63\x19\xfa\xce\x3f\x1f\x2f\x33\xb2\x57\x7d\xd7\xbe\xe8\xae\x6e\xbc\x36\x6f\xf0\x27\x72\x8d\x53\xc1\xf3\xd3\xe5\xa2\xc1\x5b\xe5\x20\x93\xf4\x1c\xc8\x49\xca\x6a\x10\xab\xed\x05\x29\x84\x3a\x5f\xb6\xfa\xef\x51\x2b\x5e\xe6\xb7\xa1\x46\x4e\x2b\x61\xe8\xd5\x60\x9f\x3d\x51\x07\xdb\xd8\x9d\xc4\xb6\x3a\xf4\x72\x2f\x9c\x7e\xc2\x7a\x32\xaf\x9e\xb7\x48\xe5\xa5\x38\x7b\x63\xf5\x6b\xdc\x87\x80\xca\x36\x04\xae\x70\x5c\x6c\xa6\x43\xa5\x67\xc3\xd9\x2c\x81\x11\x69\xdf\x49\x75\x6d\x4a\xb7\x29\xb5\xaa\xb2\x9a\x4a\xeb\xcb\xb6\x72\xb4\x55\x56\xcf\x04\xd5\x24\x0a\x65\xd6\x6e\xe6\xb3\x19\x44\x92\x46\x46\x1d\x3a\x07\x82\x53\xed\xfb\xad\x88\x93\x45\x12\x93\x47\xbe\x2c\x06\x36\x40\x0d\xc7\xf4\xeb\xe1\x7c\x71\x08\x55\x4f\xbf\xec\xa0\x75\x91\xfe\x39\x06\x55\x96\xf8\xd5\x51\xbb\xa7\x89\x5c\x33\x98\xc3\x03\xfd\x7e\x36\x18\x60\x4a\x39\x1d\x1f\x56\xb9\x58\x44\x95\xba\x94\x5b\x1e\xaf\xf5\x5b\xe9\x1f\xa7\x6d\x8c\x22\xce\x93\xad\x7e\xd8\xe6\xa2\x8d\x80\x0a\x6a\x34\x2d\x9f\x79\x16\x57\x6c\xb6\x1d\xde\x7c\xf7\xf3\xb7\x8f\xbf\xfd\xf2\x3d\xec\x4c\x75\xcb\xa7\xc7\xe9\xa8\xa9\x12\x86\x18\xc0\xc5\xaf\xe8\x62\xd1\x1a\x58\x5b\xc2\x73\x48\x3f\x20\x26\xa2\xea\x99\x9e\x0b\x12\x45\x04\xee\x90\x3d\x34\x14\x91\x77\x37\x37\x58\x0b\xf5\x01\x6b\x5f\x5f\x1b\x77\xdb\x6d\x05\xbe\x57\xe1\x7b\x9c\x21\x35\x40\x40\x2a\xeb\x41\x6b\x5e\x6f\x00\x4c\x99\xe8\xb9\x59\x0b\x43\x1c\xc0\xc9\x73\x27\x10\x84\xfd\xe1\xf1\xa7\x1f\x5f\x40\xd7\x33\xbb\xd9\xe3\x9c\x7d\x3d\xdd\x56\x01\x98\x2f\x22\xfd\xef\x9f\x00\x00\x00\xff\xff\xe8\x2a\xf1\x30\x2c\x0f\x00\x00")

func mainglue_lua_bytes() ([]byte, error) {
	return bindata_read(
		_mainglue_lua,
		"MainGlue.lua",
	)
}

func mainglue_lua() (*asset, error) {
	bytes, err := mainglue_lua_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "MainGlue.lua", size: 3884, mode: os.FileMode(420), modTime: time.Unix(1431533815, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

var _routeglue_lua = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x84\x51\xcd\x4e\xe3\x40\x0c\x3e\x27\x4f\x61\xf5\x94\x48\x6d\x55\xed\x71\xa5\xde\x76\x6f\xc0\xa1\xc0\x19\xa5\x89\x43\x46\x9a\xce\x44\x1e\x87\xd0\x0b\xcf\x8e\x3d\xf9\x69\x28\x45\x5c\xac\x64\xfc\xfd\xd8\x9f\x37\x1b\xb8\x37\x55\x65\xb1\x2f\x08\xa1\x41\xdb\x22\x85\x35\x38\xcf\x40\x58\x58\x7b\x86\x2e\x60\xdd\x59\xe8\x0d\x37\x50\xb8\x33\x37\xc6\xbd\xc2\xb1\x63\xe0\x06\x21\x20\xbd\x21\x81\x71\x86\x21\x94\x64\x5a\x4e\xeb\xce\x95\x6c\xbc\x83\x53\xbf\x75\xd8\x67\xb5\x13\x39\xec\x03\x17\x8c\x79\x9a\x58\x5f\x16\x16\x4a\x5f\x21\xec\x61\xb5\x4a\x13\x53\x03\x9f\x5b\x14\x5c\x0e\x7b\x79\x9a\xf8\x2b\x35\x70\x69\x92\x8c\xd8\xc0\x24\xce\xdb\xaa\x3b\xb5\x8a\x4d\x13\xb4\x01\xaf\xc9\x03\x68\xa6\xaa\x37\x12\x09\xdb\xfa\xa2\x1a\x9a\x99\xea\x09\x5d\x8d\x75\x4d\xed\x8f\xf0\x1f\xad\xa2\x97\x02\x04\xec\x29\x93\x1a\x1f\x5d\x95\x0e\x65\x58\x8a\xa6\x0f\xe9\xc7\xbd\xa6\xb5\x27\x7d\x9a\xa6\x91\x68\xfe\xdd\x1d\x5e\x1e\x1e\xa7\x61\x46\xfd\x6b\xc0\xe1\xf9\x82\x50\x1f\xd1\xd4\xfe\xc7\x1e\x9c\xb1\x93\xea\x72\xa6\x88\x22\xe4\x8e\x9c\x8c\xa3\x7f\xcb\x73\x60\xd9\xf8\xa8\xb7\x96\xe3\x86\x76\xbe\x86\xfe\x88\x29\x7b\xd7\x9d\x8e\x48\x59\x6c\x82\x27\xf8\xb3\xdb\x5d\x0e\x14\x07\xb9\x95\xf2\xe8\x27\x06\xff\xd5\x40\xd9\x6b\x58\x2c\xf6\x5d\x80\x8b\xa3\xc5\xdf\xf9\x7f\x49\x36\x90\x79\xf2\x45\x42\x33\x56\x97\x79\xc2\x77\x1e\xf1\xec\x97\xe7\x1d\xa3\xb8\x15\x40\xe4\xb0\x94\x39\x84\xaf\xfe\xb1\x9f\xdd\xcc\x22\x17\x1b\xe9\xe6\x51\xf7\x33\x00\x00\xff\xff\xf8\x3a\xf5\xa9\x3a\x03\x00\x00")

func routeglue_lua_bytes() ([]byte, error) {
	return bindata_read(
		_routeglue_lua,
		"RouteGlue.lua",
	)
}

func routeglue_lua() (*asset, error) {
	bytes, err := routeglue_lua_bytes()
	if err != nil {
		return nil, err
	}

	info := bindata_file_info{name: "RouteGlue.lua", size: 826, mode: os.FileMode(420), modTime: time.Unix(1431557323, 0)}
	a := &asset{bytes: bytes, info:  info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if (err != nil) {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"MainGlue.lua": mainglue_lua,
	"RouteGlue.lua": routeglue_lua,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func func() (*asset, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"MainGlue.lua": &_bintree_t{mainglue_lua, map[string]*_bintree_t{
	}},
	"RouteGlue.lua": &_bintree_t{routeglue_lua, map[string]*_bintree_t{
	}},
}}

// Restore an asset under the given directory
func RestoreAsset(dir, name string) error {
        data, err := Asset(name)
        if err != nil {
                return err
        }
        info, err := AssetInfo(name)
        if err != nil {
                return err
        }
        err = os.MkdirAll(_filePath(dir, path.Dir(name)), os.FileMode(0755))
        if err != nil {
                return err
        }
        err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
        if err != nil {
                return err
        }
        err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
        if err != nil {
                return err
        }
        return nil
}

// Restore assets under the given directory recursively
func RestoreAssets(dir, name string) error {
        children, err := AssetDir(name)
        if err != nil { // File
                return RestoreAsset(dir, name)
        } else { // Dir
                for _, child := range children {
                        err = RestoreAssets(dir, path.Join(name, child))
                        if err != nil {
                                return err
                        }
                }
        }
        return nil
}

func _filePath(dir, name string) string {
        cannonicalName := strings.Replace(name, "\\", "/", -1)
        return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

