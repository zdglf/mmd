package util

import (
    "golang.org/x/text/encoding/japanese"
    "os"
    "errors"
    "log"
    "path"
    "strings"
)

type PMDBone struct {
    Name string
    //父bone的index。为-1表示没有父bone
    ParentBoneIndex int
    //子bone的index。为-1表示没有子bone
    TailPosBoneIndex int
    //bone 类型，可以是ik follow bone （4）或者是 co-rotate bone（9） 等
    Type int
    //ik follow bone 表示这个bone受到影响的ik bone index ,如果是 co-rotate bone 表示当前bone 的rotation
    IKParentBoneIndex int
    //这个bone的坐标xyz
    HeadPos []float32

}

func NewPMDBone(f *FileReader) (b *PMDBone, err error){
    b = new(PMDBone)
    if b.Name, err = f.GetStringTrim(20, japanese.ShiftJIS.NewDecoder());err !=nil{
        return
    }
    if b.ParentBoneIndex, err = f.GetInt16Little();err !=nil{
        return
    }
    if b.TailPosBoneIndex, err = f.GetInt16Little();err !=nil{
        return
    }
    if b.Type, err = f.GetUInt8Little();err !=nil{
        return
    }
    if b.IKParentBoneIndex, err = f.GetInt16Little();err !=nil{
        return
    }
    size := 3
    b.HeadPos = make([]float32, size)
    for i:=0; i < size ; i++{
        if b.HeadPos[i], err = f.GetFloatLittle();err !=nil{
            return
        }
    }
    return
}

type PMDIK struct {
    //bone 链的第一个index
    BoneIndex int
    //bone 链的最后一个index
    TargetBoneIndex int
    //bone链长度 child_bones.length
    ChainLength int
    //到达最后一个bone时的最大迭代值
    Iterations int
    //骨头与目标的的正反方向最大角度
    MaxAngle float32
    //bone index数组
    ChildBones []int
}

func NewPMDIK(f *FileReader) (ik *PMDIK, err error){
    ik = new(PMDIK)
    if ik.BoneIndex, err = f.GetInt16Little(); err!=nil{
        return
    }
    if ik.TargetBoneIndex ,err = f.GetInt16Little();err!=nil{
        return
    }
    if ik.ChainLength, err = f.GetUInt8Little(); err != nil{
        return
    }
    if ik.Iterations, err = f.GetInt16Little(); err != nil{
        return
    }
    if ik.MaxAngle, err = f.GetFloatLittle(); err != nil{
        return
    }
    ik.ChildBones = make([]int, ik.ChainLength)
    for i :=0; i< ik.ChainLength;i++{
        if ik.ChildBones[i], err = f.GetInt16Little(); err != nil{
            return
        }
    }
    return
}

type PMDJoint struct {
    //约束名
    Name string
    //第一个刚体 index
    RigidbodyA int
    //第二个刚体 index （类似于描述两个刚体的约束。）
    RigidbodyB int
    //xyz位置与第一个刚体相关
    Pos []float32
    //xyz旋转
    Rot []float32
    //xyz位置（最小值）
    ConstrainPos1 []float32
    //xyz位置（最大值）
    ConstrainPos2 []float32
    //xyz旋转（最小值）
    ConstrainRot1 []float32
    //xyz旋转（最大值）
    ConstrainRot2 []float32
    //xyz刚度（弹簧约束）
    SpringPos []float32
    //xyz旋转刚度（弹簧约束）
    SpringRot []float32
}

