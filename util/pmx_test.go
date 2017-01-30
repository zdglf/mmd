package util

import (
    "testing"
    "os"
    "log"
)

func TestPMX_Load(t *testing.T) {
    p := new(PMX)
    path, _ := os.Getwd()
    err := p.Load(path+"/test_data", "test.pmx")
    if err != nil{
        t.Error(err)
    }
    if p.Name == "銀獅式初音ミクV3_C_ver1.10"{
        log.Println(p.Name)
    }else{
        t.Error(p.Name, "is not 銀獅式初音ミクV3_C_ver1.10")
    }
    //log.Println(p.Name)
    //log.Println(p.Comment)
    //log.Println(p.EnglishName)
    //log.Println(p.EnglishComment)
}