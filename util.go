package kodetex

import (
  "log"
)

func PanicOn(error error) {
  if error != nil {
    panic(error)
  }
}

func LogOn(error error) {
  if error != nil {
    log.Println(error)
  }
}
