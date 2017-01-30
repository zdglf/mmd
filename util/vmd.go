package util

import (
    "os"
    "golang.org/x/text/encoding/japanese"
    "errors"
)

const (
    VMD_V1 = "Vocaloid Motion Data file"
    VMD_V2 = "Vocaloid Motion Data 0002"
)

type BoneKeyFrame struct {
    Name string
    Frame int
    Location []float32
    Rotation []float32
    Interpolation []int
}

type MorphKeyFrame struct {
    Name string
    Frame int
    Weight float32
}

type CameraKeyFrame struct {

    Frame int
    Distance float32
    Location []float32
    Rotation []float32
    Interpolation []int
    ViewAngle int
    NoPerspective int
}

type LightKeyFrame struct {
    Frame int
    Color []float32
    Location []float32
}

type SelfShadowKeyFrame struct {
    Frame int
    Mode int
    Distance float32
}

type VMD struct {
    header string
    ModelName string

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
    if v.header, err = v.fr.GetStringUTF8Trim(30); err != nil{
        return
    }
    switch v.header {
    case VMD_V2:
        if v.ModelName, err = v.fr.GetStringTrim(20, japanese.ShiftJIS.NewDecoder()); err != nil{
            return
        }
    default:
        err = errors.New("vmd version error; version: "+ v.header)
        return
    }


    return nil
}