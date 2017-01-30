package util

import (
    "testing"
    "os"
    "log"
)

func TestPMX_Load(t *testing.T) {
    p := new(PMX)
    path, _ := os.Getwd()
    log.Println(path)
    err := p.Load(path+"/test_data", "test.pmx")
    if err != nil{
        t.Error(err)
    }
    log.Println(p.Name)
    log.Println(p.Comment)
    log.Println(p.EnglishName)
    //log.Println(p.EnglishComment)
}