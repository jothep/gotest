package main
import (
  "fmt"
  "os"
  "flag"
  "strings"
  "strconv"
  "encoding/csv"
  "encoding/json"
  "io"
)
type User struct {
  Id int `json:"id"`
  Name string `json:"name"`
  Home string `json:"home"`
  Shell string `json:"shell"`
}

func main() {
  path, format := parseFlags()
  users := collectUsers()

  var output io.Writer
  if path != ""{
    f, err := os.Create(path)
    handleError(err)
    defer f.Close()
    output = f
  } else {
    output = os.Stdout
  }

  if format == "json" {
    data, err :=  json.MarshalIndent(users, "", "\t")
    handleError(err)
    output.Write(data)
  } else if format == "csv" {
    output.Write([]byte("name,id,home,shell\n"))
    writer := csv.NewWriter(output)
    for _, user := range users {
      err := writer.Write([]string{user.Name, strconv.Itoa(user.Id), user.Home, user.Shell})
      handleError(err)
    }
    writer.Flush()
  }
}

func parseFlags() (path, format string) {
  flag.StringVar(&path, "path", "", "the path to export file.")
  flag.StringVar(&format, "format", "json", "the output format for the user information.Avaliable options are 'csv' and 'json'. The default option is json.")
  flag.Parse()

  format = strings.ToLower(format)
  if format != "csv" && format != "json" {
    fmt.Println("Error: invalid format. Use 'json' or 'csv' instead.")
    flag.Usage()
    os.Exit(1)
  }
  return
}

func handleError(err error) {
  if err != nil {
  fmt.Println("Error:", err)
  os.Exit(1)
  }
}

func collectUsers()(users []User) {
  f, err := os.Open("/etc/passwd")
  handleError(err)
  defer f.Close()

  reader := csv.NewReader(f)
  reader.Comma = ':'

  lines, err := reader.ReadAll()
  handleError(err)

  for _, line := range lines {
    id, err := strconv.ParseInt(line[2], 10, 64)
    handleError(err)

    if id < 1000 {
      continue
    }

    user := User{
      Name: line[0],
      Id: int(id),
      Home: line[5],
      Shell: line[6],
    }

    users = append(users, user)
  }
  return
  
}
