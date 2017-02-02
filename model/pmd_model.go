package model

import (
    "github.com/zdglf/mmd/util"
    "github.com/zdglf/mmd/gles2"
    "log"

    "strings"
    "regexp"
    "fmt"
    "path"
    "unsafe"
)

type GLBuf struct {
    size int
    buffer int
}

type AttrArrBuf struct {
    size int
    array []float32
    attribute string
}

type PMDModel struct {
    program int32
    pmd *util.PMD
    programMap map[string]int32

    cameraPosition []float32
    ignoreCameraMotion bool
    rotx int
    roty int
    distance float32
    center []float32
    fovy float32
    drawEdge bool
    edgeThickness float32
    edgeColor []float32
    lightDirection []float32
    lightDistance float32
    lightColor []float32
    drawSelfShadow bool
    drawAxes bool
    drawCenterPoint bool
    fps float32
    realFps float32
    playing bool
    frame int
    upPos []float32
    x int
    y int
    width int
    height int

    vbuffers map[string]GLBuf
    ibuffer int32

    textureManager *TextureManager

    
}

func (m *PMDModel)LoadFile(filePath string, fileName string) bool  {
    m.pmd = new(util.PMD)
    log.Println("start pmd load")
    if err := m.pmd.Load(filePath, fileName); err != nil{
        m.pmd = nil
        log.Println(err)
        return false
    }else {
        log.Println(m.pmd.Name)
        log.Println(m.pmd.Comment)
        return true;
    }
}

func (m *PMDModel)InitShader(vShader string, fShader string) bool {
    vertexShader := gles2.CreateShader(gles2.VERTEX_SHADER)
    gles2.ShaderSource(vertexShader,1, []string{vShader}, []int32{int32(len(vShader))})
    gles2.CompileShader(vertexShader)
    var compileStatus int32
    gles2.GetShaderiv(vertexShader, gles2.COMPILE_STATUS, &compileStatus)
    if compileStatus == 0{
        var maxLength int32 = 1024
        infoBytes := make([]byte, maxLength)
        var realLength int32
        gles2.GetShaderInfoLog(vertexShader, maxLength, &realLength, infoBytes)
        log.Println("vshader", string(infoBytes[:realLength]))
        gles2.DeleteShader(vertexShader)
        vertexShader = 0
    }
    if vertexShader == 0{
        return false
    }

    fragmentShader := gles2.CreateShader(gles2.FRAGMENT_SHADER)
    gles2.ShaderSource(fragmentShader, 1, []string{fShader}, []int32{int32(len(fShader))})
    gles2.CompileShader(fragmentShader)
    gles2.GetShaderiv(fragmentShader, gles2.COMPILE_STATUS, &compileStatus)
    if compileStatus == 0{
        var maxLength int32 = 1024
        infoBytes := make([]byte, maxLength)
        var realLength int32
        gles2.GetShaderInfoLog(fragmentShader, maxLength, &realLength, infoBytes)
        log.Println("fshader", string(infoBytes[:realLength]))
        gles2.DeleteShader(fragmentShader)
        fragmentShader = 0
    }
    if fragmentShader == 0{
        return false
    }
    m.program = gles2.CreateProgram()
    gles2.AttachShader(m.program, vertexShader)
    gles2.AttachShader(m.program, fragmentShader)
    gles2.LinkProgram(m.program)
    var linkStatus int32
    gles2.GetProgramiv(m.program, gles2.LINK_STATUS, &linkStatus)
    if linkStatus == 0{
        var maxLength int32 = 1024
        infoBytes := make([]byte, maxLength)
        var realLength int32
        gles2.GetProgramInfoLog(m.program, maxLength, &realLength, infoBytes)
        log.Println("link", string(infoBytes[:realLength]))
        m.program = 0
    }
    if m.program == 0{
        return false
    }
    gles2.UseProgram(m.program)

    attributes := ""
    uniforms := ""

    shaders := []string{vShader, fShader}
    for _,shader := range shaders{
        re := regexp.MustCompile("\\/\\*[\\s\\S]*?\\*\\/")
        tmp := re.ReplaceAllString(shader, "")
        re = regexp.MustCompile("\\/\\/[^\\n]*")
        tmp = re.ReplaceAllString(tmp, "")
        datas := strings.Split(tmp, ";")
        for _,d := range datas{
            re = regexp.MustCompile("(\\w+)(\\[\\d+\\])?\\s*$")
            t := re.FindString(d)
            if strings.Contains(d, "uniform"){
                if !strings.Contains(uniforms, t) {
                    uniforms = uniforms + ";" + t
                }

            }else if strings.Contains(d, "attribute"){
                if !strings.Contains(attributes, t) {
                    attributes = attributes + ";" + t
                }

            }

        }

    }
    as := strings.Split(attributes, ";")
    us := strings.Split(uniforms, ";")
    m.programMap = make(map[string]int32)
    for _, a := range as{
        m.programMap[a] = gles2.GetAttribLocation(m.program, a)
        gles2.EnableVertexAttribArray(m.programMap[a])
    }
    for _, u := range us{
        m.programMap[u] = gles2.GetUniformLocation(m.program, u)
    }
    log.Println(len(as), as)
    log.Println(len(us), us)

    return true
}
func (m *PMDModel)InitParam(x int, y int, width int, height int, toonDir string){
    m.cameraPosition = []float32{0.0, 0.0, -15.0}
    m.ignoreCameraMotion = false
    m.rotx = 0
    m.roty = 0
    m.distance = 15.0
    m.center =  []float32{0.0, 10.0, 0.0}
    m.fovy = 40
    m.drawEdge = false
    m.edgeThickness = 0.004
    m.edgeColor = []float32{0.0, 0.0, 0.0, 1.0}
    m.lightDirection = []float32{1.0, -1.0, -1.0}
    m.lightDistance = 100.
    m.lightColor = []float32{0.6, 0.6, 0.6, 1.0}
    m.drawSelfShadow = true
    m.drawAxes = true
    m.drawCenterPoint = false
    m.fps = 30.
    m.realFps = m.fps
    m.playing = false
    m.frame =-1
    m.upPos = []float32{0,1,0.0}
    m.x = x
    m.y = y
    m.width = width
    m.height = height
    m.initVertices()
    m.initIndices()
    m.initTextures(toonDir)
}