func NewPMDJoint(f *FileReader)(j *PMDJoint, err error)  {
    j = new(PMDJoint)
    if j.Name, err = f.GetStringTrim(20, japanese.ShiftJIS.NewDecoder());err!=nil{
        return
    }
    if j.RigidbodyA, err = f.GetUInt32Little(); err != nil{
        return
    }
    if j.RigidbodyB, err = f.GetUInt32Little(); err != nil{
        return
    }
    size := 3
    j.Pos = make([]float32, size)
    for i :=0; i< size;i++{
        if j.Pos[i], err = f.GetFloatLittle(); err != nil{
            return
        }
    }
    size = 3
    j.Rot = make([]float32, size)
    for i :=0; i< size;i++{
        if j.Rot[i], err = f.GetFloatLittle(); err != nil{
            return
        }
    }
    size = 3
    j.ConstrainPos1 = make([]float32, size)
    for i :=0; i< size;i++{
        if j.ConstrainPos1[i], err = f.GetFloatLittle(); err != nil{
            return
        }
    }
    size = 3
    j.ConstrainPos2 = make([]float32, size)
    for i :=0; i< size;i++{
        if j.ConstrainPos2[i], err = f.GetFloatLittle(); err != nil{
            return
        }
    }
    size = 3
    j.ConstrainRot1 = make([]float32, size)
    for i :=0; i< size;i++{
        if j.ConstrainRot1[i], err = f.GetFloatLittle(); err != nil{
            return
        }
    }
    size = 3
    j.ConstrainRot2 = make([]float32, size)
    for i :=0; i< size;i++{
        if j.ConstrainRot2[i], err = f.GetFloatLittle(); err != nil{
            return
        }
    }
    size = 3
    j.SpringPos = make([]float32, size)
    for i :=0; i< size;i++{
        if j.SpringPos[i], err = f.GetFloatLittle(); err != nil{
            return
        }
    }
    size = 3
    j.SpringRot = make([]float32, size)
    for i :=0; i< size;i++{
        if j.SpringRot[i], err = f.GetFloatLittle(); err != nil{
            return
        }
    }
    return
}

type PMDMaterial struct {
    //rgb
    Diffuse []float32
    Alpha float32
    Shininess float32
    //rgb
    Specular []float32
    //rgb
    Ambient []float32
    //toon材质的index [0 ,10]
    ToonIndex int
    //材质是否画上边缘线
    EdgeFlag int
    //表示有多少vertex受到这个材料影响
    FaceVertCount int;
    //材质的文件名（可以为空）
    TextureFileName string

    Textures map[string]int32
}

func NewPMDMaterial(f *FileReader) (m *PMDMaterial, err error)  {
    m = new(PMDMaterial)
    size := 3
    m.Diffuse = make([]float32, size)
    for i :=0; i< size;i++{
        if m.Diffuse[i], err = f.GetFloatLittle(); err != nil{
            return
        }
    }
    if m.Alpha, err = f.GetFloatLittle(); err != nil{
        return
    }
    if m.Shininess, err = f.GetFloatLittle(); err != nil{
        return
    }
    size = 3
    m.Specular = make([]float32, size)
    for i :=0; i< size;i++{
        if m.Specular[i], err = f.GetFloatLittle(); err != nil{
            return
        }
    }
    size = 3
    m.Ambient = make([]float32, size)
    for i :=0; i< size;i++{
        if m.Ambient[i], err = f.GetFloatLittle(); err != nil{
            return
        }
    }
    if m.ToonIndex, err = f.GetUInt8Little(); err != nil{
        return
    }
    if m.EdgeFlag, err = f.GetUInt8Little(); err != nil{
        return
    }
    if m.FaceVertCount, err = f.GetUInt32Little(); err != nil{
        return
    }
    size = 20
    if m.TextureFileName, err = f.GetStringTrim(size, japanese.ShiftJIS.NewDecoder()); err != nil{
        return
    }
    return
}

type PMDMorph struct {
    //名字
    Name string
    //类型   4种 Eyebrow（1） eye(2) lip(3), other(0)
    Count int
    //受到这个morph影响的vertices
    Type int
    //受到这个morph影响的vertices 数量
    Data []*PMDMorphData

}

type PMDMorphData struct {
    Index int
    X float32
    Y float32
    Z float32
}

