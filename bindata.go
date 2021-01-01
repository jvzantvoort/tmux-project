// Code generated by go-bindata. DO NOT EDIT.
// sources:
// messages/bash (1.319kB)
// messages/usage_archive (96B)
// messages/usage_create (92B)
// messages/usage_edit (76B)
// messages/usage_list (96B)
// messages/usage_projectinit (73B)
// messages/usage_shell (4B)

package tmuxproject

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
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
		return nil, fmt.Errorf("read %q: %w", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("read %q: %w", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes  []byte
	info   os.FileInfo
	digest [sha256.Size]byte
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

var _messagesBash = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xe4\x94\xd1\x6f\xd3\x30\x10\xc6\xdf\xfd\x57\x7c\xb8\xd1\xba\x22\xa5\x69\x2a\x3a\x09\xaa\x20\x10\x2a\x82\x87\x6e\x12\xe5\xad\xaa\x2a\xd7\xb9\xac\x06\xd7\xa9\xe2\x64\x1b\xa4\xf9\xdf\x91\x5b\x8f\x35\xa3\x15\xbc\xe3\xa7\x38\xf7\xdd\x77\xe7\x5f\x72\xee\xbc\x88\x56\xca\x44\x2b\x61\xd7\x8c\x75\xec\x9a\xb4\x96\x6b\x92\xdf\x91\x2a\x2b\x56\x9a\x92\xd9\x87\x61\x3c\x1a\xb1\x94\xa4\x16\x05\x21\x7c\xc0\x6c\x32\x9b\x7d\xbe\xb9\xbe\x7e\x3f\x9d\x24\xc1\x65\xb9\xa9\x1e\x9c\x78\xab\xc5\x8f\x70\x43\xd6\x8a\x5b\x42\xb8\x45\xb7\x33\xeb\xf6\xce\x59\xc6\x83\xd7\x03\xa6\x32\xcc\x11\x66\xe0\xc1\xa7\x9b\xe9\x24\xea\xbb\x1e\x22\x67\xd7\x4f\xa3\xa0\x3e\xaa\xd2\xf4\xc9\xdc\x71\x2c\x58\xb9\x26\xc3\x00\x9b\x57\x85\x24\xfc\x6b\x22\xcb\x14\x63\x59\x65\x64\xa9\x72\x83\xe5\xd2\x29\xb7\x4b\x6d\x2f\x7b\xac\x66\xc0\xfe\x00\xda\x22\xfc\x88\x6e\xa7\xb6\x64\xad\xca\xcd\xd2\x88\x0d\x35\x5d\x0c\xdf\x46\x29\xdd\x45\xa6\xd2\x9a\x35\x27\x5c\x94\xf1\x2e\x3a\x97\x42\x83\x18\x90\xe5\x05\x08\xca\x80\x07\xf5\xbb\x37\xc3\x86\x8f\x91\xe6\x98\xcf\xc1\x03\xe2\x48\x12\xf0\x20\xe6\x58\x2c\x70\x71\x81\x82\xca\xaa\x30\x18\x38\x89\x71\xc9\xfe\x45\x7c\xb2\x98\x28\x4b\x21\xd7\xad\x82\xbe\xdd\xdf\x7b\x99\x9b\x4c\xdd\x66\x4a\x13\x73\xa0\x0e\xd1\x24\x88\xc7\x0c\x47\xb1\x84\x07\xb5\x43\xd7\x3c\x67\xe7\x13\x9a\x7e\x21\xb9\x33\x98\xfb\x0f\x54\x3f\xe5\x36\xfb\xde\x77\xbb\xbf\xdb\xa5\x94\x89\x4a\x97\x8f\x5e\x5f\x27\x5f\xa6\x09\xb7\xb2\x20\x32\xe1\x70\x74\x25\x73\x9d\x17\xfc\x80\xff\x44\x91\xc3\x61\x11\xa6\x08\x4b\x17\x7c\x6c\x8d\x9f\x44\x63\xe8\xfe\x3f\xe1\x62\xe8\x1e\xa1\x3d\x4f\xa4\x20\x5b\x6d\xe8\x0c\x0c\x29\x2c\x21\x88\xa1\xdc\x06\xd0\x76\xa7\x95\x2d\x97\x5e\x60\x7b\xfb\xaa\xe1\xb6\xc8\xbf\x91\x2c\xe1\x62\x08\xb3\xf1\x78\x2f\x7e\xd9\x7b\xc6\xcd\x2d\x07\xe2\x67\xab\x99\xf6\x9f\x1d\x7b\xdd\x99\x9b\x65\xf0\xea\xca\x0b\x54\x76\x34\x53\x6d\xc3\xe0\xf2\x69\x66\x7b\x5e\xee\x6f\x82\xc3\x6a\xcf\x47\x1b\xcd\x41\x41\xda\xd2\x1f\x72\x47\xf2\x84\x36\x53\xfe\x61\x7f\x6e\xb2\x42\xb2\x86\xfd\x0a\x00\x00\xff\xff\x01\xe7\xb9\xc0\x27\x05\x00\x00")

func messagesBashBytes() ([]byte, error) {
	return bindataRead(
		_messagesBash,
		"messages/bash",
	)
}

func messagesBash() (*asset, error) {
	bytes, err := messagesBashBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "messages/bash", size: 1319, mode: os.FileMode(0664), modTime: time.Unix(1609536830, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xd4, 0x22, 0x52, 0xa8, 0x10, 0xf8, 0xfa, 0xa6, 0x4d, 0x17, 0x92, 0xb7, 0xd0, 0x5, 0x6f, 0xc9, 0x72, 0xa3, 0xf4, 0xff, 0x78, 0x22, 0x34, 0x3e, 0x20, 0xe8, 0x36, 0x66, 0xa1, 0x71, 0x14, 0xed}}
	return a, nil
}

var _messagesUsage_archive = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x2c\x4a\xce\xc8\x2c\x4b\x55\x48\x54\x28\x28\xca\xcf\x4a\x4d\x2e\xe1\xe2\x4a\x84\x0a\x45\xeb\x96\xc5\x2a\x44\xeb\x26\xd6\x40\x05\xf2\x12\x73\x53\x15\x6c\x90\x38\x76\xb1\x0a\xba\x79\x0a\x35\x0a\xba\x50\xad\x10\x05\x48\x1c\x3b\x05\x2e\x40\x00\x00\x00\xff\xff\xe7\x9a\xee\x67\x60\x00\x00\x00")

func messagesUsage_archiveBytes() ([]byte, error) {
	return bindataRead(
		_messagesUsage_archive,
		"messages/usage_archive",
	)
}

func messagesUsage_archive() (*asset, error) {
	bytes, err := messagesUsage_archiveBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "messages/usage_archive", size: 96, mode: os.FileMode(0664), modTime: time.Unix(1609532582, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0xcd, 0xca, 0x84, 0x4d, 0x4f, 0xa6, 0xde, 0x22, 0x56, 0xca, 0xa5, 0x47, 0x76, 0x16, 0x4a, 0xc0, 0xbc, 0x36, 0xea, 0x5e, 0xbc, 0x7a, 0x83, 0xb9, 0xf2, 0x8f, 0xaf, 0x8d, 0x60, 0x7f, 0x19, 0x9c}}
	return a, nil
}

