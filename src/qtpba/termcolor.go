package qtpba

import "fmt"

const (
	colorRed    = 31
	colorGreen  = 32
	colorYellow = 33
	colorString = "\033[%dm%v\033[0m"
)

func RED(v interface{}) string {
	return fmt.Sprintf(colorString, colorRed, v)
}

func GREEN(v interface{}) string {
	return fmt.Sprintf(colorString, colorGreen, v)
}

func YELLOW(v interface{}) string {
	return fmt.Sprintf(colorString, colorYellow, v)
}
