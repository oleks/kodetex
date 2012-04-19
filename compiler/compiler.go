package main

import (
  "io"
  "log"
  "os"
  "path"
  "regexp"
  "calmach/kodetex"
)

type Settings struct {
  extensionRegex *regexp.Regexp
}

func main() {
  settings := loadSettings();

  file, error := os.Open(".");
  kodetex.PanicOn(error);

  var fileInfos []os.FileInfo;
  var fileInfo os.FileInfo;
  var name, extension string;

  for {

    fileInfos, error = file.Readdir(1);
    if error != nil {
      if error != io.EOF {
        panic(error);
      }
      break;
    }

    fileInfo = fileInfos[0];
    if fileInfo.IsDir() {
      continue;
    }

    name = fileInfo.Name();

    extension = path.Ext(name);
    if !settings.extensionRegex.MatchString(extension) {
      continue;
    }

    log.Printf("Parsing %s.", fileInfo.Name());

    parseFile(name);
  }
}

func loadSettings() (settings *Settings) {

  kodetex.ChangeToCurrentProjectDirectory();

  extensionRegex := regexp.MustCompile(kodetex.FileExtensions); 
  
  return &Settings{
    extensionRegex: extensionRegex,
  };
}

func parseFile(relativePath string) {

  defer func() {
    if error := recover() ; error != nil {
      log.Printf(
        "Failed to parse file %s. The following error occured:\n%s.\n",
        relativePath, error);
    }
  }()

  file, error := os.Open(relativePath);
  kodetex.PanicOn(error);

  parser := &kodetex.CppParser {
    File: file,
  };

  parser.Parse();
}
