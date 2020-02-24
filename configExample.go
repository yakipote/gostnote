package main

import (
	"time"

	"github.com/jessevdk/go-assets"
)

var _Assets1be4991216147bffc0fef814d096ae692e1d953a = "[Firebase]\nkeyPath = \"path/to/key/json\"\nstorageBucket = \"firebase bucket host\"\n"

// Assets returns go-assets FileSystem
var Assets = assets.NewFileSystem(map[string][]string{"/": []string{"configExample.toml"}}, map[string]*assets.File{
	"/": &assets.File{
		Path:     "/",
		FileMode: 0x800001ed,
		Mtime:    time.Unix(1577025103, 1577025103322916958),
		Data:     nil,
	}, "/configExample.toml": &assets.File{
		Path:     "/configExample.toml",
		FileMode: 0x1b4,
		Mtime:    time.Unix(1574959117, 1574959117176805351),
		Data:     []byte(_Assets1be4991216147bffc0fef814d096ae692e1d953a),
	}}, "")
