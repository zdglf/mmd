package model

import (
    "github.com/zdglf/mmd/util"
    "github.com/zdglf/mmd/gles2"
    "github.com/go-gl/mathgl/mgl32"
    "log"

    "strings"
    "fmt"
    "path"
    "unsafe"
)

type GLBuf struct {
    size int
    buffer int32
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

    cameraPosition mgl32.Vec3
    ignoreCameraMotion bool
    rotx float32
    roty float32
    distance float32
    center mgl32.Vec3
    fovy float32
    drawEdge bool
    edgeThickness float32
    edgeColor []float32
    lightDirection mgl32.Vec3
    lightDistance float32
    lightColor []float32
    drawSelfShadow bool
    drawAxes bool
    drawCenterPoint bool
    fps float32
    realFps float32
    playing bool
    frame int
    upPos mgl32.Vec3
    x int32
    y int32
    width int32
    height int32

    vbuffers map[string]GLBuf
    ibuffer int32

    textureManager *TextureManager

    viewMatrix mgl32.Mat4
    pMatrix mgl32.Mat4
    modelMatrix mgl32.Mat4
    mvMatrix mgl32.Mat4
    nMatrix mgl32.Mat4

    
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
    m.programMap = make(map[string]int32)

    var count int32
    gles2.GetProgramiv(m.program, gles2.ACTIVE_ATTRIBUTES, &count)

    for i:= 0;i <int(count);i++{
        var bufSize int32 = 20
        var realSize int32
        var attr_type int32
        var nameSize int32
        buff := make([]byte, bufSize)
        gles2.GetActiveAttrib(m.program, int32(i), bufSize, &nameSize, &realSize, &attr_type, buff)
        name := string(buff[:nameSize])
        m.programMap[name] = int32(i)
        gles2.EnableVertexAttribArray(int32(i))
        log.Println(i,name)
    }
    gles2.GetProgramiv(m.program, gles2.ACTIVE_UNIFORMS, &count)
    for i:= 0;i <int(count);i++{
        var bufSize int32 = 20
        var realSize int32
        var attr_type int32
        var nameSize int32
        buff := make([]byte, bufSize)
        gles2.GetActiveUniform(m.program, int32(i), bufSize, &nameSize, &realSize, &attr_type, buff)
        name := string(buff[:nameSize])
        m.programMap[name] = int32(i)
        log.Println(i, name)
    }

    return true
}
func (m *PMDModel)InitParam(x int32, y int32, width int32, height int32, toonDir string){
    m.cameraPosition = mgl32.Vec3{0.0, 0.0, -15.0}
    m.ignoreCameraMotion = false
    m.rotx = 0
    m.roty = 0
    m.distance = 15.0
    m.center =  mgl32.Vec3{0.0, 10.0, 0.0}
    m.fovy = 40
    m.drawEdge = false
    m.edgeThickness = 0.004
    m.edgeColor = []float32{0.0, 0.0, 0.0, 1.0}
    m.lightDirection = mgl32.Vec3{1.0, -1.0, -1.0}
    m.lightDistance = 100.
    m.lightColor = []float32{0.6, 0.6, 0.6, 1.0}
    m.drawSelfShadow = true
    m.drawAxes = true
    m.drawCenterPoint = false
    m.fps = 30.
    m.realFps = m.fps
    m.playing = false
    m.frame =-1
    m.upPos = mgl32.Vec3{0,1,0.0}
    m.x = x
    m.y = y
    m.width = width
    m.height = height
    m.initVertices()
    m.initIndices()
    if err := m.initTextures(toonDir);err!=nil{
        log.Println(err)
    }
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
            gles2.BufferData(gles2.ARRAY_BUFFER, 4*len(tmp.array), unsafe.Pointer(&tmp.array[0]), gles2.STATIC_DRAW)
            m.vbuffers[tmp.attribute] = GLBuf{tmp.size, buffer[0]}
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
    log.Println("index count:", len(indices),"buf:",m.ibuffer)

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

func (m *PMDModel)computeMatrices(){
    m.cameraPosition = mgl32.Vec3{0.0, 0.0,-m.distance}
    m.upPos = mgl32.Vec3{0, 1,0.0}
    m.viewMatrix = mgl32.LookAtV(m.cameraPosition, m.center, m.upPos)
    gles2.Enable(gles2.CULL_FACE)
    gles2.Enable(gles2.DEPTH_TEST)
    gles2.Viewport(m.x,m.y,m.width,m.height)

    ratio := float32(m.width) / float32(m.height)
    left := -ratio
    right := ratio
    var bottom float32 = -1.0
    var top float32 = 1.0
    var near float32 = 1.0
    var far float32 = 60.0

    m.pMatrix = mgl32.Frustum(left, right, bottom, top, near, far)

}

func (m *PMDModel)DebugInitParam(x int32, y int32, width int32, height int32){
    m.vbuffers = make(map[string]GLBuf)
    if m.pmd != nil{
        length := len(m.pmd.Vertices)

        positions := make([]float32, 4*length)
        normals := make([]float32, 3*length)
        colors := make([]float32, 4*length)

        for i:= 0; i< length; i++ {
            vertex := m.pmd.Vertices[i]
            normals[3 * i] = vertex.NX
            normals[3 * i + 1] = vertex.NY
            normals[3 * i + 2] = vertex.NZ
            positions[3 * i] = vertex.X
            positions[3 * i + 1] = vertex.Y
            positions[3 * i + 2] = vertex.Z
            colors[4 * i] = 1.
            colors[4 * i + 1] = 1.
            colors[4 * i + 2] = 1.
            colors[4 * i + 3] = 1.
        }
        tmpArr := make([]AttrArrBuf, 0)
        tmpArr = append(tmpArr, AttrArrBuf{4, colors, "a_Color"})
        tmpArr = append(tmpArr, AttrArrBuf{3, positions, "a_Position"})
        tmpArr = append(tmpArr, AttrArrBuf{3, normals, "a_Normal"})
        for _, tmp := range tmpArr{
            buffer := make([]int32, 1)
            gles2.GenBuffers(1, buffer)
            gles2.BindBuffer(gles2.ARRAY_BUFFER, buffer[0])
            gles2.BufferData(gles2.ARRAY_BUFFER, 4*len(tmp.array), unsafe.Pointer(&tmp.array[0]), gles2.STATIC_DRAW)
            m.vbuffers[tmp.attribute] = GLBuf{tmp.size, buffer[0]}
        }
        gles2.BindBuffer(gles2.ARRAY_BUFFER, 0)

        indices := m.pmd.Triangles
        buffer := make([]int32, 1)
        gles2.GenBuffers(1, buffer)
        gles2.BindBuffer(gles2.ELEMENT_ARRAY_BUFFER, buffer[0])
        gles2.BufferData(gles2.ELEMENT_ARRAY_BUFFER, 4*len(indices), unsafe.Pointer(&indices[0]), gles2.STATIC_DRAW)
        m.ibuffer = buffer[0]
        log.Println("index count:", len(indices),"buf:",m.ibuffer)
        m.x = x
        m.y = y
        m.width = width
        m.height = height

    }
}

func (m *PMDModel)DebugRender(){
    gles2.ClearColor(0.5, 0.5, 0.5, 0.5)
    gles2.ClearDepthf(1)
    gles2.Enable(gles2.DEPTH_TEST)
    gles2.Enable(gles2.CULL_FACE)
    gles2.BindFramebuffer(gles2.FRAMEBUFFER, 0)
    gles2.Clear(gles2.COLOR_BUFFER_BIT | gles2.DEPTH_BUFFER_BIT)
    m.distance = 15.0
    m.cameraPosition = mgl32.Vec3{0.0, 0.0,-m.distance}
    m.center =  mgl32.Vec3{0.0, 10.0, 0.0}
    m.upPos = mgl32.Vec3{0, 1,0.0}
    m.viewMatrix = mgl32.LookAtV(m.cameraPosition, m.center, m.upPos)
    log.Println("viewMatrix", m.viewMatrix)
    gles2.Viewport(m.x,m.y,m.width,m.height)

    ratio := float32(m.width) / float32(m.height)
    left := -ratio
    right := ratio
    var bottom float32 = -1.0
    var top float32 = 1.0
    var near float32 = 1.0
    var far float32 = 60.0

    m.pMatrix = mgl32.Frustum(left, right, bottom, top, near, far)
    log.Println("pMatrix", m.pMatrix)
    m.modelMatrix = mgl32.Ident4()
    log.Println("modelMatrix", m.modelMatrix)
    m.mvMatrix = m.viewMatrix.Mul4(m.modelMatrix)

    mvpMatrix := m.pMatrix.Mul4(m.mvMatrix)

    lightPosition := mgl32.Vec3{50,50,50}

    gles2.Uniform3fv(m.programMap["u_LightPos"], 1, &lightPosition[0])
    log.Println("u_LightPos", lightPosition)
    gles2.UniformMatrix4fv(m.programMap["u_MVPMatrix"], 1, byte(0), &mvpMatrix[0])
    log.Println("u_MVPMatrix", mvpMatrix)
    gles2.UniformMatrix4fv(m.programMap["u_MVMatrix"], 1, byte(0), &m.mvMatrix[0])
    log.Println("u_MVMatrix", m.mvMatrix)

    for attr, vb := range m.vbuffers{
        gles2.BindBuffer(gles2.ARRAY_BUFFER, vb.buffer)
        gles2.VertexAttribPointer(m.programMap[attr], int32(vb.size), gles2.FLOAT, byte(0), 0, nil)
        log.Println("attr", attr)
        gles2.EnableVertexAttribArray(m.programMap[attr])
    }
    gles2.BindBuffer(gles2.ELEMENT_ARRAY_BUFFER, m.ibuffer)
    index := 0
    gles2.Enable(gles2.CULL_FACE)
    gles2.Enable(gles2.BLEND)
    gles2.DrawElements(gles2.TRIANGLES, int32(len(m.pmd.Triangles)), gles2.UNSIGNED_INT, unsafe.Pointer(&index))
    gles2.Disable(gles2.CULL_FACE)
    gles2.Disable(gles2.BLEND)
}

func (m *PMDModel)Render(){
    m.computeMatrices()
    gles2.ClearColor(0.5, 0.5, 0.5, 1)
    gles2.ClearDepthf(1)
    gles2.Enable(gles2.DEPTH_TEST)

    gles2.BindFramebuffer(gles2.FRAMEBUFFER, 0)
    gles2.Viewport(m.x, m.y, m.width, m.height)
    gles2.Clear(gles2.COLOR_BUFFER_BIT | gles2.DEPTH_BUFFER_BIT)
    for attr, vb := range m.vbuffers{
        gles2.BindBuffer(gles2.ARRAY_BUFFER, vb.buffer)
        gles2.VertexAttribPointer(m.programMap[attr], int32(vb.size), gles2.FLOAT, byte(0), 0, nil)
    }
    gles2.BindBuffer(gles2.ELEMENT_ARRAY_BUFFER, m.ibuffer)
    m.setSelfShadowTexture()
    m.setUniforms()
    gles2.Enable(gles2.CULL_FACE)
    gles2.Enable(gles2.BLEND)
    gles2.BlendFuncSeparate(gles2.SRC_ALPHA, gles2.ONE_MINUS_SRC_ALPHA, gles2.SRC_ALPHA, gles2.DST_ALPHA)
    offset := 0
    materials := m.pmd.Materials
    for _,material := range materials{
        m.renderMaterial(material, offset)
        offset += material.FaceVertCount
    }
    gles2.Disable(gles2.BLEND)
    offset = 0
    for _,material := range materials{
        m.renderEdge(material, offset)
        offset += material.FaceVertCount
    }

    gles2.Disable(gles2.CULL_FACE)
    gles2.Flush()
}

func (m *PMDModel)setSelfShadowTexture()  {
    
}
func (m *PMDModel)setUniforms()  {
    m.modelMatrix = mgl32.Ident4()
    m.mvMatrix = m.viewMatrix.Mul4(m.modelMatrix)
    m.nMatrix = m.mvMatrix.Inv()
    m.nMatrix = m.nMatrix.Transpose()

    gles2.Uniform1f(m.programMap["uEdgeThickness"], m.edgeThickness)
    gles2.Uniform3fv(m.programMap["uEdgeColor"], 1, &m.edgeColor[0])
    gles2.UniformMatrix4fv(m.programMap["uMVMatrix"], 1, byte(0), &m.mvMatrix[0])
    gles2.UniformMatrix4fv(m.programMap["uPMatrix"], 1, byte(0), &m.pMatrix[0])
    gles2.UniformMatrix4fv(m.programMap["uNMatrix"], 1, byte(0), &m.nMatrix[0])

    ld := m.lightDirection.Normalize()
    ld4 := m.nMatrix.Mul4x1(ld.Vec4(0))
    ld = ld4.Vec3()
    gles2.Uniform3fv(m.programMap["uLightDirection"], 1, &ld[0])
    gles2.Uniform3fv(m.programMap["uLightColor"], 1, &m.lightColor[0])

    gles2.Uniform1i(m.programMap["uSelfShadow"], 0)

    gles2.Uniform1i(m.programMap["uGenerateShadowMap"], 0)
    gles2.Uniform1i(m.programMap["uAxis"], 0)
    gles2.Uniform1i(m.programMap["uCenterPoint"], 0)
    

}

func (m *PMDModel)renderMaterial(material *util.PMDMaterial, offset int)  {
    gles2.Uniform3fv(m.programMap["uAmbientColor"],1, &material.Ambient[0])
    gles2.Uniform3fv(m.programMap["uSpecularColor"],1, &material.Specular[0])
    gles2.Uniform3fv(m.programMap["uDiffuseColor"], 1, &material.Diffuse[0])
    gles2.Uniform1f(m.programMap["uAlpha"], material.Alpha)
    gles2.Uniform1f(m.programMap["uShininess"], material.Shininess)
    gles2.Uniform1i(m.programMap["uEdge"], 0)
    textures := material.Textures
    gles2.ActiveTexture(gles2.TEXTURE0)
    gles2.BindTexture(gles2.TEXTURE_2D, textures["toon"])
    gles2.Uniform1i(m.programMap["uToon"], 0)
    if _,ok := textures["regular"];ok{
        gles2.ActiveTexture(gles2.TEXTURE1)
        gles2.BindTexture(gles2.TEXTURE_2D, textures["regular"])
        gles2.Uniform1i(m.programMap["uTexture"], 1)
    }
    if _, ok := textures["regular"];ok{
        gles2.Uniform1i(m.programMap["uUseTexture"], 1)
    }else{
        gles2.Uniform1i(m.programMap["uUseTexture"], 0)
    }

    _, sph_ok := textures["sph"]
    _, spa_ok := textures["spa"]
    if sph_ok||spa_ok {
        gles2.ActiveTexture(gles2.TEXTURE2)
        if sph_ok{
            gles2.BindTexture(gles2.TEXTURE_2D, textures["sph"])
            gles2.Uniform1i(m.programMap["uIsSphereMapAdditive"], 0)
        }else{
            gles2.BindTexture(gles2.TEXTURE_2D, textures["spa"])
            gles2.Uniform1i(m.programMap["uIsSphereMapAdditive"], 1)
        }

        gles2.Uniform1i(m.programMap["uSphereMap"], 2)
        gles2.Uniform1i(m.programMap["uUseSphereMap"], 1)
    } else {
        gles2.Uniform1i(m.programMap["uUseSphereMap"], 0)
    }
    gles2.CullFace(gles2.FRONT)
    index := offset * 4
    gles2.DrawElements(gles2.TRIANGLES, int32(material.FaceVertCount), gles2.UNSIGNED_INT, unsafe.Pointer(&index))
}

func (m *PMDModel)renderEdge(material *util.PMDMaterial, offset int)  {
    if (!m.drawEdge || material.EdgeFlag==0) {
        return
    }
    index := offset * 4
    gles2.Uniform1i(m.programMap["uEdge"], 1)
    gles2.CullFace(gles2.BACK)
    gles2.DrawElements(gles2.TRIANGLES, int32(material.FaceVertCount), gles2.UNSIGNED_INT, unsafe.Pointer(&index))
    gles2.CullFace(gles2.FRONT);
    gles2.Uniform1i(m.programMap["uEdge"], 0)
}
