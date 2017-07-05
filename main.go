// main
package main

import (
	"fmt"
	"os"

	"github.com/therecipe/qt/widgets"
)

func main() {
	fmt.Println("Hello World!")
	widgets.NewQApplication(len(os.Args), os.Args)
	mainWindows := NewMainwindow()
	mainWindows.Show()
	widgets.QApplication_Exec()
}
