package util

import (
    "golang.org/x/text/encoding"
    "os"
    "errors"
    "golang.org/x/text/encoding/unicode"
    "log"
    "path"
)

const(
    VERTEX_BDEF1 int = iota
    VERTEX_BDEF2
    VERTEX_BDEF4
    VERTEX_SDEF
    VERTEX_QDEF
)

const(
    DISPLAYFRAME_BONE int = iota
    DISPLAYFRAME_MORPH
)

type PMXVertex struct {


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

    //count*4;
    AddVec4 [] float32
    //0 = BDEF1, 1 = BDEF2, 2 = BDEF4, 3 = SDEF, 4 = QDEF
    WeightType int
    //0 = BDEF1
    BoneIndex1 int//Weight = 1.0

    //1 = BDEF2
    BoneIndex2 int
    Bone1Weight float32//Bone 2 weight = 1.0 - Bone 1 weight

    //2 = BDEF4 4 = QDEF
    BoneIndex3 int
    BoneIndex4 int
    Bone2Weight float32
    Bone3Weight float32
    Bone4Weight float32

    //3 = SDEF
    C []float32
    R0 []float32
    R1 []float32

    Scale float32
}

func NewPMXVertex(f *FileReader, sizes []int, decoder *encoding.Decoder)(v *PMXVertex, err error)  {
    addCount := sizes[1]
    boneIndexSize := sizes[5]
    v = new(PMXVertex)
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
    count := addCount*4
    v.AddVec4 = make([]float32, count)
    for i :=0; i< count;i++{
        if v.AddVec4[i], err = f.GetFloatLittle(); err != nil{
            return
        }
    }
    if v.WeightType, err = f.GetInt8Little(); err != nil{
        return
    }
    switch v.WeightType {
    case VERTEX_BDEF1:
        if v.BoneIndex1, err = f.GetIntLittle(boneIndexSize); err != nil{
            return
        }
    case VERTEX_BDEF2:
        if v.BoneIndex1, err = f.GetIntLittle(boneIndexSize); err != nil{
            return
        }
        if v.BoneIndex2, err = f.GetIntLittle(boneIndexSize); err != nil{
            return
        }
        if v.Bone1Weight, err = f.GetFloatLittle(); err != nil{
            return
        }
    case VERTEX_BDEF4, VERTEX_QDEF:
        if v.BoneIndex1, err = f.GetIntLittle(boneIndexSize); err != nil{
            return
        }
        if v.BoneIndex2, err = f.GetIntLittle(boneIndexSize); err != nil{
            return
        }
        if v.BoneIndex3, err = f.GetIntLittle(boneIndexSize); err != nil{
            return
        }
        if v.BoneIndex4, err = f.GetIntLittle(boneIndexSize); err != nil{
            return
        }
        if v.Bone1Weight, err = f.GetFloatLittle(); err != nil{
            return
        }
        if v.Bone2Weight, err = f.GetFloatLittle(); err != nil{
            return
        }
        if v.Bone3Weight, err = f.GetFloatLittle(); err != nil{
            return
        }
        if v.Bone4Weight, err = f.GetFloatLittle(); err != nil{
            return
        }
    case VERTEX_SDEF:
        if v.BoneIndex1, err = f.GetIntLittle(boneIndexSize); err != nil{
            return
        }
        if v.BoneIndex2, err = f.GetIntLittle(boneIndexSize); err != nil{
            return
        }
        if v.Bone1Weight, err = f.GetFloatLittle(); err != nil{
            return
        }
        count = 3
        v.C = make([]float32, count)
        for i :=0; i< count;i++{
            if v.C[i], err = f.GetFloatLittle(); err != nil{
                return
            }
        }
        count = 3
        v.R0 = make([]float32, count)
        for i :=0; i< count;i++{
            if v.R0[i], err = f.GetFloatLittle(); err != nil{
                return
            }
        }
        count = 3
        v.R1 = make([]float32, count)
        for i :=0; i< count;i++{
            if v.R1[i], err = f.GetFloatLittle(); err != nil{
                return
            }
        }
    }
    if v.Scale, err = f.GetFloatLittle(); err != nil{
        return
    }
    return
}

type PMXIK struct {
    BoneIndex int

    HasLimit int//When equal to 1, use angle limits
    MinLimit []float32
    MaxLimit []float32
}

func NewPMXIK(f *FileReader, sizes []int)(ik *PMXIK, err error) {
    boneIndexSize := sizes[5]
    ik = new(PMXIK)
    if ik.BoneIndex, err = f.GetIntLittle(boneIndexSize); err != nil{
        return
    }
    if ik.HasLimit, err = f.GetInt8Little(); err != nil{
        return
    }
    if ik.HasLimit == 1{
        count := 3
        ik.MinLimit = make([]float32, count)
        for i :=0; i< count;i++{
            if ik.MinLimit[i], err = f.GetFloatLittle(); err != nil{
                return
            }
        }

        count = 3
        ik.MaxLimit = make([]float32, count)
        for i :=0; i< count;i++{
            if ik.MaxLimit[i], err = f.GetFloatLittle(); err != nil{
                return
            }
        }
    }
    return
}

type PMXDisplayFrame struct {
    Type int//0 == Bone 1 == Morph;
    Index int
}

