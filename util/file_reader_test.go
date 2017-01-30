package util

import (
    "testing"
    "os"
    "log"
    "golang.org/x/text/encoding/japanese"
)

func TestFileReader_GetFloatBig(t *testing.T) {
    fr := new(FileReader)
    f, err := os.Open("test_data/test")
    if err!=nil{
        t.Error(err)
        return
    }
    defer f.Close()
    fr.Set(f)
    if str, err := fr.GetString(3, japanese.ShiftJIS.NewDecoder()); err == nil{
        if str != "Pmd"{
            t.Error("str should be Pmd but ", str)
        }
    }else{
        t.Error(err)
    }
    if version, err := fr.GetFloatLittle(); err == nil{
        if(version!=1.0){
            t.Error("version should be 1")
        }
    }else{
        t.Error(err)
    }
    if name, err := fr.GetString(20, japanese.ShiftJIS.NewDecoder()); err == nil{
        log.Println(name)
    }else{
        t.Error(err)
    }
    if comment, err := fr.GetString(256, japanese.ShiftJIS.NewDecoder()); err == nil{
        log.Println(comment)
    }else{
        t.Error(err)
    }



}