func (m *PMDModel)initVertices() {
    m.vbuffers = make(map[string]GLBuf)
    if m.pmd != nil{
        length := len(m.pmd.Vertices)
        weight := make([]float32, length)
        vectors1 := make([]float32, 3*length)
        vectors2 := make([]float32, 3*length)
        rotations1 := make([]float32, 4*length)
        rotations2 := make([]float32, 4*length)
        positions1 := make([]float32, 3*length)
        positions2 := make([]float32, 3*length)
        morphVec := make([]float32, 3*length)
        normals := make([]float32, 3*length)
        uvs := make([]float32, 2*length)
        edge := make([]float32, length)

        for i:= 0; i< length; i++ {
            vertex := m.pmd.Vertices[i]
            bone1 := m.pmd.Bones[vertex.BoneNum1]
            bone2 := m.pmd.Bones[vertex.BoneNum2]
            weight[i] = float32(vertex.BoneWeight)
            vectors1[3 * i] = vertex.X - bone1.HeadPos[0]
            vectors1[3 * i + 1] = vertex.Y - bone1.HeadPos[1]
            vectors1[3 * i + 2] = vertex.Z - bone1.HeadPos[2]
            vectors2[3 * i] = vertex.X - bone2.HeadPos[0]
            vectors2[3 * i + 1] = vertex.Y - bone2.HeadPos[1]
            vectors2[3 * i + 2] = vertex.Z - bone2.HeadPos[2]
            positions1[3 * i] = bone1.HeadPos[0]
            positions1[3 * i + 1] = bone1.HeadPos[1]
            positions1[3 * i + 2] = bone1.HeadPos[2]
            positions2[3 * i] = bone2.HeadPos[0]
            positions2[3 * i + 1] = bone2.HeadPos[1]
            positions2[3 * i + 2] = bone2.HeadPos[2]
            rotations1[4 * i + 3] = 1
            rotations2[4 * i + 3] = 1
            normals[3 * i] = vertex.NX
            normals[3 * i + 1] = vertex.NY
            normals[3 * i + 2] = vertex.NZ
            uvs[2 * i] = vertex.U
            uvs[2 * i + 1] = vertex.V
            edge[i] = 1. - float32(vertex.EdgeFlag)
        }
        tmpArr := make([]AttrArrBuf, 0)
        tmpArr = append(tmpArr, AttrArrBuf{1, weight, "aBoneWeight"})
        tmpArr = append(tmpArr, AttrArrBuf{3, vectors1, "aVectorFromBone1"})
        tmpArr = append(tmpArr, AttrArrBuf{3, vectors2, "aVectorFromBone2"})
        tmpArr = append(tmpArr, AttrArrBuf{4, rotations1, "aBone1Rotation"})
        tmpArr = append(tmpArr, AttrArrBuf{4, rotations2, "aBone2Rotation"})
        tmpArr = append(tmpArr, AttrArrBuf{3, positions1, "aBone1Position"})
        tmpArr = append(tmpArr, AttrArrBuf{3, positions2, "aBone2Position"})
        tmpArr = append(tmpArr, AttrArrBuf{3, morphVec, "aMultiPurposeVector"})
        tmpArr = append(tmpArr, AttrArrBuf{3, normals, "aVertexNormal"})
        tmpArr = append(tmpArr, AttrArrBuf{2, uvs, "aTextureCoord"})
        tmpArr = append(tmpArr, AttrArrBuf{1, edge, "aVertexEdge"})
        for _, tmp := range tmpArr{
            buffer := make([]int32, 1)
            gles2.GenBuffers(1, buffer)
            gles2.BindBuffer(gles2.ARRAY_BUFFER, buffer[0])
            gles2.BufferData(gles2.ARRAY_BUFFER, tmp.size*len(tmp.array), unsafe.Pointer(&tmp.array[0]), gles2.STATIC_DRAW)
            m.vbuffers[tmp.attribute] = GLBuf{tmp.size, int(buffer[0])}
        }
        gles2.BindBuffer(gles2.ARRAY_BUFFER, 0)

    }

}