func NewPMDMorph(f *FileReader)(m *PMDMorph, err error)  {
    m = new(PMDMorph)

    size := 20
    if m.Name, err = f.GetStringTrim(size, japanese.ShiftJIS.NewDecoder()); err != nil{
        return
    }
    if m.Count, err = f.GetUInt32Little(); err != nil{
        return
    }
    if m.Type, err = f.GetUInt8Little(); err != nil{
        return
    }
    m.Data = make([]*PMDMorphData, m.Count)
    size = m.Count
    for i :=0; i< size;i++{
        m.Data[i] = new(PMDMorphData)
        if m.Data[i].Index, err = f.GetUInt32Little(); err != nil{
            return
        }
        if m.Data[i].X, err = f.GetFloatLittle(); err != nil{
            return
        }
        if m.Data[i].Y, err = f.GetFloatLittle(); err != nil{
            return
        }
        if m.Data[i].Z, err = f.GetFloatLittle(); err != nil{
            return
        }
    }
    return 

}

type PMDRigidBody struct {
    //刚体的名字
    Name string
    //这个刚体影响的bone index
    RelBoneIndex int
    //碰撞组index
    GroupIndex int
    //碰撞组掩码
    GroupMask int
    //刚体形状 可以指定 sphere box oval
    ShapeType int
    //刚体的宽度（所有形状）
    ShapeW float32
    //刚体的高度 或者radius(如果类型是 sphere 或capsule)
    ShapeH float32
    //刚体的深度 (类型为box 有效)
    ShapeD float32
    //xyz与bone 相关
    Pos []float32
    //xyz旋转
    Rot []float32
    //刚体质量
    Weight float32
    //线性阻尼系数
    PosDim float32
    //角阻尼系数
    RotDim float32
    //恢复系数（反冲系数）这个应该是乳摇效果的实现。
    Recoil float32
    //摩擦系数
    Friction float32
    //类型   kinematic（只受到骨骼影响） 或simulated（身体影响骨骼，身体受到物理引擎影响）或aligned(骨骼和物理引擎影响)
    Type int

}

func NewPMDRigidBody(f *FileReader) (r *PMDRigidBody, err error)  {
    r = new(PMDRigidBody)

    size := 20
    if r.Name, err = f.GetStringTrim(size, japanese.ShiftJIS.NewDecoder()); err != nil{
        return
    }
    if r.RelBoneIndex, err = f.GetUInt16Little(); err != nil{
        return
    }

    if r.GroupIndex, err = f.GetUInt8Little(); err != nil{
        return
    }
    if r.GroupMask, err = f.GetUInt16Little(); err != nil{
        return
    }

    if r.Type, err = f.GetUInt8Little(); err != nil{
        return
    }

    if r.ShapeW, err = f.GetFloatLittle(); err != nil{
        return
    }

    if r.ShapeH, err = f.GetFloatLittle(); err != nil{
        return
    }

    if r.ShapeD, err = f.GetFloatLittle(); err != nil{
        return
    }
    size = 3
    r.Pos = make([]float32, size)
    for i :=0; i< size;i++{
        if r.Pos[i], err = f.GetFloatLittle(); err != nil{
            return
        }
    }
    size = 3
    r.Rot = make([]float32, size)
    for i :=0; i< size;i++{
        if r.Rot[i], err = f.GetFloatLittle(); err != nil{
            return
        }
    }

    if r.Weight, err = f.GetFloatLittle(); err != nil{
        return
    }
    if r.PosDim, err = f.GetFloatLittle(); err != nil{
        return
    }
    if r.RotDim, err = f.GetFloatLittle(); err != nil{
        return
    }
    if r.Recoil, err = f.GetFloatLittle(); err != nil{
        return
    }
    if r.Friction, err = f.GetFloatLittle(); err != nil{
        return
    }
    if r.Type, err = f.GetUInt8Little(); err != nil{
        return
    }
    return
}

type PMDVertex struct {
    X float32
    Y float32
    Z float32
    //normal x
    NX float32
    //normal y
    NY float32
    //normal z
    NZ float32
    //texture 的u
    U float32
    //texture 的v
    V float32
    //bone 1在bone list的index
    BoneNum1 int
    //bone 2的bone list的index
    BoneNum2 int
    //表示bone 1 的weight  bone2的weight = 100-bone1
    BoneWeight int
    //是否画上边缘线
    EdgeFlag int

}