func NewPMXDisplayFrame(f *FileReader, sizes []int)(df *PMXDisplayFrame, err error) {
    boneIndexSize := sizes[5]
    morphIndexSize := sizes[6]
    df = new(PMXDisplayFrame)
    if df.Type, err = f.GetInt8Little(); err != nil{
        return
    }
    switch df.Type {
    case DISPLAYFRAME_BONE:
        if df.Index, err = f.GetIntLittle(boneIndexSize); err != nil{
            return
        }
    case DISPLAYFRAME_MORPH:
        if df.Index, err = f.GetIntLittle(morphIndexSize); err != nil{
            return
        }
    }
    return

}

type PMXDisplayData struct {
    Name string
    NameEnglish string
    SpecialFlag int
    FrameCount int
    DisplayFrames []*PMXDisplayFrame
}

func NewPMXDisplayData(f *FileReader, sizes[]int, decoder *encoding.Decoder)( dd *PMXDisplayData, err error) {
    dd = new(PMXDisplayData)
    var count int

    if count, err = f.GetInt32Little(); err != nil{
        return
    }
    if dd.Name, err = f.GetString(count, decoder); err != nil{
        return
    }
    if count, err = f.GetInt32Little(); err != nil{
        return
    }
    if dd.NameEnglish, err = f.GetString(count, decoder); err != nil{
        return
    }
    if dd.SpecialFlag, err = f.GetInt8Little(); err != nil{
        return
    }
    if dd.FrameCount, err = f.GetInt32Little(); err != nil{
        return
    }
    count = dd.FrameCount
    dd.DisplayFrames = make([]*PMXDisplayFrame, count)
    for i :=0; i< count;i++{
        if dd.DisplayFrames[i], err = NewPMXDisplayFrame(f, sizes); err != nil{
            return
        }
    }
    return
}


type PMXBone struct {

    Name string
    NameEnglish string

    Pos []float32
    ParentBoneIndex int
    Layer int
    Flags []byte
    TailPosition []float32//If indexed tail position flag is set then this is a bone index
    //or
    BoneIndex int
    //Inherit Bone  Used if either of the inherit flags are set. See Inherit Bone
    ParentIndex int
    ParentInfluence float32
    //Fixed axis Used if fixed axis flag is set. See Bone Fixed Axis
    AxisDirection []float32

    //Local co-ordinate Used if local co-ordinate flag is set. See Bone Local Co-ordinate
    XVector []float32
    ZVector []float32

    //External parent Used if external parent deform flag is set. See Bone External Parent
    ExternalParentIndex int

    //IK Used if IK flag is set. See Bone IK
    TargetIndex int
    LoopCount int
    LimitRadian float32
    LinkCount int
    Pmxiks []*PMXIK
}

func NewPMXBone(f *FileReader, sizes[]int, decoder *encoding.Decoder)(b *PMXBone, err error){
    boneIndexSize := sizes[5]
    b = new(PMXBone)
    var count int

    if count, err = f.GetInt32Little(); err != nil{
        return
    }
    if b.Name, err = f.GetString(count, decoder); err != nil{
        return
    }
    if count, err = f.GetInt32Little(); err != nil{
        return
    }
    if b.NameEnglish, err = f.GetString(count, decoder); err != nil{
        return
    }
    count = 3
    b.Pos = make([]float32, count)
    for i :=0; i< count;i++{
        if b.Pos[i], err = f.GetFloatLittle(); err != nil{
            return
        }
    }

    if b.ParentBoneIndex, err = f.GetIntLittle(boneIndexSize); err != nil{
        return
    }
    if b.Layer, err = f.GetInt32Little(); err != nil{
        return
    }
    if b.Flags, err = f.GetBytes(2); err != nil{
        return
    }
    //Indexed tail position
    if isBitSet(b.Flags[0], 0){
        if b.BoneIndex, err = f.GetIntLittle(boneIndexSize); err != nil{
            return
        }
    }else{
        count = 3
        b.TailPosition = make([]float32, count)
        for i :=0; i< count;i++{
            if b.TailPosition[i], err = f.GetFloatLittle(); err != nil{
                return
            }
        }
    }

    //Inherit
    if isBitSet(b.Flags[1], 0) || isBitSet(b.Flags[1], 1){
        if b.ParentIndex, err = f.GetIntLittle(boneIndexSize); err != nil{
            return
        }
        if b.ParentInfluence, err = f.GetFloatLittle(); err != nil{
            return
        }
    }

    //Fixed axis
    if isBitSet(b.Flags[1], 2){

        count = 3
        b.AxisDirection = make([]float32, count)
        for i :=0; i< count;i++{
            if b.AxisDirection[i], err = f.GetFloatLittle(); err != nil{
                return
            }
        }
    }
    //Local co-ordinate
    if isBitSet(b.Flags[1], 3){

        count = 3
        b.XVector = make([]float32, count)
        for i :=0; i< count;i++{
            if b.XVector[i], err = f.GetFloatLittle(); err != nil{
                return
            }
        }

        count = 3
        b.ZVector = make([]float32, count)
        for i :=0; i< count;i++{
            if b.ZVector[i], err = f.GetFloatLittle(); err != nil{
                return
            }
        }
    }
    //External parent deform
    if isBitSet(b.Flags[1], 5){
        if b.ExternalParentIndex, err = f.GetIntLittle(boneIndexSize); err != nil{
            return
        }
    }
    //IK
    if isBitSet(b.Flags[0], 5){
        if b.TargetIndex, err = f.GetIntLittle(boneIndexSize); err != nil{
            return
        }
        if b.LoopCount, err = f.GetInt32Little(); err != nil{
            return
        }
        if b.LimitRadian, err = f.GetFloatLittle(); err != nil{
            return
        }
        if b.LinkCount, err = f.GetInt32Little(); err != nil{
            return
        }
        count = b.LinkCount
        b.Pmxiks = make([]*PMXIK, count)
        for i :=0; i< count;i++{
            if b.Pmxiks[i], err = NewPMXIK(f, sizes); err != nil{
                return
            }
        }

    }
    return
}

