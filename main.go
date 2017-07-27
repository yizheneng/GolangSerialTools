// main
package main

import (
	"fmt"
	"os"

	"github.com/therecipe/qt/widgets"
)

func main() {
	fmt.Println("Hello World!")
	app := widgets.NewQApplication(len(os.Args), os.Args)
	mainWindows := NewMainwindow(app)

	mainWindows.Show()
	widgets.QApplication_Exec()
}