func NewPMDVertex(f*FileReader)( v* PMDVertex, err error) {
    v = new(PMDVertex)
    if v.X, err = f.GetFloatLittle(); err != nil{
        return
    }
    if v.Y, err = f.GetFloatLittle(); err != nil{
        return
    }
    if v.Z, err = f.GetFloatLittle(); err != nil{
        return
    }
    if v.NX, err = f.GetFloatLittle(); err != nil{
        return
    }
    if v.NY, err = f.GetFloatLittle(); err != nil{
        return
    }
    if v.NZ, err = f.GetFloatLittle(); err != nil{
        return
    }
    if v.U, err = f.GetFloatLittle(); err != nil{
        return
    }
    if v.V, err = f.GetFloatLittle(); err != nil{
        return
    }
    if v.BoneNum1, err = f.GetUInt16Little(); err != nil{
        return
    }
    if v.BoneNum2, err = f.GetUInt16Little(); err != nil{
        return
    }
    if v.BoneWeight, err = f.GetUInt8Little(); err != nil{
        return
    }
    if v.EdgeFlag, err = f.GetUInt8Little(); err != nil{
        return
    }
    return
}

type PMDBoneTable struct {
    Index int
    GroupIndex int
}

type PMD struct {
    Directory string
    FileName string
    fr * FileReader

    Version float32

    Name string
    Comment string

    Vertices []*PMDVertex

    Triangles []int32

    Materials []*PMDMaterial

    Bones [] *PMDBone

    Iks [] *PMDIK

    Morphs []*PMDMorph

    MorphOrders []int

    BoneGroupNames [] string

    BoneTables []*PMDBoneTable
    EnglishFlag int

    EnglishName string

    EnglishComment string

    EnglishBoneNames []string

    EnglishMorphNames []string

    EnglishBoneGroupNames [] string

    ToonFileNames [] string

    RigidBodies []*PMDRigidBody

    Joints []*PMDJoint





}

func (p *PMD)Load(directory, fileName string)(err error)  {
    p.Directory = directory
    p.FileName = fileName
    var f *os.File
    if f, err = os.Open(path.Join(directory, fileName));err!=nil {
        return
    }
    defer f.Close()
    p.fr = new(FileReader)
    p.fr.Set(f)
    if err = p.checkHeader(); err != nil{
        return
    }

    log.Println("header ok; version", p.Version)
    if err = p.parseName(); err != nil{
        return
    }
    log.Println("name ok")
    if err = p.parseVertices(); err != nil{
        return
    }
    log.Println("vertices ok")
    if err = p.parseTriangles(); err != nil{
        return
    }
    log.Println("triangles ok")
    if err = p.parseMaterials(); err != nil{
        return
    }
    log.Println("materials ok")
    if err = p.parseBones(); err != nil{
        return
    }
    log.Println("bones ok")
    if err = p.parseIKs(); err != nil{
        return
    }
    log.Println("iks ok")
    if err = p.parseMorphs(); err != nil{
        return
    }
    log.Println("morphs ok")
    if err = p.parseMorphOrder(); err != nil{
        return
    }
    log.Println("morphorder ok")
    if err = p.parseBoneGroupName(); err != nil{
        return
    }
    log.Println("bonegroup name ok")
    if err = p.parseBoneTable(); err != nil{
        return
    }
    log.Println("bone table ok")
    if err = p.parseEnglishFlag(); err != nil{
        return
    }
    log.Println("english flag ok")
    if p.EnglishFlag != 0{
        if err = p.parseEnglishName(); err != nil{
            return
        }
        log.Println("englishName ok")
        if err = p.parseEnglishBoneNames(); err != nil{

        }
        log.Println("englishBoneName ok")
        if err = p.parseEnglishMorphNames(); err != nil{

        }
        log.Println("englishmorphName ok")
        if err = p.parseEnglishBoneGroupNames(); err != nil{
            return
        }
        log.Println("englishbone group Name ok")
    }
    if err = p.parseToonFileNames(); err != nil{
        return
    }
    log.Println("toon file name ok")
    if err = p.parseRigidBodys(); err != nil{
        return
    }
    log.Println("rigid body ok")
    if err = p.parseJoints(); err != nil{
        return
    }
    log.Println("joints ok")



    return
}
func (p *PMD)checkHeader() (err error){
    const PMDHEADER string  = "Pmd"
    var header string
    if header, err = p.fr.GetString(3, japanese.ShiftJIS.NewDecoder());err != nil{
        return
    }
    if header != PMDHEADER{
        err = errors.New("Pmd Header check fail")
        return
    }
    if p.Version, err = p.fr.GetFloatLittle();err !=nil{
        return
    }
    if p.Version != 1.0{
        err = errors.New("Version is  not 1.0")
        return
    }
    return

}

