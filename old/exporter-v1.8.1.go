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
    Port string `yaml:"port"`
}


type slaveStatus struct {
    Slave_SQL_Running string		// yes:no
    Position string					// 1:0
    Slave_SQL_Running_State string	// Slave has read all relay log; waiting for more updates
    Last_SQL_Error string			// errno + error

    Slave_IO_Running string			// yes:no
    Slave_IO_State string			// Waiting for master to send event
    Last_IO_Error string			// errno + error

    Last_Error string				// errno + error
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


func MySQLSlaveStatus() (map[string]string, error) {
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
        replicationStatuses, err := MySQLSlaveStatus()
        if err != nil {
            fmt.Println(err)
        }
        
        // выделим отдельно те статусы, которые мы обрабатывать пока не будем
        ignorStatuses := [17]string{"Master_Host", "Master_User", "Master_SSL_Allowed", "Executed_Gtid_Set", "Slave_IO_State", "Master_Log_File", "Slave_IO_Running", "Master_UUID", "Until_Condition", "Retrieved_Gtid_Set", "Relay_Master_Log_File", "Relay_Log_File", "Master_Info_File", "Slave_SQL_Running_State", "Replicate_Wild_Ignore_Table", "Master_SSL_Verify_Server_Cert", "Slave_SQL_Running"}
        
        // эта переменная передает статус репликации (ошибки и states), когда репликация выключена
        var slaveSQLRunningState = slaveStatus{ replicationStatuses["Slave_SQL_Running"], "0", replicationStatuses["Slave_SQL_Running_State"], replicationStatuses["Last_SQL_Error"], replicationStatuses["Slave_IO_Running"], replicationStatuses["Slave_IO_State"], replicationStatuses["Last_IO_Error"], replicationStatuses["Last_Error"] }
        
        // создаем триггер для grafana
        if slaveSQLRunningState.Slave_SQL_Running == "Yes"{
        	slaveSQLRunningState.Position = "1"
        } else {
        	slaveSQLRunningState.Position = "0"
        }
        
        // вырезаем все, что не будем обрабатывать
        for _, ignor := range ignorStatuses {
        	delete(replicationStatuses, ignor)
        }
        
        
        // переменная под текс ответа для prometheus
        var responseText string
        responseText += fmt.Sprintf("# MySQL replication slave status exporter for Prometheus&Grafana\n")
        for status, value := range replicationStatuses {
        	responseText += fmt.Sprintf("mysql_%s{method=\"slave\"} %s\n", status, value)
        }
        
        // добавляем сообщения
        responseText += fmt.Sprintf("mysql_Slave_SQL_Running{")
        responseText += fmt.Sprintf("%s=%s, ", "Slave_SQL_Running", slaveSQLRunningState.Slave_SQL_Running)
        responseText += fmt.Sprintf("%s=%s, ", "Slave_SQL_Running_State", slaveSQLRunningState.Slave_SQL_Running_State)
        responseText += fmt.Sprintf("%s=%s, ", "Last_SQL_Error", slaveSQLRunningState.Last_SQL_Error)
        responseText += fmt.Sprintf("%s=%s, ", "Slave_IO_Running", slaveSQLRunningState.Slave_IO_Running) 
        responseText += fmt.Sprintf("%s=%s, ", "Slave_IO_State", slaveSQLRunningState.Slave_IO_State)
        responseText += fmt.Sprintf("%s=%s, ", "Last_IO_Error", slaveSQLRunningState.Last_IO_Error)
        responseText += fmt.Sprintf("%s=%s, ", "Last_Error", slaveSQLRunningState.Last_Error)
        responseText += fmt.Sprintf("} %s\n", slaveSQLRunningState.Position)

        // отдаем текст по запросу
        fmt.Fprintf(w, responseText)
    }
    
    // лезем в конфиг за портом
    list, err := getConf("conf.yaml")
    if err != nil {
        fmt.Println(err)
    }
    
    // http сервер
    // в корень / поместим index.html (из ./static), в /metricsm - метрики
    http.Handle("/", http.FileServer(http.Dir("/etc/replication")))
    http.HandleFunc("/metrics", message)
    http.ListenAndServe(":" + list.Port, nil)
}
