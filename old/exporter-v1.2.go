package main

import (
    "database/sql"
    "fmt"
    "net/http"

    _ "github.com/go-sql-driver/mysql"
)


func MySQLSlaveStatus(password string) (map[string]int64, error) {
    db, err := sql.Open("mysql", fmt.Sprintf("replication:%s@/%s", password, "cnc_factory"))
    if err!= nil {
        panic(err.Error())
    }
    defer db.Close()
    
    var sqlSelectStatus, myQuery string
    var sqlSelectValue int64
    myQuery =  // "SHOW SLAVE STATUS\\G"
    // fmt.Println(myQuery)
    
    columns, err := rows.Columns()
    if err != nil {
        log.Fatal(err)
    }
    
    selectStatuses := make(map[string]int64)
    rows, err := db.Query(myQuery)
    if err != nil {
        return selectStatuses, err
    }
    rows.Scan()
    for rows.Next() {
        rows.Scan(&sqlSelectStatus, &sqlSelectValue)
        selectStatuses[sqlSelectStatus] = sqlSelectValue
    }
    return selectStatuses, nil
}

func main() {
    // http сервер
    message := func(w http.ResponseWriter, r *http.Request) {
        // лезем в базу за нужными данными
        replicationStatuses, err := MySQLSlaveStatus("tSj9ERz@5*DvSos?")
        if err != nil {
            fmt.Println(err)
        }
        // переменная под текс ответа mysql
        var responseText string // переменная под текс ответа для prometheus
        responseText += fmt.Sprintf("MAIN mysql_custom_messages\n")
        for status, value := range replicationStatuses {
            responseText += fmt.Sprintf("mysql_custom_message{method=\"slave\", status=\"%s\"} %d\n", status, value)
        }
        
        // отдаем текст по запросу
        fmt.Fprintf(w, responseText)
    }

    http.HandleFunc("/metrics", message)
    http.ListenAndServe(":9092", nil)
}

