package main

import (
    "database/sql"
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"

    _ "github.com/go-sql-driver/mysql"
    "github.com/go-yaml/yaml"
)


type conf struct {
    Username string `yaml:"username"` // важно указывать переменные именно с большой буквы
    Password string `yaml:"password"`
    Database string `yaml:"database"`
}


func getConf(filename string) (*conf, error) {
    
    yamlFile, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    
    configs := &conf{}
    err = yaml.Unmarshal(yamlFile, configs)
    if err != nil {
        return nil, fmt.Errorf("in file %q: %v", filename, err)
    }
    
    return configs, nil
}


func MySQLSlaveStatus(password string) (map[string]string, error) {
	// достаем конфиг, читаем данные для авторизации
	authorization, err := getConf("conf.yaml")
    if err != nil {
        fmt.Println(err)
    }
    
    // авторизуемся в MySQL DB под replication (password: Yes)
    db, err := sql.Open("mysql", fmt.Sprintf(authorization.Username + ":%s@/%s", authorization.Password, authorization.Database))
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
    
    columns, err := rows.Columns()
    if err != nil {
        fmt.Println(err)
    }
    
    // создаем пустой массив values длины len(columns)
    values := make([]interface{}, len(columns))

    // для каждого ключа, который мы находим при обходе массива columns,
    // устанавливаем в массиве values место пустым типом sql.RawBytes (аналогичный []byte)
    for key, _ := range columns {
        values[key] = new(sql.RawBytes)
    }
    
    // "rows.Next()"" рекурсивная функция
    for rows.Next() {
        //"values..." сообщает Golang использовать каждый непустой слот
        err = rows.Scan(values...)
        if err != nil {
            fmt.Println(err)
        }
    }

    for index, columnName := range columns {
        // преобразуем sql.RawBytes в строки при помощи "fmt.Sprintf"
        columnValue := fmt.Sprintf("%s", values[index])

        // удаляем "&" из переменных слоя
        columnValue = strings.Replace(columnValue, "&", "", -1)
        if len(columnValue) == 0 {
            continue
        }

        selectStatuses[columnName] = columnValue
    }

    return selectStatuses, nil
}

func main() {
    message := func(w http.ResponseWriter, r *http.Request) {
        // лезем в БД
        replicationStatuses, err := MySQLSlaveStatus("strong_passwd")
        if err != nil {
            fmt.Println(err)
        }
        // переменная под текс ответа
        var responseText string // переменная под текс ответа для prometheus
        responseText += fmt.Sprintf("MAIN mysql_custom_message\n")
        for status, value := range replicationStatuses {
            responseText += fmt.Sprintf("mysql_custom_message{method=\"slave\", status=\"%s\"} %s\n", status, value)
        }

        // отдаем текст по запросу
        fmt.Fprintf(w, responseText)
    }
    
    // http сервер
    // в корень / поместим index.html (из ./static), в /metricsm - метрики
    http.Handle("/", http.FileServer(http.Dir("static")))
    http.HandleFunc("/metrics", message)
    http.ListenAndServe(":9092", nil)
}
