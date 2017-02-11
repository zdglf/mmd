package model

import (
	"log"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/zdglf/mmd/util"
	"github.com/zdglf/mmd/gles2"
	"path"
	"unsafe"
)

type PMXModel struct {

	program int32
	programMap map[string]int32
	pmx *util.PMX
	vbuffers map[string]GLBuf

	textureManager *TextureManager


	normals []float32
	uvs []float32
	positions []float32
	morphVec []float32
	edge []float32


	x int32
	y int32
	width int32
	height int32

	distance float32
	center mgl32.Vec3
	cameraPosition mgl32.Vec3
	upPos mgl32.Vec3
	edgeThickness float32
	edgeColor []float32
	lightDirection mgl32.Vec3
	lightColor []float32
	drawEdge bool

	viewMatrix mgl32.Mat4
	pMatrix mgl32.Mat4
	modelMatrix mgl32.Mat4
	mvMatrix mgl32.Mat4
	nMatrix mgl32.Mat4

}

func (m *PMXModel)LoadFile(filePath string, fileName string) bool  {
	m.pmx = new(util.PMX)
	log.Println("start pmd load")
	if err := m.pmx.Load(filePath, fileName); err != nil{
		m.pmx = nil
		log.Println(err)
		return false
	}else {
		log.Println(m.pmx.Name)
		log.Println(m.pmx.Comment)
		return true
	}
}

func (m *PMXModel)InitShader(vShader string, fShader string) bool {
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
	log.Println("InitShader:", gles2.GetError())

	return true
}
func (m *PMXModel)InitParam(x int32, y int32, width int32, height int32, toonDir string){
	m.distance = 15.0
	m.cameraPosition = mgl32.Vec3{0.0, 0.0, -15.0}
	m.upPos = mgl32.Vec3{0,1,0.0}
	m.center =  mgl32.Vec3{0.0, 10.0, 0.0}


	m.drawEdge = false
	m.edgeThickness = 0.004
	m.edgeColor = []float32{0.0, 0.0, 0.0, 1.0}
	m.lightDirection = mgl32.Vec3{1.0, -1.0, -1.0}
	m.lightColor = []float32{0.6, 0.6, 0.6, 1.0}


	m.x = x
	m.y = y
	m.width = width
	m.height = height


	m.initVertices()
	if err := m.initTextures(toonDir);err!=nil{
		log.Println(err)
	}
	log.Println("InitParam:", gles2.GetError())

}

func (m *PMXModel)initVertices(){
	m.vbuffers = make(map[string]GLBuf)
	if m.pmx != nil{
		length := len(m.pmx.Vertices)
		m.positions = make([]float32, 3*length)
		m.morphVec = make([]float32, 3*length)
		m.normals = make([]float32, 3*length)
		m.uvs = make([]float32, 2*length)
		m.edge = make([]float32, length)

		for i:= 0; i< length; i++ {
			vertex := m.pmx.Vertices[i]
			m.positions[ 3 * i] = vertex.X
			m.positions[ 3 * i + 1] = vertex.Y
			m.positions[ 3 * i + 2] = vertex.Z
			m.normals[3 * i] = vertex.NX
			m.normals[3 * i + 1] = vertex.NY
			m.normals[3 * i + 2] = vertex.NZ
			m.uvs[2 * i] = vertex.U
			m.uvs[2 * i + 1] = vertex.V

		}
		tmpArr := make([]AttrArrBuf, 0)
		tmpArr = append(tmpArr, AttrArrBuf{3, m.positions, "aPosition"})
		tmpArr = append(tmpArr, AttrArrBuf{3, m.morphVec, "aMultiPurposeVector"})
		tmpArr = append(tmpArr, AttrArrBuf{3, m.normals, "aVertexNormal"})
		tmpArr = append(tmpArr, AttrArrBuf{2, m.uvs, "aTextureCoord"})
		tmpArr = append(tmpArr, AttrArrBuf{1, m.edge, "aVertexEdge"})
		for _, tmp := range tmpArr{
			buffer := make([]int32, 1)
			gles2.GenBuffers(1, buffer)
			gles2.BindBuffer(gles2.ARRAY_BUFFER, buffer[0])
			gles2.BufferData(gles2.ARRAY_BUFFER, 4*len(tmp.array), unsafe.Pointer(&tmp.array[0]), gles2.STATIC_DRAW)
			m.vbuffers[tmp.attribute] = GLBuf{tmp.size, buffer[0]}

		}
		gles2.BindBuffer(gles2.ARRAY_BUFFER, 0)
		log.Println("initVertices:", gles2.GetError())

	}
}

