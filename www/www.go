// Code generated for package www by go-bindata DO NOT EDIT. (@generated)
// sources:
// www/index.html
package www

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net/http"
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

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}


type assetFile struct {
	*bytes.Reader
	name            string
	childInfos      []os.FileInfo
	childInfoOffset int
}

type assetOperator struct{}

// Open implement http.FileSystem interface
func (f *assetOperator) Open(name string) (http.File, error) {
	var err error
	if len(name) > 0 && name[0] == '/' {
		name = name[1:]
	}
	content, err := Asset(name)
	if err == nil {
		return &assetFile{name: name, Reader: bytes.NewReader(content)}, nil
	}
	children, err := AssetDir(name)
	if err == nil {
		childInfos := make([]os.FileInfo, 0, len(children))
		for _, child := range children {
			childPath := filepath.Join(name, child)
			info, errInfo := AssetInfo(filepath.Join(name, child))
			if errInfo == nil {
				childInfos = append(childInfos, info)
			} else {
				childInfos = append(childInfos, newDirFileInfo(childPath))
			}
		}
		return &assetFile{name: name, childInfos: childInfos}, nil
	} else {
		// If the error is not found, return an error that will
		// result in a 404 error. Otherwise the server returns
		// a 500 error for files not found.
		if strings.Contains(err.Error(), "not found") {
			return nil, os.ErrNotExist
		}
		return nil, err
	}
}

// Close no need do anything
func (f *assetFile) Close() error {
	return nil
}

// Readdir read dir's children file info
func (f *assetFile) Readdir(count int) ([]os.FileInfo, error) {
	if len(f.childInfos) == 0 {
		return nil, os.ErrNotExist
	}
	if count <= 0 {
		return f.childInfos, nil
	}
	if f.childInfoOffset+count > len(f.childInfos) {
		count = len(f.childInfos) - f.childInfoOffset
	}
	offset := f.childInfoOffset
	f.childInfoOffset += count
	return f.childInfos[offset : offset+count], nil
}

// Stat read file info from asset item
func (f *assetFile) Stat() (os.FileInfo, error) {
	if len(f.childInfos) != 0 {
		return newDirFileInfo(f.name), nil
	}
	return AssetInfo(f.name)
}

// newDirFileInfo return default dir file info
func newDirFileInfo(name string) os.FileInfo {
	return &bindataFileInfo{
		name:    name,
		size:    0,
		mode:    os.FileMode(2147484068), // equal os.FileMode(0644)|os.ModeDir
		modTime: time.Time{}}
}

// AssetFile return a http.FileSystem instance that data backend by asset
func AssetFile() http.FileSystem {
	return &assetOperator{}
}

var _indexHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xb4\x54\x4d\x6b\x1b\x31\x10\xbd\xeb\x57\x4c\x75\xda\x25\x66\x17\x92\x5b\xa3\xdd\x43\x62\xd3\x1a\xdc\xa4\xd0\xf4\x50\xe8\x45\x96\x26\x96\x40\x96\x8c\xa4\xb5\x1b\x4a\xfe\x7b\xd1\x47\x9c\x40\xb1\x0b\x85\xec\x65\xf5\xf1\x46\xf3\xde\xcc\x93\xd8\x87\xf9\xfd\xed\xc3\x8f\xaf\x0b\x50\x71\x6b\x46\xc2\xea\x6f\xed\xe4\xd3\x48\x08\x53\x97\xa0\xe5\x40\x83\xe2\x74\x0c\x8a\xc3\x1a\xb5\xdd\xc0\x7a\xd2\x26\xb2\x5e\x5d\xa6\x80\xab\x8c\x88\xdc\x6f\x30\xd2\x51\x4c\xde\xa3\x8d\x50\xe6\xac\x57\x57\x23\x61\xc2\x49\xcc\x28\x37\xc5\xdd\x14\xe9\xc8\xfa\xb4\x94\x12\x04\xe1\xf5\x2e\x8e\x64\xcf\x3d\xe0\x3e\x7e\x73\x93\x17\x08\x03\x58\x3c\xc0\x62\x8f\xb6\xae\x34\xb4\xc7\x34\x0b\x5d\xc8\x73\xda\x5e\x93\x1c\x63\x74\x88\x68\x61\x80\xc7\xc9\x8a\xa8\x9d\x6d\xb0\x85\xdf\x04\x20\x6d\x5a\x3c\x2c\x0c\x6e\x13\x9f\x01\xa4\x13\x53\x1a\x76\xc2\x23\x8f\x58\x37\x1a\x1a\xf1\x57\x4c\xc7\x95\x90\x92\xe5\x2d\x7c\x83\xb1\x62\x6f\x9e\x96\xb2\x79\xd1\x90\x08\x00\x08\x67\x83\x33\xd8\x19\xb7\x69\xb0\x2c\xbd\x26\xed\xb4\xb5\xe8\x3f\x3f\x7c\x59\xc1\x00\xd8\x49\x1e\x39\x5c\x00\xfd\x69\x59\xbf\xf6\x23\x4d\x29\xab\x28\xbe\xdb\xa1\x95\xb7\x4a\x1b\xd9\xbc\xc6\xb7\xd7\xe4\xb9\xa8\x2c\xd5\x3c\xa1\xf2\xdf\x94\x6b\x73\x4e\x52\xae\x2c\xde\xd2\xad\x31\x1f\x81\xc2\x45\xe5\x7e\x64\x23\x0c\x72\x7f\x82\x4c\x32\xc9\x19\x26\xc9\x48\xb9\xd6\x41\xf1\x13\xe5\x81\x9b\xef\xcb\xd5\x7c\x79\xf7\x89\xfe\x4f\x4f\x00\x0e\x4a\x1b\x84\xa6\x6a\x7a\xd4\x3e\xc4\x5c\xd8\xc2\x31\x7d\x75\xcb\xe3\xd6\xed\xb1\x14\xfd\x6f\x74\x3a\xea\xf9\x45\xb1\x74\x16\xdf\x4f\xf0\xfc\xfe\x6e\x41\x73\x75\x8f\x37\xa0\xe3\x52\x66\xfb\xaf\xb2\xbf\xd1\x1f\x35\xce\xaa\xe5\xdb\xb3\x60\x8b\x07\x3a\x2b\x8d\x3a\x0f\x4c\xca\xe8\x2c\x0b\x3c\x0f\xac\x1e\x9a\x55\x2f\xb6\x84\xb0\xbe\x5e\x5e\x48\xe3\xf2\x62\xb0\x3e\x3f\x20\x7f\x02\x00\x00\xff\xff\xf6\x80\x55\xb3\x57\x04\x00\x00")

func indexHtmlBytes() ([]byte, error) {
	return bindataRead(
		_indexHtml,
		"index.html",
	)
}

func indexHtml() (*asset, error) {
	bytes, err := indexHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "index.html", size: 1111, mode: os.FileMode(420), modTime: time.Unix(1562189197, 0)}
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
	"index.html": indexHtml,
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
	"index.html": &bintree{indexHtml, map[string]*bintree{}},
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