type PMXJoint struct {
    Name string
    NameEnglish string
    Type int
    RigidBodyIndex1 int
    RigidBodyIndex2 int
    Position []float32
    Rotation []float32
    PositionMin []float32
    PositionMax []float32
    RotationMin []float32
    RotationMax []float32
    PositionSpring []float32
    RotationSpring []float32
}

func NewPMXJoint(f *FileReader, sizes[]int, decoder *encoding.Decoder) (j*PMXJoint, err error){
    rigidbodyIndexSize := sizes[7]
    j = new(PMXJoint)
    var count int

    if count, err = f.GetInt32Little(); err != nil{
        return
    }
    if j.Name, err = f.GetString(count, decoder); err != nil{
        return
    }
    if count, err = f.GetInt32Little(); err != nil{
        return
    }
    if j.NameEnglish, err = f.GetString(count, decoder); err != nil{
        return
    }
    if j.Type, err = f.GetInt8Little(); err != nil{
        return
    }
    if j.RigidBodyIndex1, err = f.GetIntLittle(rigidbodyIndexSize); err != nil{
        return
    }
    if j.RigidBodyIndex2, err = f.GetIntLittle(rigidbodyIndexSize); err != nil{
        return
    }
    count = 3
    j.Position = make([]float32, count)
    for i :=0; i< count;i++{
        if j.Position[i], err = f.GetFloatLittle(); err != nil{
            return
        }
    }
    count = 3
    j.Rotation = make([]float32, count)
    for i :=0; i< count;i++{
        if j.Rotation[i], err = f.GetFloatLittle(); err != nil{
            return
        }
    }
    count = 3
    j.PositionMin = make([]float32, count)
    for i :=0; i< count;i++{
        if j.PositionMin[i], err = f.GetFloatLittle(); err != nil{
            return
        }
    }
    count = 3
    j.PositionMax = make([]float32, count)
    for i :=0; i< count;i++{
        if j.PositionMax[i], err = f.GetFloatLittle(); err != nil{
            return
        }
    }
    count = 3
    j.RotationMin = make([]float32, count)
    for i :=0; i< count;i++{
        if j.RotationMin[i], err = f.GetFloatLittle(); err != nil{
            return
        }
    }
    count = 3
    j.RotationMax = make([]float32, count)
    for i :=0; i< count;i++{
        if j.RotationMax[i], err = f.GetFloatLittle(); err != nil{
            return
        }
    }
    count = 3
    j.PositionSpring = make([]float32, count)
    for i :=0; i< count;i++{
        if j.PositionSpring[i], err = f.GetFloatLittle(); err != nil{
            return
        }
    }
    count = 3
    j.RotationSpring = make([]float32, count)
    for i :=0; i< count;i++{
        if j.RotationSpring[i], err = f.GetFloatLittle(); err != nil{
            return
        }
    }
    return

}

const (

    MATERIAL_TEXTURE_REFERENCE int = iota
    MATERIAL_INTERNAL_REFERENCE

)
const(
    MATERIAL_MODE_DISABLE int = iota
    MATERIAL_MODE_SPH
    MATERIAL_MODE_SPA
)
type PMXMaterial struct {
    Name string
    NameEnglish string
    DiffuseColor []float32
    Alpha float32
    SpecularColor []float32
    SpecularStrength float32
    AmbientColor []float32
    DrawingFlag int
    EdgeColor []float32
    EdgeScale float32
    TextureIndex int
    EnvironmentIndex int
    EnvironmentBlendMode int
    ToonReference int//	0 = Texture reference, 1 = internal reference
    ToonValue int// toonReference 是1的时候，使用byte. 是0的时候，使用index;
    MetaData string
    //影响的vertex 数量
    SurfaceCount int
}

func NewPMXMaterial(f *FileReader, sizes[]int , decoder *encoding.Decoder)(m *PMXMaterial,err error)  {
    textureIndexSize := sizes[3]
    m = new(PMXMaterial)
    var count int

    if count, err = f.GetInt32Little(); err != nil{
        return
    }
    if m.Name, err = f.GetString(count, decoder); err != nil{
        return
    }
    if count, err = f.GetInt32Little(); err != nil{
        return
    }
    if m.NameEnglish, err = f.GetString(count, decoder); err != nil{
        return
    }
    count = 3
    m.DiffuseColor = make([]float32, count)
    for i :=0; i< count;i++{
        if m.DiffuseColor[i], err = f.GetFloatLittle(); err != nil{
            return
        }
    }
    if m.Alpha, err = f.GetFloatLittle(); err != nil{
        return
    }

    count = 3
    m.SpecularColor = make([]float32, count)
    for i :=0; i< count;i++{
        if m.SpecularColor[i], err = f.GetFloatLittle(); err != nil{
            return
        }
    }

    if m.SpecularStrength, err = f.GetFloatLittle(); err != nil{
        return
    }
    count = 3
    m.AmbientColor = make([]float32, count)
    for i :=0; i< count;i++{
        if m.AmbientColor[i], err = f.GetFloatLittle(); err != nil{
            return
        }
    }
    if m.DrawingFlag, err = f.GetInt8Little(); err != nil{
        return
    }
    count = 4
    m.EdgeColor = make([]float32, count)
    for i :=0; i< count;i++{
        if m.EdgeColor[i], err = f.GetFloatLittle(); err != nil{
            return
        }
    }
    if m.EdgeScale, err = f.GetFloatLittle(); err != nil{
        return
    }
    if m.TextureIndex, err = f.GetIntLittle(textureIndexSize); err != nil{
        return
    }
    if m.EnvironmentIndex, err = f.GetIntLittle(textureIndexSize); err != nil{
        return
    }
    if m.EnvironmentBlendMode, err = f.GetInt8Little(); err != nil{
        return
    }
    if m.ToonReference, err = f.GetInt8Little(); err != nil{
        return
    }
    if m.ToonReference == MATERIAL_INTERNAL_REFERENCE{
        if m.ToonValue, err = f.GetInt8Little(); err != nil{
            return
        }
    }else{
        if m.ToonValue, err = f.GetIntLittle(textureIndexSize); err != nil{
            return
        }
    }

    if count, err = f.GetInt32Little(); err != nil{
        return
    }
    if m.MetaData, err = f.GetString(count, decoder); err != nil{
        return
    }

    if m.SurfaceCount, err = f.GetInt32Little(); err != nil{
        return
    }
    return
}