var _messagesUsage_create = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x2e\x4a\x4d\x2c\x49\x55\x48\x54\xc8\x4b\x2d\x57\x28\x28\xca\xcf\x4a\x4d\x2e\xd1\xe3\xe2\x4a\x86\x08\x47\xeb\x96\xd4\xe8\x42\x45\x4b\x2a\x0b\x52\x15\x6c\x40\xa4\x5d\xac\x82\x6e\x9e\x82\x4d\x5e\x62\x6e\xaa\x9d\x42\x8d\x02\x4c\x01\x88\x0f\x13\x8d\xd6\x2d\x8b\xe5\x02\x04\x00\x00\xff\xff\xf0\x64\x3d\x02\x5c\x00\x00\x00")

func messagesUsage_createBytes() ([]byte, error) {
	return bindataRead(
		_messagesUsage_create,
		"messages/usage_create",
	)
}

func messagesUsage_create() (*asset, error) {
	bytes, err := messagesUsage_createBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "messages/usage_create", size: 92, mode: os.FileMode(0664), modTime: time.Unix(1609532582, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x6b, 0xbd, 0x87, 0x53, 0x8e, 0xa6, 0xe5, 0xac, 0xa1, 0x1c, 0x4e, 0xb7, 0x92, 0x8f, 0x6b, 0xe1, 0xe8, 0x8, 0x53, 0x44, 0xad, 0x3f, 0xbe, 0x12, 0x9b, 0x4, 0x92, 0x39, 0xc6, 0x8e, 0xde, 0x7a}}
	return a, nil
}

