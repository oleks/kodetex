package kodetex

import (
  "os"
)

func ChangeToCurrentProjectDirectory() {
  error := os.Chdir(ProjectRoot)
  PanicOn(error)
}