const(
    MORPH_GROUP int = iota
    MORPH_VERTEX
    MORPH_BONE
    MORPH_UV
    MORPH_UV1
    MORPH_UV2
    MORPH_UV3
    MORPH_UV4
    MORPH_MATERIAL
    MORPH_FLIP
    MORPH_IMPULSE
)

type PMXMorphGroup struct {
    GroupMorphIndex int
    GroupInfluence float32
}

type PMXMorphVertex struct {
    VertexIndex int
    VertexTranslation []float32
}

type PMXMorphBone struct {
    BoneIndex int
    BoneTranslation []float32
    BoneRotation []float32
}

type PMXMorphUV struct {
    UVVertexIndex int
    UVFloats []float32
}

type PMXMorphMaterial struct {
    MaterialIndex int
    MaterialUnknow int
    MaterialDiffuse []float32
    MaterialSpecular []float32
    MaterialSpecularity float32
    MaterialAmbient []float32
    MaterialEdgeColor []float32
    MaterialEdgeSize float32
    MaterialTextureTint []float32
    MaterialEnviromentTint []float32
    MaterialToonTint []float32
}

type PMXMorphFlip struct {

    FlipMorphIndex int
    FlipInfluence float32
}
type PMXMorphImpluse struct {
    ImpluseRigidBodyIndex int
    ImpluseFlag int
    ImpluseMovementSpeed []float32
    ImpluseRotationTorque []float32
}

type PMXMorph struct {
    Name string
    NameEnglish string
    PanelType int
    MorphType int
    OffsetSize int

    MorphGroups []*PMXMorphGroup
    MorphVertexes []*PMXMorphVertex
    MorphBones []*PMXMorphBone
    MorphUVs []*PMXMorphUV
    MorphMaterials []*PMXMorphMaterial
    MorphFlips []*PMXMorphFlip
    MorphImpluses []*PMXMorphImpluse

}