func (p *PMD)parseName() (err error){
    if p.Name, err = p.fr.GetStringTrim(20, japanese.ShiftJIS.NewDecoder());err != nil{
        return
    }
    if p.Comment, err = p.fr.GetStringTrim(256, japanese.ShiftJIS.NewDecoder());err != nil{
        return
    }
    return
}

func (p *PMD)parseVertices() (err error){
    var count int
    if count, err = p.fr.GetUInt32Little();err !=nil{
        return
    }
    p.Vertices = make([]*PMDVertex, count)
    for i:= 0; i < count;i++{
        if p.Vertices[i], err = NewPMDVertex(p.fr);err != nil{
            return
        }
    }
    return
}

func (p *PMD)parseTriangles() (err error){
    var count int
    if count, err = p.fr.GetUInt32Little();err !=nil{
        return
    }
    p.Triangles = make([]int32, count)
    var data int
    for i:= 0; i < count;i+=3{
        if data, err = p.fr.GetUInt16Little();err != nil{
            return
        }
        p.Triangles[i+1]=int32(data)
        if data, err = p.fr.GetUInt16Little();err != nil{
            return
        }
        p.Triangles[i]=int32(data)
        if data, err = p.fr.GetUInt16Little();err != nil{
            return
        }
        p.Triangles[i+2]=int32(data)
    }
    return
}

func (p *PMD)parseMaterials() (err error){
    var count int
    if count, err = p.fr.GetUInt32Little();err !=nil{
        return
    }
    p.Materials = make([]*PMDMaterial, count)
    for i:= 0; i < count;i++{
        if p.Materials[i], err = NewPMDMaterial(p.fr);err != nil{
            return
        }
    }
    return
}

func (p *PMD)parseBones() (err error){
    var count int
    if count, err = p.fr.GetUInt16Little();err !=nil{
        return
    }
    p.Bones = make([]*PMDBone, count)
    for i:= 0; i < count;i++{
        if p.Bones[i], err = NewPMDBone(p.fr);err != nil{
            return
        }
    }
    return
}


func (p *PMD)parseIKs() (err error){
    var count int
    if count, err = p.fr.GetUInt16Little();err !=nil{
        return
    }
    p.Iks = make([]*PMDIK, count)
    for i:= 0; i < count;i++{
        if p.Iks[i], err = NewPMDIK(p.fr);err != nil{
            return
        }
    }
    return
}


func (p *PMD)parseMorphs() (err error){
    var count int
    if count, err = p.fr.GetUInt16Little();err !=nil{
        return
    }
    p.Morphs = make([]*PMDMorph, count)
    for i:= 0; i < count;i++{
        if p.Morphs[i], err = NewPMDMorph(p.fr);err != nil{
            return
        }
    }
    return
}

func (p *PMD)parseMorphOrder() (err error){
    var count int
    if count, err = p.fr.GetUInt8Little();err !=nil{
        return
    }
    p.MorphOrders = make([]int, count)
    for i:= 0; i < count;i++{
        if p.MorphOrders[i], err = p.fr.GetUInt16Little();err != nil{
            return
        }
    }
    return
}


