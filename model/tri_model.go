package model

import (
    "log"
    "github.com/zdglf/mmd/gles2"
    "github.com/go-gl/mathgl/mgl32"
    "unsafe"
)

type TriAngleModel struct {
    program int32
    programMap map[string]int32

    triangle2VerticesData []float32
    triangle2ColorData []float32

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

func (m *TriAngleModel)LoadFile(filePath string, fileName string) bool  {
    return true;
}

func (m *TriAngleModel)InitShader(vShader string, fShader string) bool {
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
func (m *TriAngleModel)InitParam(x int32, y int32, width int32, height int32, toonDir string){

    gles2.ClearColor(0.5, 0.5, 0.5, 0.5)
    m.x = x
    m.y = y
    m.height = height
    m.width  = width

    m.cameraPosition = mgl32.Vec3{0.0, 0.0,1.5}
    m.center =  mgl32.Vec3{0.0, 0.0, -5.0}
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
    var far float32 = 10.0

    m.pMatrix = mgl32.Frustum(left, right, bottom, top, near, far)
    log.Println("pMatrix", m.pMatrix)

    m.triangle2VerticesData = []float32{
        // X, Y, Z, 
        // R, G, B, A
        -0.5, -0.25, 0.0,

        0.5, -0.25, 0.0,


        0.0, 0.559016994, 0.0,

    }
    m.triangle2ColorData = []float32{
        1.0, 1.0, 0.0, 1.0,
        0.0, 1.0, 1.0, 1.0,
        1.0, 0.0, 1.0, 1.0,
    }




}
func (m *TriAngleModel)LoadMotion(filePath string) bool{
    return true
}
func (m *TriAngleModel)getFrameCount() int{
    return 0;
}
func (m *TriAngleModel)InitFrame(index int){

}
func (m *TriAngleModel)Render (){
    gles2.Clear(gles2.COLOR_BUFFER_BIT | gles2.DEPTH_BUFFER_BIT)

    m.modelMatrix = mgl32.Translate3D(0.0, 0.0, -5.0)

    log.Println("modelMatrix", m.modelMatrix)
    m.mvMatrix = m.viewMatrix.Mul4(m.modelMatrix)

    mvpMatrix := m.pMatrix.Mul4(m.mvMatrix)

    attr := "a_Position"
    size := 3
    gles2.VertexAttribPointer(m.programMap[attr], int32(size), gles2.FLOAT, byte(0), 0, unsafe.Pointer(&m.triangle2VerticesData[0]))
    log.Println("attr", attr)
    gles2.EnableVertexAttribArray(m.programMap[attr])

    attr = "a_Color"
    size = 4
    gles2.VertexAttribPointer(m.programMap[attr], int32(size), gles2.FLOAT, byte(0), 0, unsafe.Pointer(&m.triangle2ColorData[0]))
    log.Println("attr", attr)
    gles2.EnableVertexAttribArray(m.programMap[attr])

    gles2.UniformMatrix4fv(m.programMap["u_MVPMatrix"], 1, byte(0), &mvpMatrix[0])
    gles2.DrawArrays(gles2.TRIANGLES, 0, 3)
}

