package util

import (
    "os"
    "golang.org/x/text/encoding/japanese"
    "errors"
    "log"
)

const (
    VMD_V1 = "Vocaloid Motion Data file"
    VMD_V2 = "Vocaloid Motion Data 0002"
)

type VMDBoneKeyFrame struct {
    Name string
    Frame int
    Location []float32
    Rotation []float32
    Interpolation []int
}

func NewVMDBoneKeyFrame(f *FileReader)(b *VMDBoneKeyFrame, err error) {
    b = new(VMDBoneKeyFrame)
    if b.Name, err = f.GetStringTrim(15, japanese.ShiftJIS.NewDecoder());err!=nil{
        return
    }
    if b.Frame, err = f.GetUInt32Little();err!=nil{
        return
    }
    size := 3
    b.Location = make([]float32, size)
    for i :=0; i< size;i++{
        if b.Location[i], err = f.GetFloatLittle(); err != nil{
            return
        }
    }
    size = 4
    b.Rotation = make([]float32, size)
    for i :=0; i< size;i++{
        if b.Rotation[i], err = f.GetFloatLittle(); err != nil{
            return
        }
    }
    size = 64
    b.Interpolation = make([]int, size)
    for i :=0; i< size;i++{
        if b.Interpolation[i], err = f.GetUInt8Little(); err != nil{
            return
        }
    }
    return
}



type VMDMorphKeyFrame struct {
    Name string
    Frame int
    Weight float32
}


func NewVMDMorphKeyFrame(f *FileReader)(m *VMDMorphKeyFrame, err error) {
    m = new(VMDMorphKeyFrame)
    if m.Name, err = f.GetStringTrim(15, japanese.ShiftJIS.NewDecoder());err!=nil{
        return
    }
    if m.Frame, err = f.GetUInt32Little();err!=nil{
        return
    }

    if m.Weight, err = f.GetFloatLittle();err!=nil{
        return
    }
    return
}

type VMDCameraKeyFrame struct {

    Frame int
    Distance float32
    Location []float32
    Rotation []float32
    Interpolation []int
    ViewAngle int
    NoPerspective int
}


func NewVMDCameraKeyFrame(f *FileReader)(c *VMDCameraKeyFrame, err error) {
    c = new(VMDCameraKeyFrame)

    if c.Frame, err = f.GetUInt32Little();err!=nil{
        return
    }

    if c.Distance, err = f.GetFloatLittle();err!=nil{
        return
    }
    size := 3
    c.Location = make([]float32, size)
    for i :=0; i< size;i++{
        if c.Location[i], err = f.GetFloatLittle(); err != nil{
            return
        }
    }
    size = 3
    c.Rotation = make([]float32, size)
    for i :=0; i< size;i++{
        if c.Rotation[i], err = f.GetFloatLittle(); err != nil{
            return
        }
    }
    size = 24
    c.Interpolation = make([]int, size)
    for i :=0; i< size;i++{
        if c.Interpolation[i], err = f.GetUInt8Little(); err != nil{
            return
        }
    }
    if c.ViewAngle, err = f.GetUInt32Little();err!=nil{
        return
    }
    if c.NoPerspective, err = f.GetUInt8Little(); err != nil{
        return
    }
    return
}

type VMDLightKeyFrame struct {
    Frame int
    Color []float32
    Location []float32
}


func NewVMDLightKeyFrame(f *FileReader)(l *VMDLightKeyFrame, err error) {
    l = new(VMDLightKeyFrame)
    if l.Frame, err = f.GetUInt32Little();err!=nil{
        return
    }
    size := 3
    l.Color = make([]float32, size)
    for i :=0; i< size;i++{
        if l.Color[i], err = f.GetFloatLittle(); err != nil{
            return
        }
    }
    size = 3
    l.Location = make([]float32, size)
    for i :=0; i< size;i++{
        if l.Location[i], err = f.GetFloatLittle(); err != nil{
            return
        }
    }
    return
}

type VMDSelfShadowKeyFrame struct {
    Frame int
    Mode int
    Distance float32
}


func NewVMDSelfShadowKeyFrame(f *FileReader)(ss *VMDSelfShadowKeyFrame, err error) {
    ss = new(VMDSelfShadowKeyFrame)
    if ss.Frame, err = f.GetUInt32Little();err!=nil{
        return
    }
    if ss.Mode, err = f.GetUInt8Little();err!=nil{
        return
    }
    if ss.Distance, err = f.GetFloatLittle();err!=nil{
        return
    }
    return
}

type VMD struct {
    header string
    ModelName string

    BoneFrames []*VMDBoneKeyFrame
    LightFrames []*VMDLightKeyFrame
    SelfShadowFrames []*VMDSelfShadowKeyFrame
    CameraFrames []*VMDCameraKeyFrame
    MorphFrames []*VMDMorphKeyFrame

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
    log.Println("header ok; version: ", v.header)
    if err = v.parseBoneFrame(); err != nil{
        return
    }
    log.Println("bone frame ok")
    if err = v.parseMorphFrame(); err != nil{
        return
    }
    log.Println("morph frame ok")
    if err = v.parseCameraFrame(); err != nil{
        return
    }
    log.Println("camera frame ok")
    if err = v.parseLightFrame(); err != nil{
        return
    }
    log.Println("light frame ok")
    if err = v.parseSelfShadowFrame(); err != nil{
        return
    }
    log.Println("selfshadow frame ok")
    return

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

func (v *VMD)parseBoneFrame()(err error) {
    var count int
    if count ,err = v.fr.GetUInt32Little(); err != nil{
        return
    }
    v.BoneFrames = make([]*VMDBoneKeyFrame, count)
    for i := 0; i < count; i++{
        if v.BoneFrames[i], err = NewVMDBoneKeyFrame(v.fr);err != nil{
            return
        }
    }
    return
}


func (v *VMD)parseMorphFrame()(err error) {
    var count int
    if count ,err = v.fr.GetUInt32Little(); err != nil{
        return
    }
    v.MorphFrames = make([]*VMDMorphKeyFrame, count)
    for i := 0; i < count; i++{
        if v.MorphFrames[i], err = NewVMDMorphKeyFrame(v.fr);err != nil{
            return
        }
    }
    return
}


func (v *VMD)parseCameraFrame()(err error) {
    var count int
    if count ,err = v.fr.GetUInt32Little(); err != nil{
        return
    }
    v.CameraFrames = make([]*VMDCameraKeyFrame, count)
    for i := 0; i < count; i++{
        if v.CameraFrames[i], err = NewVMDCameraKeyFrame(v.fr);err != nil{
            return
        }
    }
    return
}


func (v *VMD)parseLightFrame()(err error) {
    var count int
    if count ,err = v.fr.GetUInt32Little(); err != nil{
        return
    }
    v.LightFrames = make([]*VMDLightKeyFrame, count)
    for i := 0; i < count; i++{
        if v.LightFrames[i], err = NewVMDLightKeyFrame(v.fr);err != nil{
            return
        }
    }
    return
}


func (v *VMD)parseSelfShadowFrame()(err error) {
    var count int
    if count ,err = v.fr.GetUInt32Little(); err != nil{
        return
    }
    v.SelfShadowFrames = make([]*VMDSelfShadowKeyFrame, count)
    for i := 0; i < count; i++{
        if v.SelfShadowFrames[i], err = NewVMDSelfShadowKeyFrame(v.fr);err != nil{
            return
        }
    }
    return
}

