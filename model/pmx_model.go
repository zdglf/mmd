package model

type PMXModel struct {

}

func (m *PMXModel)LoadFile(filePath string, fileName string) bool  {
	return true;
}

func (m *PMXModel)InitShader(vShader string, fShader string) bool {
	return true
}
func (m *PMXModel)InitParam(x int, y int, width int, height int){

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

}