func NewPMXMorph(f *FileReader, sizes[]int, decoder *encoding.Decoder)(m *PMXMorph, err error)  {
    vertexIndexSize := sizes[2];
    matrialIndexSize := sizes[4];
    boneIndexSize := sizes[5];
    morphIndexSize := sizes[6];
    rigidbodyIndexSize := sizes[7];
    m = new(PMXMorph)
    var count int

    if count, err = f.GetInt32Little(); err != nil{
        return
    }
    if m.Name, err = f.GetString(count, decoder); err != nil{
        return
    }
    if count, err = f.GetInt32Little(); err != nil{
        return
    }
    if m.NameEnglish, err = f.GetString(count, decoder); err != nil{
        return
    }
    if m.PanelType, err = f.GetInt8Little(); err != nil{
        return
    }
    if m.MorphType, err = f.GetInt8Little(); err != nil{
        return
    }
    if count, err = f.GetInt32Little(); err != nil{
        return
    }
    var size int
    switch m.MorphType {
    case MORPH_GROUP:
        m.MorphGroups = make([]*PMXMorphGroup, count)
        for j:=0; j< count;j++ {
            m.MorphGroups[j] = new(PMXMorphGroup)
            if m.MorphGroups[j].GroupMorphIndex, err = f.GetIntLittle(morphIndexSize); err != nil{
                return
            }
            if m.MorphGroups[j].GroupInfluence, err = f.GetFloatLittle(); err != nil{
                return
            }
        }
    case MORPH_VERTEX:
        m.MorphVertexes = make([]*PMXMorphVertex, count)
        for j:=0; j< count;j++ {
            m.MorphVertexes[j] = new(PMXMorphVertex)
            if m.MorphVertexes[j].VertexIndex, err = f.GetIntLittle(vertexIndexSize); err != nil{
                return
            }
            size = 3
            m.MorphVertexes[j].VertexTranslation = make([]float32, size)
            for i :=0; i< size;i++{
                if m.MorphVertexes[j].VertexTranslation[i], err = f.GetFloatLittle(); err != nil{
                    return
                }
            }

        }
    case MORPH_BONE:
        m.MorphBones = make([]*PMXMorphBone, count)
        for j:=0; j< count;j++ {
            m.MorphBones[j] = new(PMXMorphBone)
            if m.MorphBones[j].BoneIndex, err = f.GetIntLittle(boneIndexSize); err != nil{
                return
            }
            size = 3
            m.MorphBones[j].BoneTranslation = make([]float32, size)
            for i :=0; i< size;i++{
                if m.MorphBones[j].BoneTranslation[i], err = f.GetFloatLittle(); err != nil{
                    return
                }
            }
            size = 4
            m.MorphBones[j].BoneRotation = make([]float32, size)
            for i :=0; i< size;i++{
                if m.MorphBones[j].BoneRotation[i], err = f.GetFloatLittle(); err != nil{
                    return
                }
            }


        }

    case MORPH_UV, MORPH_UV1, MORPH_UV2, MORPH_UV3, MORPH_UV4:
        m.MorphUVs = make([]*PMXMorphUV, count)
        for j:=0; j< count;j++ {
            m.MorphUVs[j] = new(PMXMorphUV)
            if m.MorphUVs[j].UVVertexIndex, err = f.GetIntLittle(vertexIndexSize); err != nil{
                return
            }
            size = 4
            m.MorphUVs[j].UVFloats = make([]float32, size)
            for i :=0; i< size;i++{
                if m.MorphUVs[j].UVFloats[i], err = f.GetFloatLittle(); err != nil{
                    return
                }
            }

        }
    case MORPH_MATERIAL:
        m.MorphMaterials = make([]*PMXMorphMaterial, count)
        for j:=0; j< count;j++ {
            m.MorphMaterials[j] = new(PMXMorphMaterial)
            if m.MorphMaterials[j].MaterialIndex, err = f.GetIntLittle(matrialIndexSize); err != nil{
                return
            }
            if m.MorphMaterials[j].MaterialUnknow, err = f.GetInt8Little(); err != nil{
                return
            }
            size = 4
            m.MorphMaterials[j].MaterialDiffuse = make([]float32, size)
            for i :=0; i< size;i++{
                if m.MorphMaterials[j].MaterialDiffuse[i], err = f.GetFloatLittle(); err != nil{
                    return
                }
            }
            size = 3
            m.MorphMaterials[j].MaterialSpecular = make([]float32, size)
            for i :=0; i< size;i++{
                if m.MorphMaterials[j].MaterialSpecular[i], err = f.GetFloatLittle(); err != nil{
                    return
                }
            }
            if m.MorphMaterials[j].MaterialSpecularity, err = f.GetFloatLittle(); err != nil{
                return
            }
            size = 3
            m.MorphMaterials[j].MaterialAmbient = make([]float32, size)
            for i :=0; i< size;i++{
                if m.MorphMaterials[j].MaterialAmbient[i], err = f.GetFloatLittle(); err != nil{
                    return
                }
            }
            size = 4
            m.MorphMaterials[j].MaterialEdgeColor = make([]float32, size)
            for i :=0; i< size;i++{
                if m.MorphMaterials[j].MaterialEdgeColor[i], err = f.GetFloatLittle(); err != nil{
                    return
                }
            }

            if m.MorphMaterials[j].MaterialEdgeSize, err = f.GetFloatLittle(); err != nil{
                return
            }

            size = 4
            m.MorphMaterials[j].MaterialTextureTint = make([]float32, size)
            for i :=0; i< size;i++{
                if m.MorphMaterials[j].MaterialTextureTint[i], err = f.GetFloatLittle(); err != nil{
                    return
                }
            }
            size = 4
            m.MorphMaterials[j].MaterialEnviromentTint = make([]float32, size)
            for i :=0; i< size;i++{
                if m.MorphMaterials[j].MaterialEnviromentTint[i], err = f.GetFloatLittle(); err != nil{
                    return
                }
            }
            size = 4
            m.MorphMaterials[j].MaterialToonTint = make([]float32, size)
            for i :=0; i< size;i++{
                if m.MorphMaterials[j].MaterialToonTint[i], err = f.GetFloatLittle(); err != nil{
                    return
                }
            }


        }
    case MORPH_FLIP:
        m.MorphFlips = make([]*PMXMorphFlip, count)
        for j:=0; j< count;j++ {
            m.MorphFlips[j] = new(PMXMorphFlip)
            if m.MorphFlips[j].FlipMorphIndex, err = f.GetIntLittle(morphIndexSize); err != nil{
                return
            }
            if m.MorphFlips[j].FlipInfluence, err = f.GetFloatLittle(); err != nil{
                return
            }

        }
    case MORPH_IMPULSE:
        m.MorphImpluses = make([]*PMXMorphImpluse, count)
        for j:=0; j< count;j++ {
            m.MorphImpluses[j] = new(PMXMorphImpluse)
            if m.MorphImpluses[j].ImpluseRigidBodyIndex, err = f.GetIntLittle(rigidbodyIndexSize); err != nil{
                return
            }
            if m.MorphImpluses[j].ImpluseFlag, err = f.GetInt8Little(); err != nil{
                return
            }
            size = 3
            m.MorphImpluses[j].ImpluseMovementSpeed = make([]float32, size)
            for i :=0; i< size;i++{
                if m.MorphImpluses[j].ImpluseMovementSpeed[i], err = f.GetFloatLittle(); err != nil{
                    return
                }
            }
            size = 3
            m.MorphImpluses[j].ImpluseRotationTorque = make([]float32, size)
            for i :=0; i< size;i++{
                if m.MorphImpluses[j].ImpluseRotationTorque[i], err = f.GetFloatLittle(); err != nil{
                    return
                }
            }
        }

    }
    return

}

