package util

import (
    "testing"
    "os"
    "log"
)

func TestPMD_Load(t *testing.T) {



    p := new(PMD)
    path, _ := os.Getwd()
    log.Println(path)
    err := p.Load(path+"/test_data", "alice.pmd")
    if err != nil{
        t.Error(err)
    }
    log.Println(p.Name)
    log.Println(p.Comment)

}
