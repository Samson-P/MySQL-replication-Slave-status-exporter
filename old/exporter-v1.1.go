package main

import (
    "database/sql"
    "fmt"
    "net/http"

    _ "github.com/go-sql-driver/mysql"
)


func UserSelectCounts(password string) (map[string]int64, int64, error) {
    db, err := sql.Open("mysql", fmt.Sprintf("replication:%s@/%s", password, "cnc_factory"))
    if err!= nil {
        panic(err.Error())
    }
    defer db.Close()

    var sqlSelectUsername string
    var sqlSelectCount, sqlSelectCounts int64
    selectUsers := make(map[string]int64)
    rows, err := db.Query("SHOW SLAVE STATUS")
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
    // обычный http сервер
    message := func(w http.ResponseWriter, r *http.Request) {
        // лезем в базу за нужными данными
        users, userCounts, err := UserSelectCounts("tSj9ERz@5*DvSos?")
        if err != nil {
            fmt.Println(err)
        }

        var responseText string // переменная под текс ответа для prometheus
        responseText += fmt.Sprintf("mysql_custom_message_bmgeek{method=\"all\", base=\"%s\"} %d\n", userCounts)
        for user, usCount := range users {
            responseText += fmt.Sprintf("mysql_custom_message_dmgeek{method=\"user\", user=\"%s\", base=\"%s\"} %d\n", user, usCount)
        }
        // отдаем текст по запросу
        fmt.Fprintf(w, responseText)
    }

    http.HandleFunc("/metrics", message)
    http.ListenAndServe(":9092", nil)
}
