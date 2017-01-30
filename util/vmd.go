package util

import "os"

const (
    VMD_V1 = "Vocaloid Motion Data file"
    VMD_V2 = "Vocaloid Motion Data 0002"
)

type VMD struct {
    header string
    fr *FileReader
}

func (v *VMD)Load(filePath string)(err error){
    var f *os.File
    if f, err = os.Open(filePath);err!=nil {
        return
    }
    defer f.Close()
    v.fr = new(FileReader)
    v.fr.Set(f)
    if err = v.checkHeader(); err != nil{
        return
    }
    return nil

}

func (v *VMD)checkHeader()(err error) {
    return nil
}