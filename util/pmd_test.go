package util

import (
    "testing"
    "os"
    "log"
)

func TestPMD_Load(t *testing.T) {



    p := new(PMD)
    path, _ := os.Getwd()
    err := p.Load(path+"/test_data", "alice.pmd")
    if err != nil{
        t.Error(err)
    }
    if p.Name == "門を開く者 アリス"{
        log.Println(p.Name)
    }else{
        t.Error(p.Name, "is not 門を開く者 アリス")
    }

    //log.Println(p.Comment)

}
