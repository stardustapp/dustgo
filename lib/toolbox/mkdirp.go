package toolbox

import (
	"log"
	"strings"

	"github.com/stardustapp/dustgo/lib/base"
	"github.com/stardustapp/dustgo/lib/inmem"
)

func Mkdirp(ctx base.Context, path string) (ok bool) {
	names := strings.Split(strings.TrimPrefix(path, "/"), "/")
	path = ""
	for _, name := range names {
		path += "/" + name
		ent, ok := ctx.Get(path)
		if !ok {
			ok := ctx.Put(path, inmem.NewFolder(name))
			if !ok {
				log.Println("mkdirp: Failed to create", path)
				return false
			}
		} else if _, ok := ent.(base.Folder); !ok {
			log.Println("mkdirp:", path, "already exists, and isn't a Folder")
			return false
		}
	}
	return true
}