type PMXRigidBody struct {
    Name string
    NameEnglish string
    RelatedBoneIndex int
    GroupId int
    NonCollisionGroup int
    Shape int
    ShapeSize []float32
    ShapePosition []float32
    ShapeRotation []float32
    Mass float32
    MoveAttenuation float32
    RotationDamping float32
    Repulsion float32
    FrictionForce float32
    PhysicsMode int
}

func NewPMXRigidBody(f *FileReader, sizes[]int, decoder *encoding.Decoder)(rb *PMXRigidBody, err error)    {
    boneIndexSize := sizes[5]
    rb = new(PMXRigidBody)
    var count int

    if count, err = f.GetInt32Little(); err != nil{
        return
    }
    if rb.Name, err = f.GetString(count, decoder); err != nil{
        return
    }
    if count, err = f.GetInt32Little(); err != nil{
        return
    }
    if rb.NameEnglish, err = f.GetString(count, decoder); err != nil{
        return
    }
    if rb.RelatedBoneIndex, err = f.GetIntLittle(boneIndexSize); err != nil{
        return
    }
    if rb.GroupId, err = f.GetInt8Little(); err != nil{
        return
    }
    if rb.NonCollisionGroup, err = f.GetInt16Little(); err != nil{
        return
    }
    if rb.Shape, err = f.GetInt8Little(); err != nil{
        return
    }

    count = 3
    rb.ShapeSize = make([]float32, count)
    for i :=0; i< count;i++{
        if rb.ShapeSize[i], err = f.GetFloatLittle(); err != nil{
            return
        }
    }
    count = 3
    rb.ShapePosition = make([]float32, count)
    for i :=0; i< count;i++{
        if rb.ShapePosition[i], err = f.GetFloatLittle(); err != nil{
            return
        }
    }
    count = 3
    rb.ShapeRotation = make([]float32, count)
    for i :=0; i< count;i++{
        if rb.ShapeRotation[i], err = f.GetFloatLittle(); err != nil{
            return
        }
    }
    if rb.Mass, err = f.GetFloatLittle(); err != nil{
        return
    }
    if rb.MoveAttenuation, err = f.GetFloatLittle(); err != nil{
        return
    }
    if rb.RotationDamping, err = f.GetFloatLittle(); err != nil{
        return
    }
    if rb.Repulsion, err = f.GetFloatLittle(); err != nil{
        return
    }
    if rb.FrictionForce, err = f.GetFloatLittle(); err != nil{
        return
    }
    if rb.PhysicsMode, err = f.GetInt8Little(); err != nil{
        return
    }
    return
}

type PMXAnchorRigidBody struct {
    RigidBodyIndex int
    VertexIndex int
    NearMode int
}

type PMXSoftBody struct {
    Name string
    NameEnglish string
    ShapeType int
    MaterialIndex int
    Group int
    NoCollisionMask int
    Flags int
    B_LinkCreateDistance int
    ClustersCount int
    TotalMass float32
    CollisionMargin float32
    AerodynamicsModel int
    ConfigVCF float32
    ConfigDP float32
    ConfigDG float32
    ConfigLF float32
    ConfigPR float32
    ConfigVC float32
    ConfigDF float32
    ConfigMT float32
    ConfigCHR float32
    ConfigKHR float32
    ConfigSHR float32
    ConfigAHR float32
    ClusterSRHR_CL float32
    ClusterSKHR_CL float32
    ClusterSSHR_CL float32
    ClusterSR_SPLT_CL float32
    ClusterSK_SPLT_CL float32
    ClusterSS_SPLT_CL float32
    InterationV_IT int
    InterationP_IT int
    InterationD_IT int
    InterationC_IT int
    MaterialLST int
    MaterialAST int
    MaterialVST int
    AnchorRigidBodyCount int
    AnchorRigidBodies []*PMXAnchorRigidBody
    VertexPinCount int
    VertexPinIndexs []int
}

