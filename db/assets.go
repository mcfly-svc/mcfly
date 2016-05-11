// Code generated by go-bindata.
// sources:
// db/migrations/0001_init_db.down.sql
// db/migrations/0001_init_db.up.sql
// DO NOT EDIT!

package db

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _dbMigrations0001_init_dbDownSql = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x72\x09\xf2\x0f\x50\xf0\xf4\x73\x71\x8d\x50\xf0\x74\x53\x70\x8d\xf0\x0c\x0e\x09\x56\x28\x28\xca\xcf\x4a\x4d\x2e\x89\x2f\x4e\x2d\x2a\xcb\x4c\x4e\x8d\x2f\x05\x32\xf2\x12\x73\x53\xe3\x41\x84\x35\x17\x2e\x3d\x65\x99\x29\xa9\x45\x89\xc9\xc9\xa9\xc5\xc5\x25\xf9\xd9\xa9\x79\x60\x7d\xd6\x5c\x10\xf5\x21\x8e\x4e\x3e\xae\x48\xea\x93\x4a\x33\x73\x52\xac\xb1\xcb\x81\xf4\xc5\x43\x1d\x81\x43\x09\x41\x59\xb0\x63\xe2\x21\xae\x89\x07\x3b\x07\x87\xda\xdc\xc4\xa2\xe2\xd2\x82\x4c\x54\xc7\x46\x06\x60\x33\xce\x9a\x0b\x10\x00\x00\xff\xff\xa2\x6a\x6e\x7d\x2f\x01\x00\x00")

func dbMigrations0001_init_dbDownSqlBytes() ([]byte, error) {
	return bindataRead(
		_dbMigrations0001_init_dbDownSql,
		"db/migrations/0001_init_db.down.sql",
	)
}

func dbMigrations0001_init_dbDownSql() (*asset, error) {
	bytes, err := dbMigrations0001_init_dbDownSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "db/migrations/0001_init_db.down.sql", size: 303, mode: os.FileMode(420), modTime: time.Unix(1462573876, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _dbMigrations0001_init_dbUpSql = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\x9c\x92\xcd\x6e\x82\x40\x14\x85\xd7\xf0\x14\x77\x07\x26\xbc\x41\x57\x14\x6f\x13\x53\xa5\x16\x35\xa9\x2b\x82\x30\xd1\x51\x1c\xcc\xfc\x18\x1f\xbf\xcc\xc8\x38\x58\x69\x62\xcb\x6a\x7e\x2e\xe7\x7c\xf7\xdc\x49\x32\x8c\x97\x08\xcb\xf5\x1c\xe1\xc4\x9b\x33\xad\x08\x87\x78\x01\x98\xae\x66\x10\x06\x5b\x2a\x77\x6a\x13\x44\x10\x54\xbc\x39\x6d\x9a\x8b\x5e\x4a\x5e\x9c\xa9\xd0\xab\x3d\x61\x07\xca\x44\x30\x7a\xf1\x7d\xab\x14\xbf\x4e\x11\x8e\x05\x17\xea\x44\x73\x25\x5a\xb9\xd0\xf7\x68\xe5\xb9\xaf\x3d\xa3\x45\x0d\xf3\x6c\x32\x8b\xb3\x35\xbc\xe3\x3a\xf2\x3d\x56\x1c\x89\x2b\x91\xe4\x22\xdb\xc3\xa2\x2c\x89\x10\xb9\x6c\x0e\x84\xb9\x0b\xff\xc1\xce\x92\xe7\xfd\x1f\xb4\xaf\xbd\xb0\xba\x76\x1f\xb9\x2b\xc3\xa8\xdd\x61\xd8\x17\xfa\x44\xe0\xe9\xea\xbc\xd7\x0e\x65\x92\x6c\xdb\x26\x33\x7c\xc3\x0c\xd3\x04\x17\x77\xcd\x87\xb4\x1a\x69\xde\x0e\x77\x95\x4e\x3e\x57\x08\x93\x74\x8c\x5f\x37\xea\xab\x9b\x31\xbb\x06\xf6\x91\xfe\xd6\xd1\x8d\x1f\x3a\x8e\xa1\x28\xf6\xa4\x94\xcf\x85\xbe\x2b\x58\x55\x93\x9f\xa9\x8b\x46\xf1\x92\xe4\x8a\xd7\x83\xe7\x77\x99\xda\xcd\xe3\x4c\x36\x8a\xd6\xd5\x73\x18\x1d\xb3\x8b\x75\x20\xd4\xae\x46\xe7\x69\xc0\xc5\xee\xe1\xb5\x18\xc7\x67\xf1\x4c\x7e\xbd\xb0\xfe\x31\x57\x3d\xa7\x31\x4e\xb1\xd5\x4c\xe2\x45\x12\x8f\xf1\xef\xbd\x0c\x6a\xf4\xd2\x81\xb0\x03\x8b\xc0\x29\x9b\x17\xf5\x1d\x00\x00\xff\xff\x0c\xfe\x65\x7d\xba\x03\x00\x00")

func dbMigrations0001_init_dbUpSqlBytes() ([]byte, error) {
	return bindataRead(
		_dbMigrations0001_init_dbUpSql,
		"db/migrations/0001_init_db.up.sql",
	)
}

func dbMigrations0001_init_dbUpSql() (*asset, error) {
	bytes, err := dbMigrations0001_init_dbUpSqlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "db/migrations/0001_init_db.up.sql", size: 954, mode: os.FileMode(420), modTime: time.Unix(1462993270, 0)}
	a := &asset{bytes: bytes, info: info}
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
	if err != nil {
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
	"db/migrations/0001_init_db.down.sql": dbMigrations0001_init_dbDownSql,
	"db/migrations/0001_init_db.up.sql": dbMigrations0001_init_dbUpSql,
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
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"db": &bintree{nil, map[string]*bintree{
		"migrations": &bintree{nil, map[string]*bintree{
			"0001_init_db.down.sql": &bintree{dbMigrations0001_init_dbDownSql, map[string]*bintree{}},
			"0001_init_db.up.sql": &bintree{dbMigrations0001_init_dbUpSql, map[string]*bintree{}},
		}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
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

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

