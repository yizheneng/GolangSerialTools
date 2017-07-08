// main
package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/therecipe/qt/widgets"
)

func main() {
	qssFile := "stylesheet.css"
	fmt.Println("Hello World!")
	file, err := os.OpenFile(qssFile, os.O_RDONLY, 0666)
	defer file.Close()

	app := widgets.NewQApplication(len(os.Args), os.Args)
	mainWindows := NewMainwindow()

	if err == nil {
		qssString, err := ioutil.ReadAll(file)
		if err == nil {
			fmt.Println(string(qssString))
			app.SetStyleSheet(string(qssString))
		}
	}

	mainWindows.Show()
	widgets.QApplication_Exec()
}