func NewPMXSoftBody(f *FileReader, sizes[]int, decoder *encoding.Decoder)(sb *PMXSoftBody, err error)    {
    vertexIndexSize := sizes[2];
    materialIndexSize := sizes[4];
    rigidbodyIndexSize := sizes[7];
    sb = new(PMXSoftBody)
    var count int

    if count, err = f.GetInt32Little(); err != nil{
        return
    }
    if sb.Name, err = f.GetString(count, decoder); err != nil{
        return
    }
    if count, err = f.GetInt32Little(); err != nil{
        return
    }
    if sb.NameEnglish, err = f.GetString(count, decoder); err != nil{
        return
    }
    if sb.ShapeType, err = f.GetInt8Little(); err != nil{
        return
    }
    if sb.MaterialIndex, err = f.GetIntLittle(materialIndexSize); err != nil{
        return
    }
    if sb.Group, err = f.GetInt8Little(); err != nil{
        return
    }
    if sb.NoCollisionMask, err = f.GetInt16Little(); err != nil{
        return
    }
    if sb.Flags, err = f.GetInt8Little(); err != nil{
        return
    }
    if sb.B_LinkCreateDistance, err = f.GetInt32Little(); err != nil{
        return
    }

    if sb.ClustersCount, err = f.GetInt32Little(); err != nil{
        return
    }
    if sb.TotalMass, err = f.GetFloatLittle(); err != nil{
        return
    }
    if sb.CollisionMargin, err = f.GetFloatLittle(); err != nil{
        return
    }
    if sb.AerodynamicsModel, err = f.GetInt32Little(); err != nil{
        return
    }

    if sb.ConfigVCF, err = f.GetFloatLittle(); err != nil{
        return
    }
    if sb.ConfigDP, err = f.GetFloatLittle(); err != nil{
        return
    }
    if sb.ConfigDG, err = f.GetFloatLittle(); err != nil{
        return
    }
    if sb.ConfigLF, err = f.GetFloatLittle(); err != nil{
        return
    }
    if sb.ConfigPR, err = f.GetFloatLittle(); err != nil{
        return
    }
    if sb.ConfigVC, err = f.GetFloatLittle(); err != nil{
        return
    }

    if sb.ConfigDF, err = f.GetFloatLittle(); err != nil{
        return
    }
    if sb.ConfigMT, err = f.GetFloatLittle(); err != nil{
        return
    }
    if sb.ConfigCHR, err = f.GetFloatLittle(); err != nil{
        return
    }
    if sb.ConfigKHR, err = f.GetFloatLittle(); err != nil{
        return
    }
    if sb.ConfigSHR, err = f.GetFloatLittle(); err != nil{
        return
    }
    if sb.ConfigAHR, err = f.GetFloatLittle(); err != nil{
        return
    }

    if sb.ClusterSRHR_CL, err = f.GetFloatLittle(); err != nil{
        return
    }
    if sb.ClusterSKHR_CL, err = f.GetFloatLittle(); err != nil{
        return
    }
    if sb.ClusterSSHR_CL, err = f.GetFloatLittle(); err != nil{
        return
    }
    if sb.ClusterSR_SPLT_CL, err = f.GetFloatLittle(); err != nil{
        return
    }
    if sb.ClusterSK_SPLT_CL, err = f.GetFloatLittle(); err != nil{
        return
    }
    if sb.ClusterSS_SPLT_CL, err = f.GetFloatLittle(); err != nil{
        return
    }
    if sb.InterationV_IT, err = f.GetInt32Little(); err != nil{
        return
    }

    if sb.InterationP_IT, err = f.GetInt32Little(); err != nil{
        return
    }

    if sb.InterationD_IT, err = f.GetInt32Little(); err != nil{
        return
    }

    if sb.InterationC_IT, err = f.GetInt32Little(); err != nil{
        return
    }

    if sb.MaterialLST, err = f.GetInt32Little(); err != nil{
        return
    }

    if sb.MaterialAST, err = f.GetInt32Little(); err != nil{
        return
    }

    if sb.MaterialVST, err = f.GetInt32Little(); err != nil{
        return
    }
    if sb.AnchorRigidBodyCount, err = f.GetInt32Little(); err != nil{
        return
    }
    count = sb.AnchorRigidBodyCount
    sb.AnchorRigidBodies = make([]*PMXAnchorRigidBody, count)
    for j:=0; j< count;j++ {
        sb.AnchorRigidBodies[j] = new(PMXAnchorRigidBody)
        if sb.AnchorRigidBodies[j].RigidBodyIndex, err = f.GetIntLittle(rigidbodyIndexSize); err != nil{
            return
        }
        if sb.AnchorRigidBodies[j].VertexIndex, err = f.GetIntLittle(vertexIndexSize); err != nil{
            return
        }
        if sb.AnchorRigidBodies[j].NearMode, err = f.GetInt8Little(); err != nil{
            return
        }

    }

    if sb.VertexPinCount, err = f.GetInt32Little(); err != nil{
        return
    }
    count = sb.VertexPinCount
    sb.VertexPinIndexs = make([]int, count)
    for j:=0; j< count;j++ {
        if sb.VertexPinIndexs[j], err = f.GetIntLittle(vertexIndexSize); err != nil{
            return
        }
    }
    return
}

type PMX struct {
    Directory string
    FileName string
    fr * FileReader

    Name string
    EnglishName string
    Comment string
    EnglishComment string
    decoder *encoding.Decoder
    globalSizes []int
    Version float32

    Vertices []*PMXVertex

    Triangles []int

    Materials []*PMXMaterial

    Textures []string

    Bones [] *PMXBone

    Morphs []*PMXMorph

    DisplayDatas []*PMXDisplayData

    RigidBodies []*PMXRigidBody

    Joints []*PMXJoint

    SoftBodies []*PMXSoftBody

}