func (m *PMDModel)initIndices() {
    indices := m.pmd.Triangles
    buffer := make([]int32, 1)
    gles2.GenBuffers(1, buffer)
    gles2.BindBuffer(gles2.ELEMENT_ARRAY_BUFFER, buffer[0])
    gles2.BufferData(gles2.ELEMENT_ARRAY_BUFFER, 4*len(indices), unsafe.Pointer(&indices[0]), gles2.STATIC_DRAW)
    m.ibuffer = buffer[0]

}

func (m *PMDModel)initTextures(toonDir string)(err error) {
    m.textureManager = NewTextureManager()
    toonFiles := []string{"toon00.bmp", "toon01.bmp", "toon02.bmp", "toon03.bmp", "toon04.bmp", "toon05.bmp",
        "toon06.bmp", "toon07.bmp", "toon08.bmp", "toon09.bmp", "toon10.bmp"}
    materials := m.pmd.Materials
    for _, material := range materials{
        if material.Textures == nil{
            material.Textures = make(map[string]int32)
        }
        toonIndex := material.ToonIndex
        fileName := fmt.Sprintf("toon%02d.bmp", toonIndex)
        if m.pmd.ToonFileNames!=nil||len(m.pmd.ToonFileNames)<=toonIndex||m.pmd.ToonFileNames[toonIndex]==""{
            fileName = path.Join(toonDir, fileName)
        }else{
            isInToonFiles := false
            fileName = m.pmd.ToonFileNames[toonIndex]
            for _,toonFile := range toonFiles{
                if toonFile == fileName{
                    isInToonFiles = true
                    break
                }
            }
            if isInToonFiles{
                fileName = path.Join(toonDir, fileName)
            }else{
                fileName = path.Join(m.pmd.Directory, fileName)
            }
        }
        if material.Textures["toon"], err = m.textureManager.Get("toon", fileName); err !=nil{
            return
        }
        if material.TextureFileName != ""{
            textureFiles := strings.Split(material.TextureFileName, "*")
            for _,textureFile := range textureFiles{
                log.Println(textureFile)
                size := len(textureFile)
                if(size < 4){
                    continue
                }
                endFix := strings.ToUpper(textureFile[size-4:])
                switch endFix {
                case ".SPH":
                    if material.Textures["sph"], err = m.textureManager.Get("sph", path.Join(m.pmd.Directory, textureFile)); err !=nil{
                        return
                    }
                case ".SPA":
                    if material.Textures["spa"], err = m.textureManager.Get("spa", path.Join(m.pmd.Directory, textureFile)); err !=nil{
                        return
                    }
                default:
                    if material.Textures["regular"], err = m.textureManager.Get("regular", path.Join(m.pmd.Directory, textureFile)); err !=nil{
                        return
                    }

                    
                }
            }
        }

    }
    return

}

func (m *PMDModel)LoadMotion(filePath string) bool{
        return true
}
func (m *PMDModel)getFrameCount() int{
        return 0;
}
func (m *PMDModel)InitFrame(index int){

}
func (m *PMDModel)Render(){
}