func (p *PMD)parseBoneGroupName() (err error){
    var count int
    if count, err = p.fr.GetUInt8Little();err !=nil{
        return
    }
    p.BoneGroupNames = make([]string, count)
    for i:= 0; i < count;i++{
        if p.BoneGroupNames[i], err = p.fr.GetStringTrim(50, japanese.ShiftJIS.NewDecoder());err != nil{
            return
        }
    }
    return
}


func (p *PMD)parseBoneTable() (err error){
    var count int
    if count, err = p.fr.GetUInt32Little();err !=nil{
        return
    }
    p.BoneTables = make([]*PMDBoneTable, count)
    for i:= 0; i < count;i++{
        p.BoneTables[i] = new(PMDBoneTable)
        if p.BoneTables[i].Index, err = p.fr.GetInt16Little();err != nil{
            return
        }
        if p.BoneTables[i].GroupIndex, err = p.fr.GetUInt8Little();err != nil{
            return
        }
    }
    return
}


func (p *PMD)parseEnglishFlag() (err error){
    if p.EnglishFlag, err = p.fr.GetUInt8Little();err !=nil{
        return
    }
    return
}


func (p *PMD)parseEnglishName() (err error){
    if p.EnglishName, err = p.fr.GetStringTrim(20, japanese.ShiftJIS.NewDecoder());err != nil{
        return
    }
    if p.EnglishComment, err = p.fr.GetStringTrim(256, japanese.ShiftJIS.NewDecoder());err != nil{
        return
    }
    return
}


func (p *PMD)parseEnglishBoneNames() (err error){
    count := len(p.Bones)
    p.EnglishBoneNames = make([]string, count)
    for i:= 0; i < count;i++{
        if p.EnglishBoneNames[i], err = p.fr.GetStringTrim(20, japanese.ShiftJIS.NewDecoder());err != nil{
            return
        }
    }
    return
}


func (p *PMD)parseEnglishMorphNames() (err error){
    count := len(p.Morphs)-1
    if(count<0){
        count = 0
    }
    p.EnglishMorphNames = make([]string, count)
    for i:= 0; i < count;i++{
        if p.EnglishMorphNames[i], err = p.fr.GetStringTrim(20, japanese.ShiftJIS.NewDecoder());err != nil{
            return
        }
    }
    return
}

func (p *PMD)parseEnglishBoneGroupNames() (err error){
    count := len(p.BoneGroupNames)
    p.EnglishBoneGroupNames = make([]string, count)
    for i:= 0; i < count;i++{
        if p.EnglishBoneGroupNames[i], err = p.fr.GetStringTrim(50, japanese.ShiftJIS.NewDecoder());err != nil{
            return
        }
    }
    return
}


func (p *PMD)parseToonFileNames() (err error){

    count := 10
    p.ToonFileNames = make([] string, count)
    for i:= 0; i < count;i++{
        if p.ToonFileNames[i], err = p.fr.GetStringTrim(100, japanese.ShiftJIS.NewDecoder());err != nil{
            return
        }
        p.ToonFileNames[i] = strings.Replace(p.ToonFileNames[i], "\\", "/", -1)
    }
    return

}

func (p *PMD)parseRigidBodys() (err error){
    var count int
    if count, err = p.fr.GetUInt32Little();err !=nil{
        return
    }
    p.RigidBodies = make([]*PMDRigidBody, count)
    for i:= 0; i < count;i++{
        if p.RigidBodies[i], err = NewPMDRigidBody(p.fr);err != nil{
            return
        }
    }
    return
}


func (p *PMD)parseJoints() (err error){
    var count int
    if count, err = p.fr.GetUInt32Little();err !=nil{
        return
    }
    p.Joints = make([]*PMDJoint, count)
    for i:= 0; i < count;i++{
        if p.Joints[i], err = NewPMDJoint(p.fr);err != nil{
            return
        }
    }
    return
}