func (p *PMX)Load(directory, fileName string)(err error)  {
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
    log.Println("header ok;version ", p.Version)
    if err = p.parseVertices(); err != nil{
        return
    }
    log.Println("vertices ok")
    if err = p.parseTriangles(); err != nil{
        return
    }
    log.Println("triangles ok")
    if err = p.parseTextures(); err != nil{
        return
    }
    log.Println("textures ok")
    if err = p.parseMaterials(); err != nil{
        return
    }
    log.Println("materials ok")
    if err = p.parseBones(); err != nil {
        return
    }
    log.Println("bones ok")
    if err = p.parseMorphs(); err != nil{
        return
    }
    log.Println("morphs ok")
    if err = p.parseDisplayDatas(); err != nil{
        return
    }
    log.Println("display datas ok")
    if err = p.parseRigidBodys(); err != nil{
        return
    }
    log.Println("rigid body ok")
    if err = p.parseJoints(); err != nil{
        return
    }
    log.Println("joints ok")
    if p.Version == 2.1{
        if err = p.parseSoftBodys(); err != nil{
            return
        }
        log.Println("softbodys ok")
    }

    return
}
func (p *PMX)checkHeader()(err error)  {
    const PMXHEADER string  = "PMX "
    var header string
    count := 4
    if header, err = p.fr.GetString(count, unicode.UTF8.NewDecoder());err != nil{
        return
    }
    if header != PMXHEADER{
        err = errors.New("PMX Header check fail")
        return
    }
    if p.Version, err = p.fr.GetFloatLittle();err !=nil{
        return
    }
    if p.Version != 2.0 && p.Version != 2.1{
        err = errors.New("Version is  not 2.0 or 2.1")
        return
    }

    if count, err = p.fr.GetInt8Little(); err != nil{
        return
    }
    p.globalSizes = make([]int, count)
    for i:= 0;i< count;i++{
        if p.globalSizes[i], err = p.fr.GetInt8Little(); err != nil{
            return
        }
    }
    if p.globalSizes[0] == 0{
        p.decoder = unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
    }else{
        p.decoder = unicode.UTF8.NewDecoder()
    }
    if count, err = p.fr.GetInt32Little(); err != nil{
        return
    }
    if p.Name, err = p.fr.GetString(count, p.decoder); err != nil{
        return
    }
    if count, err = p.fr.GetInt32Little(); err != nil{
        return
    }
    if p.EnglishName, err = p.fr.GetString(count, p.decoder); err != nil{
        return
    }
    if count, err = p.fr.GetInt32Little(); err != nil{
        return
    }
    if p.Comment, err = p.fr.GetString(count, p.decoder); err != nil{
        return
    }
    if count, err = p.fr.GetInt32Little(); err != nil{
        return
    }
    if p.EnglishComment, err = p.fr.GetString(count, p.decoder); err != nil{
        return
    }

    return
}
func (p *PMX)parseVertices()(err error)  {
    var count int
    if count, err = p.fr.GetInt32Little(); err != nil{
        return
    }
    p.Vertices = make([]*PMXVertex, count)
    for i := 0;i < count; i++ {
        if p.Vertices[i], err = NewPMXVertex(p.fr, p.globalSizes, p.decoder); err != nil{
            return
        }
    }
    return
}
func (p *PMX)parseTriangles()(err error)  {
    vertexIndexSize := p.globalSizes[2]
    var count int
    if count, err = p.fr.GetInt32Little(); err != nil{
        return
    }
    p.Triangles = make([]int, count)
    for i := 0;i < count; i++ {
        if p.Triangles[i], err = p.fr.GetIntLittle(vertexIndexSize); err != nil{
            return
        }
    }
    return

}
func (p *PMX)parseTextures()(err error)  {
    var count int
    var size int
    if count, err = p.fr.GetInt32Little(); err != nil{
        return
    }
    p.Textures = make([]string, count)
    for i := 0;i < count; i++ {
        if size, err = p.fr.GetInt32Little(); err != nil{
            return
        }
        if p.Textures[i], err = p.fr.GetString(size, p.decoder); err != nil{
            return
        }
    }
    return

}
func (p *PMX)parseMaterials()(err error)  {
    var count int
    if count, err = p.fr.GetInt32Little(); err != nil{
        return
    }
    p.Materials = make([]*PMXMaterial, count)
    for i := 0;i < count; i++ {
        if p.Materials[i], err = NewPMXMaterial(p.fr, p.globalSizes, p.decoder); err != nil{
            return
        }
    }
    return
}
func (p *PMX)parseBones()(err error)  {
    var count int
    if count, err = p.fr.GetInt32Little(); err != nil{
        return
    }
    p.Bones = make([]*PMXBone, count)
    for i := 0;i < count; i++ {
        if p.Bones[i], err = NewPMXBone(p.fr, p.globalSizes, p.decoder); err != nil{
            return
        }
    }
    return
}
func (p *PMX)parseMorphs()(err error)  {
    var count int
    if count, err = p.fr.GetInt32Little(); err != nil{
        return
    }
    p.Morphs = make([]*PMXMorph, count)
    for i := 0;i < count; i++ {
        if p.Morphs[i], err = NewPMXMorph(p.fr, p.globalSizes, p.decoder); err != nil{
            return
        }
    }
    return

}
func (p *PMX)parseDisplayDatas()(err error)  {
    var count int
    if count, err = p.fr.GetInt32Little(); err != nil{
        return
    }
    p.DisplayDatas = make([]*PMXDisplayData, count)
    for i := 0;i < count; i++ {
        if p.DisplayDatas[i], err = NewPMXDisplayData(p.fr, p.globalSizes, p.decoder); err != nil{
            return
        }
    }
    return

}

func (p *PMX)parseRigidBodys()(err error)  {

    var count int
    if count, err = p.fr.GetInt32Little(); err != nil{
        return
    }
    p.RigidBodies = make([]*PMXRigidBody, count)
    for i := 0;i < count; i++ {
        if p.RigidBodies[i], err = NewPMXRigidBody(p.fr, p.globalSizes, p.decoder); err != nil{
            return
        }
    }
    return
}

func (p *PMX)parseJoints()(err error)  {
    var count int
    if count, err = p.fr.GetInt32Little(); err != nil{
        return
    }
    p.Joints = make([]*PMXJoint, count)
    for i := 0;i < count; i++ {
        if p.Joints[i], err = NewPMXJoint(p.fr, p.globalSizes, p.decoder); err != nil{
            return
        }
    }
    return

}
func (p *PMX)parseSoftBodys()(err error)  {

    var count int
    if count, err = p.fr.GetInt32Little(); err != nil{
        return
    }
    p.SoftBodies = make([]*PMXSoftBody, count)
    for i := 0;i < count; i++ {
        if p.SoftBodies[i], err = NewPMXSoftBody(p.fr, p.globalSizes, p.decoder); err != nil{
            return
        }
    }
    return
}