package main

import (
    "fmt"
    "io/ioutil"
    "log"
    
    "github.com/go-yaml/yaml"
)

type conf struct {
    Username string `yaml:"username"` // важно указывать переменные именно с большой буквы
    Password string `yaml:"password"`
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

func main() {
    list, err := getConf("passlog.yaml")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(list)
}
