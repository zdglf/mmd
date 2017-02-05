package model

import (
    "log"
    "github.com/zdglf/mmd/gles2"
    "github.com/go-gl/mathgl/mgl32"
    "unsafe"
    "github.com/zdglf/mmd/util"
)

type CubeModel struct {
    program int32
    programMap map[string]int32
    pmd *util.PMD

    cubePositionData []float32
    a_position int32
    cubeColorData []float32
    a_color int32
    cubeNormalData []float32
    a_normal int32
    indices []int32
    a_index int32

    x int32
    y int32
    width int32
    height int32

    distance float32
    center mgl32.Vec3
    cameraPosition mgl32.Vec3
    upPos mgl32.Vec3

    viewMatrix mgl32.Mat4
    pMatrix mgl32.Mat4
    modelMatrix mgl32.Mat4
    mvMatrix mgl32.Mat4
    nMatrix mgl32.Mat4
}

func (m *CubeModel)LoadFile(filePath string, fileName string) bool  {
    m.pmd = new(util.PMD)
    log.Println("start pmd load")
    if err := m.pmd.Load(filePath, fileName); err != nil{
        m.pmd = nil
        log.Println(err)
        return false
    }else {
        log.Println(m.pmd.Name)
        log.Println(m.pmd.Comment)
        return true
    }
}

