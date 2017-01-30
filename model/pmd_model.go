package model

import (
    "github.com/zdglf/mmd/util"
    "log"

)

type PMDModel struct {
    program uint32
    pmd *util.PMD
}

func (m *PMDModel)LoadFile(filePath string, fileName string) bool  {
    m.pmd = new(util.PMD)
    log.Println("start pmd load")
    if err := m.pmd.Load(filePath, fileName); err != nil{
        log.Println(err)
        return false
    }else {
        log.Println(m.pmd.Name)
        log.Println(m.pmd.Comment)
        return true;
    }
}

func (m *PMDModel)InitShader(vShader string, fShader string) bool {
    //vertexShader := gles2.CreateShader(gles2.VERTEX_SHADER)
    //gles2.ShaderSource(vertexShader,1, []string{vShader}, []int32{int32(len(vShader))})
    //gles2.CompileShader(vertexShader)
    //var compileStatus int32
    //gles2.GetShaderiv(vertexShader, gles2.COMPILE_STATUS, &compileStatus)
    //if compileStatus == 0{
    //    var maxLength int32 = 1024
    //    infoBytes := make([]byte, maxLength)
    //    var realLength int32
    //    gles2.GetShaderInfoLog(vertexShader, maxLength, &realLength, infoBytes)
    //    log.Println("vshader", string(infoBytes[:realLength]))
    //    gles2.DeleteShader(vertexShader)
    //    vertexShader = 0
    //}
    //if vertexShader == 0{
    //    return false
    //}
    //
    //fragmentShader := gles2.CreateShader(gles2.FRAGMENT_SHADER)
    //gles2.ShaderSource(fragmentShader, 1, []string{fShader}, []int32{int32(len(fShader))})
    //gles2.CompileShader(fragmentShader)
    //gles2.GetShaderiv(fragmentShader, gles2.COMPILE_STATUS, &compileStatus)
    //if compileStatus == 0{
    //    var maxLength int32 = 1024
    //    infoBytes := make([]byte, maxLength)
    //    var realLength int32
    //    gles2.GetShaderInfoLog(fragmentShader, maxLength, &realLength, infoBytes)
    //    log.Println("fshader", string(infoBytes[:realLength]))
    //    gles2.DeleteShader(fragmentShader)
    //    fragmentShader = 0
    //}
    //if fragmentShader == 0{
    //    return false
    //}
    //m.program = gles2.CreateProgram()
    //gles2.AttachShader(m.program, vertexShader)
    //gles2.AttachShader(m.program, fragmentShader)
    //gles2.LinkProgram(m.program)
    //var linkStatus int32
    //gles2.GetProgramiv(m.program, gles2.LINK_STATUS, &linkStatus)
    //if linkStatus == 0{
    //    var maxLength int32 = 1024
    //    infoBytes := make([]byte, maxLength)
    //    var realLength int32
    //    gles2.GetProgramInfoLog(m.program, maxLength, &realLength, infoBytes)
    //    log.Println("link", string(infoBytes[:realLength]))
    //    m.program = 0
    //}
    //if m.program == 0{
    //    return false
    //}
    //gles2.UseProgram(m.program)


    return true
}
func (m *PMDModel)InitParam(x int, y int, width int, height int){

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