var _messagesUsage_edit = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x4d\xc9\x2c\x51\x48\x54\x28\x28\xca\xcf\x4a\x4d\x2e\x29\x56\x48\xcd\x2b\xcb\x2c\xca\xcf\xcb\x4d\xcd\x2b\x51\x48\xcc\x4b\x51\x28\xc9\x2d\xad\x50\x48\xce\xcf\x4b\xcb\x4c\x4f\xcb\xcc\x49\xe5\xe2\x4a\x05\x69\xd0\xcd\x53\xb0\x81\x6a\xc9\x4b\xcc\x4d\xb5\x53\x88\xd6\x2d\x8b\xe5\x02\x04\x00\x00\xff\xff\x87\x41\x86\x7c\x4c\x00\x00\x00")

func messagesUsage_editBytes() ([]byte, error) {
	return bindataRead(
		_messagesUsage_edit,
		"messages/usage_edit",
	)
}

func messagesUsage_edit() (*asset, error) {
	bytes, err := messagesUsage_editBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "messages/usage_edit", size: 76, mode: os.FileMode(0664), modTime: time.Unix(1609532582, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x97, 0x68, 0x22, 0x4d, 0x74, 0xb9, 0xcd, 0xc5, 0x7b, 0x92, 0x76, 0xe4, 0x93, 0xda, 0x34, 0xb, 0x15, 0xeb, 0x46, 0xbe, 0x7b, 0xf3, 0x4b, 0xd7, 0x47, 0x20, 0x3a, 0xb2, 0x35, 0xeb, 0xa8, 0x9}}
	return a, nil
}

var _messagesUsage_list = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x0a\xc9\x48\x55\x50\xca\xc9\x2c\x2e\x51\x52\x48\xce\xcf\xcd\x4d\xcc\x4b\x51\x00\xf1\x14\x0a\x8a\xf2\xb3\x52\x93\x4b\x8a\x15\x92\x4b\x8b\x8a\x8a\x52\xf3\x4a\x72\x2a\x15\x92\xf3\xf3\xd2\x32\xd3\x4b\x8b\x52\x53\xb8\xb8\xc0\x8a\xa2\x75\xa1\xca\xf2\x12\x73\x53\x6b\x74\xf3\x14\x6c\x40\x0c\xbb\x58\x85\x68\xdd\x32\x08\xc1\x05\x08\x00\x00\xff\xff\xbc\x87\x8c\xe4\x60\x00\x00\x00")

func messagesUsage_listBytes() ([]byte, error) {
	return bindataRead(
		_messagesUsage_list,
		"messages/usage_list",
	)
}

func messagesUsage_list() (*asset, error) {
	bytes, err := messagesUsage_listBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "messages/usage_list", size: 96, mode: os.FileMode(0664), modTime: time.Unix(1609532582, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x91, 0x6f, 0xd6, 0xd5, 0xd2, 0xae, 0xf5, 0x5f, 0x25, 0x58, 0xf4, 0x17, 0xfb, 0xce, 0x2, 0xfd, 0x8f, 0xd0, 0x14, 0x74, 0xd5, 0xf1, 0x16, 0xf9, 0xd0, 0x71, 0x92, 0x2d, 0x1a, 0x10, 0x43, 0x5d}}
	return a, nil
}

var _messagesUsage_projectinit = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xf2\xcc\xcb\x2c\xc9\x4c\xcc\xc9\xac\x4a\x55\x48\x54\xc8\x4b\x2d\x57\x28\x28\xca\xcf\x4a\x4d\x2e\x51\x28\xa9\x2c\x48\xe5\xe2\xca\xcc\xcb\x2c\x51\x88\xd6\x2d\x8b\x55\xd0\x2d\x51\xa8\x51\xd0\x85\xca\x82\x24\x15\x6c\x90\x38\x76\x5c\x80\x00\x00\x00\xff\xff\xda\x9f\xd3\xa8\x49\x00\x00\x00")

func messagesUsage_projectinitBytes() ([]byte, error) {
	return bindataRead(
		_messagesUsage_projectinit,
		"messages/usage_projectinit",
	)
}

