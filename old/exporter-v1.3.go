package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "strings"

    _ "github.com/go-sql-driver/mysql"
)


func MySQLSlaveStatus(password string) (map[string]string, error) {
    db, err := sql.Open("mysql", fmt.Sprintf("replication:%s@/%s", password, "cnc_factory"))
    if err!= nil {
        panic(err.Error())
    }
    defer db.Close()
    
    var myQuery string
    myQuery = "SHOW SLAVE STATUS"
    
    selectStatuses := make(map[string]string)
    rows, err := db.Query(myQuery)
    if err != nil {
        return selectStatuses, err
    }
    
    // sql.Rows has a function which returns all column names
    // as a slice of []string. Variable "columns" represents this
    columns, err := rows.Columns()
    if err != nil {
        log.Fatal(err)
    }
    // fmt.Println(len(columns)) // debug =)

    // variable "values" is a pre-populated array of empty interfaces
    // We load an empty interface for every column 'sql.Rows' has.
    // The interfaces will allow us to call methods of any type that replaces it
    values := make([]interface{}, len(columns))

    // for every key we find while traversing array "columns"
    // set the corresponding interface in array "values" to be populated
    // with an empty sql.RawBytes type
    // sql.RawBytes is analogous to []byte
    for key, _ := range columns {
        values[key] = new(sql.RawBytes)
    }

    //Contrary to appearances, this is not a loop through every row
    // "rows.Next()"" is a recursive function that is called immediately
    // upon every row until we hit "rows.Next == false"
    // This is important because it means you must prepopulate variables or
    // arrays to the exact number of columns in the target SQL table
    // more details at: https://golang.org/pkg/database/sql/#Rows.Next
    for rows.Next() {
        //the "values..." tells Go to use every available slot to populate data
        err = rows.Scan(values...)
        if err != nil {
            log.Fatal(err)
        }
    }

    for index, columnName := range columns {
        // convert sql.RawBytes to String using "fmt.Sprintf"
        columnValue := fmt.Sprintf("%s", values[index])

        // Remove "&" from row values
        columnValue = strings.Replace(columnValue, "&", "", -1)

        // Optional: Don't display values that are NULL
        // Remove "if" to return empty NULL values
        if len(columnValue) == 0 {
            continue
        }

        // Print finished product
        // fmt.Printf("%s: %s\n", columnName, columnValue)
        selectStatuses[columnName] = columnValue
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
        responseText += fmt.Sprintf("MAIN mysql_custom_message\n")
        for status, value := range replicationStatuses {
            responseText += fmt.Sprintf("mysql_custom_message{method=\"slave\", status=\"%s\"} %s\n", status, value)
        }
        
        // отдаем текст по запросу
        fmt.Fprintf(w, responseText)
    }

    http.Handle("/", http.FileServer(http.Dir("static")))
    http.HandleFunc("/metrics", message)
    http.ListenAndServe(":9092", nil)
}
