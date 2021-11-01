package main

import (
    "database/sql"
    "fmt"
    // "io/ioutil"
    // "log"
    "net/http"
       
    _ "github.com/go-sql-driver/mysql"
    // "github.com/go-yaml"
)

type openYaml struct {
    Mysql struct {
        Database string
        Password string
    }
    Adress string
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
        selectUsers[sqlSelectUsername] = sqlSelectCount
        sqlSelectCounts += sqlSelectCount
    }
    return selectUsers, sqlSelectCounts, nil
}

func main() {
    // var configs openYaml
    // yaml_file, err := ioutil.ReadFile("/home/samson/Документы/проекты/mysql_email_error_alert/go_exporter|v1.0/config.yml")
    // if err!= nil {
    //     fmt.Println(err)
    // }
    // err = yaml.Unmarshal(yaml_file, &configs)
    // if err != nil {
    //     fmt.Println(err)
    // 
    
    // обычный http сервер
    message := func(w http.ResponseWriter, r *http.Request) {
        // лезем в базу за нужными данными
        users, userCounts, err := UserSelectCounts("tSj9ERz@5*DvSos?", "cnc_factory")
        // users, userCounts, err := UserSelectCounts(configs.Mysql.Password, configs.Mysql.Database)
        if err != nil {
            fmt.Println(err)
        }
        
        var responseText string // переменная под текс ответа для prometheus
        responseText += fmt.Sprintf("mysql_custom_message_bmgeek{method=\"all\", base=\"%s\"} %d\n", "cnc_factory", userCounts)
        // responseText += fmt.Sprintf("mysql_custom_message_bmgeek{method=\"all\", base=\"%s\"} %d\n", configs.Mysql.Database, userCounts)
        for user, usCount := range users {
            responseText += fmt.Sprintf("mysql_custom_message_dmgeek{method=\"user\", user=\"%s\", base=\"%s\"} %d\n", user, "cnc_factory", usCount)
            // responseText += fmt.Sprintf("mysql_custom_message_dmgeek{method=\"user\", user=\"%s\", base=\"%s\"} %d\n", user, configs.Mysql.Database, usCount)
        }
        // отдаем текст по запросу
        fmt.Fprintf(w, responseText)
    }
    
    http.HandleFunc("/metrics", message)
    http.ListenAndServe(":9092", nil)
}