func messagesUsage_projectinit() (*asset, error) {
	bytes, err := messagesUsage_projectinitBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "messages/usage_projectinit", size: 73, mode: os.FileMode(0664), modTime: time.Unix(1609532582, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x18, 0xfa, 0xc6, 0x80, 0x1b, 0x3e, 0xa9, 0x84, 0x4, 0xc3, 0x3c, 0xf, 0xf5, 0x3e, 0xd2, 0x7b, 0xb6, 0x20, 0x63, 0x6f, 0x75, 0x1, 0x26, 0x92, 0x2c, 0xcf, 0x90, 0xd1, 0x6, 0x3d, 0xe8, 0x36}}
	return a, nil
}

var _messagesUsage_shell = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x0a\x71\x72\xe1\x02\x04\x00\x00\xff\xff\x58\x7d\xf5\x62\x04\x00\x00\x00")

func messagesUsage_shellBytes() ([]byte, error) {
	return bindataRead(
		_messagesUsage_shell,
		"messages/usage_shell",
	)
}

func messagesUsage_shell() (*asset, error) {
	bytes, err := messagesUsage_shellBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "messages/usage_shell", size: 4, mode: os.FileMode(0664), modTime: time.Unix(1609536625, 0)}
	a := &asset{bytes: bytes, info: info, digest: [32]uint8{0x76, 0xf6, 0x2, 0xb1, 0xdf, 0xc2, 0xd5, 0x48, 0xfa, 0x1f, 0x5e, 0xdd, 0xf2, 0x93, 0x15, 0xc9, 0x47, 0xbe, 0xc3, 0x1f, 0x71, 0xc9, 0x5d, 0x49, 0x65, 0x39, 0x27, 0x20, 0x44, 0x86, 0xbe, 0x39}}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// AssetString returns the asset contents as a string (instead of a []byte).
func AssetString(name string) (string, error) {
	data, err := Asset(name)
	return string(data), err
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

// MustAssetString is like AssetString but panics when Asset would return an
// error. It simplifies safe initialization of global variables.
func MustAssetString(name string) string {
	return string(MustAsset(name))
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetDigest returns the digest of the file with the given name. It returns an
// error if the asset could not be found or the digest could not be loaded.
func AssetDigest(name string) ([sha256.Size]byte, error) {
	canonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[canonicalName]; ok {
		a, err := f()
		if err != nil {
			return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s can't read by error: %v", name, err)
		}
		return a.digest, nil
	}
	return [sha256.Size]byte{}, fmt.Errorf("AssetDigest %s not found", name)
}

// Digests returns a map of all known files and their checksums.
func Digests() (map[string][sha256.Size]byte, error) {
	mp := make(map[string][sha256.Size]byte, len(_bindata))
	for name := range _bindata {
		a, err := _bindata[name]()
		if err != nil {
			return nil, err
		}
		mp[name] = a.digest
	}
	return mp, nil
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
	"messages/bash":              messagesBash,
	"messages/usage_archive":     messagesUsage_archive,
	"messages/usage_create":      messagesUsage_create,
	"messages/usage_edit":        messagesUsage_edit,
	"messages/usage_list":        messagesUsage_list,
	"messages/usage_projectinit": messagesUsage_projectinit,
	"messages/usage_shell":       messagesUsage_shell,
}

// AssetDebug is true if the assets were built with the debug flag enabled.
const AssetDebug = false

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"},
// AssetDir("data/img") would return []string{"a.png", "b.png"},
// AssetDir("foo.txt") and AssetDir("notexist") would return an error, and
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		canonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(canonicalName, "/")
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
	"messages": {nil, map[string]*bintree{
		"bash": {messagesBash, map[string]*bintree{}},
		"usage_archive": {messagesUsage_archive, map[string]*bintree{}},
		"usage_create": {messagesUsage_create, map[string]*bintree{}},
		"usage_edit": {messagesUsage_edit, map[string]*bintree{}},
		"usage_list": {messagesUsage_list, map[string]*bintree{}},
		"usage_projectinit": {messagesUsage_projectinit, map[string]*bintree{}},
		"usage_shell": {messagesUsage_shell, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory.
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
	return os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
}

// RestoreAssets restores an asset under the given directory recursively.
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
	canonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(canonicalName, "/")...)...)
}