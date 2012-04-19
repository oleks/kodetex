package main

import (
  "io"
  "log"
  "net/http"
  "os"
  "time"
  "calmach/kodetex"
)

const (
  serverHeader = "KodeTeX v/0.0 Server v/0.0"
  address = ":8080"
)

func main() {
  setWorkingDirectory()
  initializeServer()
}

func setWorkingDirectory() {
  error := os.Chdir(kodetex.DocRoot)
  if error != nil {
    log.Fatal(error)
  }
}

func initializeServer() {
  handler := new(Handler)
  server := &http.Server {
    Addr: address,
    Handler: handler,
    ReadTimeout: 10 * time.Second,
    WriteTimeout: 10 * time.Second,
    MaxHeaderBytes: 1 << 20,
  }
  log.Fatal(server.ListenAndServe())
}

type Handler struct{}

func (*Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
  header := writer.Header()
  header.Set("Server", serverHeader)

  file, error := os.Open("index.html")
  if error != nil {
    log.Printf("Couldn't open index..\n%s\n", error)
    return
  }

  io.Copy(writer, file)

  error = file.Close()
  if error != nil {
    log.Fatal("Couldn't close index..")
  }
}


