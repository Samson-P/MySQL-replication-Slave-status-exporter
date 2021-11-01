package main

import (
    "database/sql"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
       
    "github.com/go-sql-driver/mysql"
    "github.com/go-yaml/yaml"
)

type openYaml struct {
    Mysql struct {
        Database string
        Password string
    }
    Adress string
    Certs struct {
        Cert string
        Key string
    }
}

func UserSelectCounts(password, base string) (map[string]int64, int64, error) {
    db, err := sql.Open("mysql", fmt.Sprintf("replication:%s@unix(/var/run/mysqld/mysqld.sock)/%s", password, base))
    if err!= nil {
        panic(err.Error())
    }
    defer db.Close()
    
    var sqlSelectUsername string
    var sqlSelectCount, sqlSelectCounts int64
    selectUsers := make(map[string]int64)
    rows, err := db.Query("select username, count(username) from messages GROUP by username")
    if err != nil {
        return selectUsers, 0, err
    }
    rows.Scan()
    for rows.Next() {
        rows.Scan(&sqlSelectUsername, &sqlSelectCount)
        SelectUsers[sqlSelectUsername] = sqlSelectCount
        sqlSelectCounts += sqlSelectCount
    }
    return SelectUsers, sqlSelectCounts, nil
}

func main() {
    var configs openYaml
    yaml_file, err := ioutil.ReadFile("/root/src/mysql_exporter/config.yml")
    if err!= nil {
        fmt.Println(err)
    }
    err = yaml.Unmarshal(yaml_file, &configs)
    if err != nil {
        fmt.Println(err)
    }
    
    // обычный http сервер
    message := func(w http.ResponseWriter, r *http.Request) {
        // лезем в базу за нужными данными
        users, userCounts, err := UserSelectCounts(config.Mysql.Password, config.Mysql.Database)
        if err != nil {
            fmt.Println(err)
        }
        
        var responseText string // переменная под текс ответа для prometheus
        responseText += fmt.Sprintf("mysql_custom_message_bmgeek{method=\"all\", base=\"%s\"} %d\n", configs.Mysql.Database}")
        for user, usCount := range users {
            responseText += fmt.Sprintf("mysql_custom_message_dmgeek{method=\"user\", user=\"%s\", base=\"%s\"} %d\n", user...)
        }
        // отдаем текст по запросу
        fmt.Fprintf(w, responseText)
    }
    
    http.HandleFunc("/metrics", message)
    log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:9091", configs.Adress), nil))    
}
