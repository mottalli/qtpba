package qtpba

import (
	"path"
)

func GetBaseDir() string {
	return "/home/marcelo/Documents/Programacion/qtpba"
}

func GetFullPath(thePath string) string {
	return path.Join(GetBaseDir(), thePath)
}