func (m *PMXModel)initTextures(toonDir string)(err error){
	m.textureManager = NewTextureManager()
	toonFiles := []string{"toon00.bmp", "toon01.bmp", "toon02.bmp", "toon03.bmp", "toon04.bmp", "toon05.bmp",
		"toon06.bmp", "toon07.bmp", "toon08.bmp", "toon09.bmp", "toon10.bmp"}
	materials := m.pmx.Materials
	for _, material := range materials{
		if material.Textures == nil{
			material.Textures = make(map[string]int32)
		}

		log.Println("enviroment mode ", material.EnvironmentBlendMode)
		if(material.EnvironmentIndex == -1){
			material.EnvironmentBlendMode = util.MATERIAL_MODE_DISABLE
		}
		switch material.EnvironmentBlendMode {
		case util.MATERIAL_MODE_SPH:
			log.Println("enviroment ", material.EnvironmentIndex)
			if material.Textures["sph"], err = m.textureManager.Get("sph", path.Join(m.pmx.Directory, m.pmx.Textures[material.EnvironmentIndex]));err != nil{
				return
			}

		case util.MATERIAL_MODE_SPA:
			log.Println("enviroment ", material.EnvironmentIndex)
			if material.Textures["spa"], err = m.textureManager.Get("spa", path.Join(m.pmx.Directory, m.pmx.Textures[material.EnvironmentIndex]));err != nil{
				return
			}

		case util.MATERIAL_MODE_SPS:
			log.Println("enviroment ", material.EnvironmentIndex)
			if material.Textures["sps"], err = m.textureManager.Get("sps", path.Join(m.pmx.Directory, m.pmx.Textures[material.EnvironmentIndex]));err != nil{
				return
			}

		}
		if(material.ToonValue == -1){
			material.ToonValue = 0
		}
		switch material.ToonReference {
		case util.MATERIAL_TEXTURE_REFERENCE:
			log.Println("toon", material.ToonValue)
			if material.Textures["toon"], err = m.textureManager.Get("toon", path.Join(m.pmx.Directory, m.pmx.Textures[material.ToonValue]));err != nil{
				return
			}
		case util.MATERIAL_INTERNAL_REFERENCE:
			log.Println("toon", material.ToonValue)
			if material.Textures["toon"], err = m.textureManager.Get("toon", path.Join(toonDir, toonFiles[material.ToonValue]));err != nil{
				return
			}
		}
		if material.TextureIndex != -1{
			if material.Textures["regular"], err = m.textureManager.Get("regular", path.Join(m.pmx.Directory, m.pmx.Textures[material.TextureIndex]));err != nil{
				return
			}
		}

	}

	return
}

func (m *PMXModel)LoadMotion(filePath string) bool{
	return true
}
func (m *PMXModel)getFrameCount() int{
	return 0;
}
func (m *PMXModel)InitFrame(index int){

}
func (m *PMXModel)Render (){

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

	m.setUniforms()
	gles2.Enable(gles2.CULL_FACE)
	gles2.Enable(gles2.BLEND)
	gles2.BlendFuncSeparate(gles2.SRC_ALPHA, gles2.ONE_MINUS_SRC_ALPHA, gles2.SRC_ALPHA, gles2.DST_ALPHA)
	offset := 0
	materials := m.pmx.Materials
	for _,material := range materials{
		m.renderMaterial(material, offset)
		offset += material.SurfaceCount
	}
	gles2.Disable(gles2.BLEND)
	offset = 0
	for _,material := range materials{
		m.renderEdge(material, offset)
		offset += material.SurfaceCount
	}

	gles2.Disable(gles2.CULL_FACE)
	gles2.Flush()

}

func (m *PMXModel)setUniforms()  {
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
	log.Println("setUniforms:", gles2.GetError())


}

func (m *PMXModel)renderMaterial(material *util.PMXMaterial, offset int)  {
	gles2.Uniform3fv(m.programMap["uAmbientColor"],1, &material.AmbientColor[0])
	gles2.Uniform3fv(m.programMap["uSpecularColor"],1, &material.SpecularColor[0])
	gles2.Uniform3fv(m.programMap["uDiffuseColor"], 1, &material.DiffuseColor[0])
	gles2.Uniform1f(m.programMap["uAlpha"], material.Alpha)
	gles2.Uniform1f(m.programMap["uShininess"], material.SpecularStrength)
	gles2.Uniform1i(m.programMap["uEdge"], 0)
	textures := material.Textures
	gles2.ActiveTexture(gles2.TEXTURE0)
	gles2.BindTexture(gles2.TEXTURE_2D, textures["toon"])
	gles2.Uniform1i(m.programMap["uToon"], 0)
	if _,ok := textures["regular"];ok{
		gles2.ActiveTexture(gles2.TEXTURE1)
		gles2.BindTexture(gles2.TEXTURE_2D, textures["regular"])
		gles2.Uniform1i(m.programMap["uTexture"], 1)
		gles2.Uniform1i(m.programMap["uUseTexture"], 1)

	}else{
		gles2.Uniform1i(m.programMap["uUseTexture"], 0)
	}
	gles2.Uniform1i(m.programMap["uSphereMapMode"], int32(material.EnvironmentBlendMode))
	switch material.EnvironmentBlendMode {
	case util.MATERIAL_MODE_SPH:
		gles2.ActiveTexture(gles2.TEXTURE2)
		gles2.BindTexture(gles2.TEXTURE_2D, textures["sph"])
		gles2.Uniform1i(m.programMap["uSphereMap"], 2)
	case util.MATERIAL_MODE_SPA:
		gles2.ActiveTexture(gles2.TEXTURE2)
		gles2.BindTexture(gles2.TEXTURE_2D, textures["spa"])
		gles2.Uniform1i(m.programMap["uSphereMap"], 2)
	case util.MATERIAL_MODE_SPS:
		gles2.ActiveTexture(gles2.TEXTURE2)
		gles2.BindTexture(gles2.TEXTURE_2D, textures["sps"])
		gles2.Uniform1i(m.programMap["uSphereMap"], 2)

	}

	gles2.CullFace(gles2.FRONT)
	gles2.DrawElements(gles2.TRIANGLES, int32(material.SurfaceCount), gles2.UNSIGNED_INT, unsafe.Pointer(&m.pmx.Triangles[offset]))

}

func (m *PMXModel)renderEdge(material *util.PMXMaterial, offset int)  {
	if (!m.drawEdge) {
		return
	}
	gles2.Uniform1i(m.programMap["uEdge"], 1)
	gles2.CullFace(gles2.BACK)
	gles2.DrawElements(gles2.TRIANGLES, int32(material.SurfaceCount), gles2.UNSIGNED_INT, unsafe.Pointer(&m.pmx.Triangles[offset]))
	gles2.CullFace(gles2.FRONT);
	gles2.Uniform1i(m.programMap["uEdge"], 0)

}

func (m *PMXModel)computeMatrices(){
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
