// Code generated for package publisher by go-bindata DO NOT EDIT. (@generated)
// sources:
// internal/publisher/error.md
// internal/publisher/failure.md
// internal/publisher/success.md
package publisher

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

var _errorMd = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x1c\xca\xb1\x0a\xc3\x20\x10\x06\xe0\xfd\x9e\xe2\x20\x7b\x5e\xa0\x6b\x32\x0b\xd2\xee\x8a\x1e\x8d\xa0\x9e\xc4\x5f\x28\x88\xef\x5e\xda\xf1\x83\x6f\xe3\xd3\x5a\x63\xcf\x83\xcd\xeb\x49\xc8\xf1\x7e\x50\xd7\x22\xb8\x52\x7d\xf3\xa8\xf2\x69\x12\x20\x91\x2f\xdf\x9a\x54\x89\x44\xdb\xc6\x41\x4b\xcb\x02\x61\x1d\x68\x03\x04\x85\xcf\x8c\x54\x84\xe7\xdc\x8f\x71\x7b\x24\xad\x6b\x71\x97\xa0\x35\x76\x22\xe7\x1c\xcd\xb9\x9b\xff\x5f\xeb\xe7\x6f\x00\x00\x00\xff\xff\xcf\x3e\x0d\x7e\x7b\x00\x00\x00")

func errorMdBytes() ([]byte, error) {
	return bindataRead(
		_errorMd,
		"error.md",
	)
}

func errorMd() (*asset, error) {
	bytes, err := errorMdBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "error.md", size: 123, mode: os.FileMode(420), modTime: time.Unix(1563322178, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _failureMd = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x3c\xcc\x31\x0e\xc2\x30\x0c\x40\xd1\xdd\xa7\xb0\x94\xa1\x5b\x2f\xc0\x8a\x58\x19\x38\x00\xb1\x1a\x83\x2c\x39\x75\x94\x3a\x0b\x21\x77\x47\x0a\x52\xf7\xff\x5f\xc0\x1b\x89\x72\x02\xd7\x54\x2f\x10\x63\x84\x80\xef\xca\x05\x97\x17\x89\x2e\xd0\xfb\x7a\x6f\x5e\x9a\xe3\x17\x1f\x2d\x67\xaa\xf2\xe1\x31\x66\x09\x21\xe0\x66\xb9\x28\x3b\x3f\x6d\x56\xe0\xe6\xa4\xe8\x92\x19\x7b\x5f\xaf\xad\x92\x8b\xed\x63\xe0\xc1\x9b\xed\xe9\x80\x79\x9e\xea\x5f\xfa\x05\x00\x00\xff\xff\xe5\xa4\x3f\x8c\x86\x00\x00\x00")

func failureMdBytes() ([]byte, error) {
	return bindataRead(
		_failureMd,
		"failure.md",
	)
}

func failureMd() (*asset, error) {
	bytes, err := failureMdBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "failure.md", size: 134, mode: os.FileMode(420), modTime: time.Unix(1563330808, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _successMd = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x1c\xca\x31\x0a\x03\x30\x08\x00\xc0\xdd\x57\x08\xee\x79\x45\xf7\x0e\xfd\x40\x82\x75\x08\x24\x31\x54\x9d\xc4\xbf\x97\x76\x3c\x38\xc2\x57\x30\x8b\x19\x10\x11\xb2\xee\xbb\xc4\x05\x35\xfc\x86\x83\xab\x8f\x85\x3e\xb7\x60\x66\x7b\xc4\x67\xf8\xd4\x53\x85\x26\xac\xe7\x6d\x00\xbd\x77\xc8\x6c\xcf\xff\xaf\xfa\xf9\x1b\x00\x00\xff\xff\x87\x85\x33\xa2\x53\x00\x00\x00")

func successMdBytes() ([]byte, error) {
	return bindataRead(
		_successMd,
		"success.md",
	)
}

func successMd() (*asset, error) {
	bytes, err := successMdBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "success.md", size: 83, mode: os.FileMode(420), modTime: time.Unix(1563322064, 0)}
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
	"error.md":   errorMd,
	"failure.md": failureMd,
	"success.md": successMd,
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
	"error.md":   &bintree{errorMd, map[string]*bintree{}},
	"failure.md": &bintree{failureMd, map[string]*bintree{}},
	"success.md": &bintree{successMd, map[string]*bintree{}},
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