func (m *CubeModel)InitShader(vShader string, fShader string) bool {
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
func (m *CubeModel)InitParam(x int32, y int32, width int32, height int32, toonDir string){

    gles2.ClearColor(0.5, 0.5, 0.5, 0.5)
    gles2.Enable(gles2.CULL_FACE)

    // Enable depth testing
    gles2.Enable(gles2.DEPTH_TEST)
    gles2.Enable(gles2.BLEND)
    gles2.BlendFunc(gles2.ONE, gles2.ONE)

    m.x = x
    m.y = y
    m.height = height
    m.width  = width

    m.cameraPosition = mgl32.Vec3{5.0, 5.0,-5}
    m.center =  mgl32.Vec3{0.0, 0.0,0.0}
    m.upPos = mgl32.Vec3{0.0, 1.0,0.0}
    m.viewMatrix = mgl32.LookAtV(m.cameraPosition, m.center, m.upPos)
    log.Println("viewMatrix", m.viewMatrix)

    gles2.Viewport(m.x,m.y,m.width,m.height)

    ratio := float32(m.width) / float32(m.height)
    left := -ratio
    right := ratio
    var bottom float32 = -1.0
    var top float32 = 1.0
    var near float32 = 1.0
    var far float32 = 20.0

    m.pMatrix = mgl32.Frustum(left, right, bottom, top, near, far)
    log.Println("pMatrix", m.pMatrix)
    if m.pmd == nil {
        length := len(m.pmd.Vertices)

        m.cubePositionData = make([]float32, 3 * length)
        m.cubeNormalData = make([]float32, 3 * length)
        m.cubeColorData = make([]float32, 4 * length)

        for i := 0; i < length; i++ {
            vertex := m.pmd.Vertices[i]
            m.cubeNormalData[3 * i] = vertex.NX
            m.cubeNormalData[3 * i + 1] = vertex.NY
            m.cubeNormalData[3 * i + 2] = vertex.NZ
            m.cubePositionData[3 * i] = vertex.X
            m.cubePositionData[3 * i + 1] = vertex.Y
            m.cubePositionData[3 * i + 2] = vertex.Z
            m.cubeColorData[4 * i] = 1.
            m.cubeColorData[4 * i + 1] = 1.
            m.cubeColorData[4 * i + 2] = 0.
            m.cubeColorData[4 * i + 3] = 1.
        }
    }else {

        m.indices = []int32{
            0 , 1, 2, 3, 4, 5,
            6 , 7, 8, 9,10,11,
            12,13,14,15,16,17,
            18,19,20,21,22,23,
            24,25,26,27,28,29,
            30,31,32,33,34,35,
        }

        m.cubePositionData = []float32{
            // In OpenGL counter-clockwise winding is default. This means that when we look at a triangle,
            // if the points are counter-clockwise we are looking at the "front". If not we are looking at
            // the back. OpenGL has an optimization where all back-facing triangles are culled, since they
            // usually represent the backside of an object and aren't visible anyways.

            // Front face
            -1.0, 1.0, 1.0,
            -1.0, -1.0, 1.0,
            1.0, 1.0, 1.0,
            -1.0, -1.0, 1.0,
            1.0, -1.0, 1.0,
            1.0, 1.0, 1.0,

            // Right face
            1.0, 1.0, 1.0,
            1.0, -1.0, 1.0,
            1.0, 1.0, -1.0,
            1.0, -1.0, 1.0,
            1.0, -1.0, -1.0,
            1.0, 1.0, -1.0,

            // Back face
            1.0, 1.0, -1.0,
            1.0, -1.0, -1.0,
            -1.0, 1.0, -1.0,
            1.0, -1.0, -1.0,
            -1.0, -1.0, -1.0,
            -1.0, 1.0, -1.0,

            // Left face
            -1.0, 1.0, -1.0,
            -1.0, -1.0, -1.0,
            -1.0, 1.0, 1.0,
            -1.0, -1.0, -1.0,
            -1.0, -1.0, 1.0,
            -1.0, 1.0, 1.0,

            // Top face
            -1.0, 1.0, -1.0,
            -1.0, 1.0, 1.0,
            1.0, 1.0, -1.0,
            -1.0, 1.0, 1.0,
            1.0, 1.0, 1.0,
            1.0, 1.0, -1.0,

            // Bottom face
            1.0, -1.0, -1.0,
            1.0, -1.0, 1.0,
            -1.0, -1.0, -1.0,
            1.0, -1.0, 1.0,
            -1.0, -1.0, 1.0,
            -1.0, -1.0, -1.0,
        }

        m.cubeColorData = []float32{
            // Front face (red)
            1.0, 0.0, 0.0, 1.0,
            1.0, 0.0, 0.0, 1.0,
            1.0, 0.0, 0.0, 1.0,
            1.0, 0.0, 0.0, 1.0,
            1.0, 0.0, 0.0, 1.0,
            1.0, 0.0, 0.0, 1.0,

            // Right face (green)
            0.0, 1.0, 0.0, 1.0,
            0.0, 1.0, 0.0, 1.0,
            0.0, 1.0, 0.0, 1.0,
            0.0, 1.0, 0.0, 1.0,
            0.0, 1.0, 0.0, 1.0,
            0.0, 1.0, 0.0, 1.0,

            // Back face (blue)
            0.0, 0.0, 1.0, 1.0,
            0.0, 0.0, 1.0, 1.0,
            0.0, 0.0, 1.0, 1.0,
            0.0, 0.0, 1.0, 1.0,
            0.0, 0.0, 1.0, 1.0,
            0.0, 0.0, 1.0, 1.0,

            // Left face (yellow)
            1.0, 1.0, 0.0, 1.0,
            1.0, 1.0, 0.0, 1.0,
            1.0, 1.0, 0.0, 1.0,
            1.0, 1.0, 0.0, 1.0,
            1.0, 1.0, 0.0, 1.0,
            1.0, 1.0, 0.0, 1.0,

            // Top face (cyan)
            0.0, 1.0, 1.0, 1.0,
            0.0, 1.0, 1.0, 1.0,
            0.0, 1.0, 1.0, 1.0,
            0.0, 1.0, 1.0, 1.0,
            0.0, 1.0, 1.0, 1.0,
            0.0, 1.0, 1.0, 1.0,

            // Bottom face (magenta)
            1.0, 0.0, 1.0, 1.0,
            1.0, 0.0, 1.0, 1.0,
            1.0, 0.0, 1.0, 1.0,
            1.0, 0.0, 1.0, 1.0,
            1.0, 0.0, 1.0, 1.0,
            1.0, 0.0, 1.0, 1.0,
        }

        m.cubeNormalData = []float32{
            // Front face
            0.0, 0.0, 1.0,
            0.0, 0.0, 1.0,
            0.0, 0.0, 1.0,
            0.0, 0.0, 1.0,
            0.0, 0.0, 1.0,
            0.0, 0.0, 1.0,

            // Right face
            1.0, 0.0, 0.0,
            1.0, 0.0, 0.0,
            1.0, 0.0, 0.0,
            1.0, 0.0, 0.0,
            1.0, 0.0, 0.0,
            1.0, 0.0, 0.0,

            // Back face
            0.0, 0.0, -1.0,
            0.0, 0.0, -1.0,
            0.0, 0.0, -1.0,
            0.0, 0.0, -1.0,
            0.0, 0.0, -1.0,
            0.0, 0.0, -1.0,

            // Left face
            -1.0, 0.0, 0.0,
            -1.0, 0.0, 0.0,
            -1.0, 0.0, 0.0,
            -1.0, 0.0, 0.0,
            -1.0, 0.0, 0.0,
            -1.0, 0.0, 0.0,

            // Top face
            0.0, 1.0, 0.0,
            0.0, 1.0, 0.0,
            0.0, 1.0, 0.0,
            0.0, 1.0, 0.0,
            0.0, 1.0, 0.0,
            0.0, 1.0, 0.0,

            // Bottom face
            0.0, -1.0, 0.0,
            0.0, -1.0, 0.0,
            0.0, -1.0, 0.0,
            0.0, -1.0, 0.0,
            0.0, -1.0, 0.0,
            0.0, -1.0, 0.0,
        }
        buffer := make([]int32, 1)
        gles2.GenBuffers(1, buffer)
        gles2.BindBuffer(gles2.ARRAY_BUFFER, buffer[0])
        gles2.BufferData(gles2.ARRAY_BUFFER, 4*len(m.cubeColorData), unsafe.Pointer(&m.cubeColorData[0]), gles2.STATIC_DRAW)
        m.a_color = buffer[0]


        buffer = make([]int32, 1)
        gles2.GenBuffers(1, buffer)
        gles2.BindBuffer(gles2.ARRAY_BUFFER, buffer[0])
        gles2.BufferData(gles2.ARRAY_BUFFER, 4*len(m.cubePositionData), unsafe.Pointer(&m.cubePositionData[0]), gles2.STATIC_DRAW)
        m.a_position = buffer[0]


        buffer = make([]int32, 1)
        gles2.GenBuffers(1, buffer)
        gles2.BindBuffer(gles2.ELEMENT_ARRAY_BUFFER, buffer[0])
        gles2.BufferData(gles2.ELEMENT_ARRAY_BUFFER, 4*len(m.indices), unsafe.Pointer(&m.indices[0]), gles2.STATIC_DRAW)
        m.a_index = buffer[0]
    }



}
func (m *CubeModel)LoadMotion(filePath string) bool{
    return true
}
func (m *CubeModel)getFrameCount() int{
    return 0;
}
func (m *CubeModel)InitFrame(index int){

}
func (m *CubeModel)Render (){
    gles2.ClearColor(0.5, 0.5, 0.5, 0)
    gles2.Enable(gles2.CULL_FACE)

    gles2.Enable(gles2.DEPTH_TEST)
    gles2.Clear(gles2.COLOR_BUFFER_BIT | gles2.DEPTH_BUFFER_BIT)
    m.modelMatrix = mgl32.Ident4()
    log.Println("modelMatrix", m.modelMatrix)
    m.mvMatrix = m.viewMatrix.Mul4(m.modelMatrix)

    mvpMatrix := m.pMatrix.Mul4(m.mvMatrix)

    attr := "a_Position"
    size := 3
    gles2.BindBuffer(gles2.ARRAY_BUFFER,m.a_position)
    gles2.VertexAttribPointer(m.programMap[attr], int32(size), gles2.FLOAT, byte(0), 0, nil)
    log.Println("attr", attr)
    gles2.EnableVertexAttribArray(m.programMap[attr])

    attr = "a_Color"
    size = 4
    gles2.BindBuffer(gles2.ARRAY_BUFFER,m.a_color)
    gles2.VertexAttribPointer(m.programMap[attr], int32(size), gles2.FLOAT, byte(0), 0, nil)
    log.Println("attr", attr)
    gles2.EnableVertexAttribArray(m.programMap[attr])
    gles2.BindBuffer(gles2.ELEMENT_ARRAY_BUFFER,m.a_index)
    //attr = "a_Normal"
    //size = 3
    //gles2.VertexAttribPointer(m.programMap[attr], int32(size), gles2.FLOAT, byte(0), 0, unsafe.Pointer(&m.cubeNormalData[0]))
    //log.Println("attr", attr)
    //gles2.EnableVertexAttribArray(m.programMap[attr])

    //light := []float32{0,5.0,5.0}

    //gles2.UniformMatrix4fv(m.programMap["u_MVMatrix"], 1, byte(0), &m.mvMatrix[0])
    gles2.UniformMatrix4fv(m.programMap["u_MVPMatrix"], 1, byte(0), &mvpMatrix[0])
    //gles2.Uniform3fv(m.programMap["u_LightPos"], 1, &light[0])

    index := 0
    gles2.DrawElements(gles2.TRIANGLES, int32(36), gles2.UNSIGNED_INT, unsafe.Pointer(&index))
    //gles2.DrawArrays(gles2.TRIANGLES, 0, 36)
    e := gles2.GetError()

    var maxLength int32 = 1024
    infoBytes := make([]byte, maxLength)
    var realLength int32
    gles2.GetProgramInfoLog(m.program, maxLength, &realLength, infoBytes)
    log.Println("runtime", string(infoBytes[:realLength]), e)

}

