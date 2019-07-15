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

var _indexHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x54\x4d\x6b\x1b\x31\x10\x3d\x57\xbf\x62\xaa\xd3\x9a\x98\x5d\x48\x6e\x8d\x76\x0f\x89\x4d\x6b\x70\x93\x42\xd3\x43\xa1\x17\x59\x1a\xef\x0a\x64\xc9\x48\x5a\xbb\xa1\xe4\xbf\x17\x7d\xb8\x4e\x08\x76\x4b\xc9\x5e\x56\x1f\x6f\x34\x6f\xde\x3c\x89\xbd\x9f\xdd\xdf\x3e\x7c\xff\x32\x87\x21\x6c\x74\x47\x58\xf9\xad\xac\x7c\xec\x08\x61\xc3\x25\x28\xd9\x52\x3f\x70\xda\xf9\x81\xc3\x0a\x95\xe9\x61\x35\x2a\x1d\x58\x33\x5c\xc6\x80\xab\x84\x08\xdc\xf5\x18\x68\x27\x46\xe7\xd0\x04\xc8\x73\xd6\x0c\x57\x1d\x61\xc2\x4a\x4c\x28\x3b\x86\xed\x18\x68\xc7\x9a\xb8\x14\x13\x78\xe1\xd4\x36\x74\x64\xc7\x1d\xe0\x2e\x7c\xb5\xa3\x13\x08\x2d\x18\xdc\xc3\x7c\x87\xa6\xac\x54\xb4\xc1\x38\xf3\xb5\x4f\x73\x3a\xb9\x26\x29\x46\x2b\x1f\xd0\x40\x0b\xeb\xd1\x88\xa0\xac\xa9\x70\x02\xbf\x08\x40\xdc\x34\xb8\x9f\x6b\xdc\x44\x3e\x2d\x48\x2b\xc6\x38\xac\x85\x43\x1e\xb0\x6c\x54\x34\xe0\xcf\x10\x8f\xcb\x21\x39\xcb\x73\x78\x8f\xa1\x60\x6f\x1e\x17\xb2\x3a\xd4\x10\x09\x00\x08\x6b\xbc\xd5\x58\x6b\xdb\x57\x98\x97\x8e\x49\x6b\x65\x0c\xba\x4f\x0f\x9f\x97\xd0\x02\xd6\x92\x07\x0e\x17\x40\x7f\x18\xd6\xac\x5c\x47\x63\xca\x52\x14\xdf\x6e\xd1\xc8\xdb\x41\x69\x59\x1d\xe3\x27\xd7\xe4\x29\x57\x99\xd5\x3c\x51\xe5\xdf\x29\x97\xe6\x9c\xa4\x5c\x58\x3c\xa7\x5b\x62\x3e\x00\x85\x8b\xc2\xfd\x0f\x1b\xa1\x91\xbb\x57\x64\xde\xbd\x3c\xb8\x90\x8b\xa6\x39\xc3\x2c\x1a\x2b\x69\xef\x07\xfe\x52\xae\x83\x17\x6a\xcd\x7d\x48\x4e\x58\xc8\xa8\x1e\xdc\x7c\x5b\x2c\x67\x8b\xbb\x8f\xf4\x7f\x5a\x06\xb0\x1f\x94\x46\xa8\x4a\xc9\x6b\xe5\x7c\x48\xba\x67\x3d\xe3\x57\xb6\x1c\x6e\xec\x0e\x73\x4f\x5e\xa3\xe3\x51\x4f\x07\x41\xa4\x35\x78\xa2\x39\x6f\x5f\xff\xec\xfe\x6e\x4e\x53\x2f\x8e\x18\x2e\x65\x82\x2c\xd3\x6d\x40\x57\x51\x6d\x7b\x3a\x2d\xb7\x63\x72\x16\x19\x9d\xa7\x4c\x44\xa7\xbe\x9e\x07\xfb\x51\x08\xf4\x9e\x4e\x53\xcd\xe7\xb1\x6b\xae\xf4\xe8\xf0\x9f\xb0\xc5\xa1\xd3\xe2\xf4\x09\x21\xac\x29\x4f\x03\xc4\x71\x7e\x8f\x58\x93\x9e\xa7\xdf\x01\x00\x00\xff\xff\xb8\xee\x71\x92\xb5\x04\x00\x00")

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

	info := bindataFileInfo{name: "index.html", size: 1205, mode: os.FileMode(420), modTime: time.Unix(1563289201, 0)}
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
