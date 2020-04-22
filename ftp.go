package main

import (
  "bufio"
  "flag"
  "fmt"
  "log"
  "os"
  "strings"
  "time"

  "github.com/jlaffaye/ftp"
)

func main() {
  // parse args
  host := flag.String("host", "127.0.0.1", "ftp host")
  port := flag.String("port", "21", "ftp port")
  username := flag.String("username", "ftpuser", "ftp username")
  password := flag.String("password", "ftpassword", "ftp password")
  product := flag.String("product", "", "product name as sub dir")
  version := flag.String("version", "", "version as sub dir")
  upload_file := flag.String("file", "", "file need to be uploaded")
  flag.Parse()

  if *product == "" || *version == "" || *upload_file == "" {
    flag.Usage()
    log.Fatal("ERROR: --product, --version, --file can not be empty")
  }

  upload_dest_dir := []string{"App", "app", "auto_build", *product, "Android", *version}

  // ftp connect
  c, err := ftp.Dial(fmt.Sprintf("%s:%s", *host, *port), ftp.DialWithTimeout(5*time.Second))
  if err != nil {
    log.Fatal(err)
  }

  // ftp login
  err = c.Login(*username, *password)
  if err != nil {
    log.Fatal(err)
  }

  // go to upload destnation
  for i := 0; i < len(upload_dest_dir); i++ {
    err := c.ChangeDir(upload_dest_dir[i])
    if err != nil {
      log.Print(err)
      err := c.MakeDir(upload_dest_dir[i])
      if err != nil {
        log.Fatal(err)
      }
      log.Printf("%s created", upload_dest_dir[i])
      err = c.ChangeDir(upload_dest_dir[i])
      if err != nil {
        log.Fatal(err)
      }
    }
  }
  currentdir, _ := c.CurrentDir()
  log.Printf("Current dir: %s", currentdir)

  // upload file
  filename := strings.Split(*upload_file, "/")[len(strings.Split(*upload_file, "/"))-1]
  if err == nil {
    log.Printf("%s is deleted", filename)
  }

  file, err := os.Open(*upload_file)
  if err != nil {
    log.Fatal(err)
  } else {
    log.Printf("========== Start upload %s ============", *upload_file)
  }
  defer file.Close()

  reader := bufio.NewReader(file)

  err = c.Stor(filename, reader)
  if err != nil {
    log.Fatal(err)
  } else {
    log.Printf("========== End upload %s ============", *upload_file)
  }

  // ftp disconnection
  if err := c.Quit(); err != nil {
    log.Fatal(err)
  }
}
