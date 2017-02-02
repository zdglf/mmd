package model

import (
    "github.com/zdglf/mmd/gles2"
    "image/jpeg"
    "image/png"
    "golang.org/x/image/bmp"

    "strings"
)

type TextureManager struct {
    stores map[string]uint32
}

func (tm *TextureManager)init()  {
    tm.stores = make(map[string]uint32)
}

func (tm *TextureManager)Get(t string, url string)(texture uint32, err error){
    var ok bool
    if texture, ok = tm.stores[url];ok{
        return
    }
    textures := make([]uint32, 1)
    gles2.GenTextures(1, textures)
    texture = textures[0]
    gles2.BindTexture(gles2.TEXTURE_2D, texture)

    ext := strings.Split(url, ".")
    switch strings.ToUpper(ext[len(ext)-1]) {
    case "JPG", "JPEG":
    case "PNG":
    case "BMP":

        
    }
    //gles2.TexImage2D(gles2.TEXTURE_2D, 0, gles2.RGBA, gles2.RGBA, gles2.UNSIGNED_BYTE, );
    //
    gles2.TexParameteri(gles2.TEXTURE_2D, gles2.TEXTURE_MAG_FILTER, gles2.LINEAR)
    gles2.TexParameteri(gles2.TEXTURE_2D, gles2.TEXTURE_MAG_FILTER, gles2.LINEAR_MIPMAP_LINEAR)
    if t == "toon"{
        gles2.TexParameteri(gles2.TEXTURE_2D, gles2.TEXTURE_WRAP_S, gles2.CLAMP_TO_EDGE)
        gles2.TexParameteri(gles2.TEXTURE_2D, gles2.TEXTURE_WRAP_T, gles2.CLAMP_TO_EDGE)
    }else{
        gles2.TexParameteri(gles2.TEXTURE_2D, gles2.TEXTURE_WRAP_S, gles2.REPEAT)
        gles2.TexParameteri(gles2.TEXTURE_2D, gles2.TEXTURE_WRAP_S, gles2.REPEAT)
    }
    gles2.GenerateMipmap(gles2.TEXTURE_2D)
    gles2.BindTexture(gles2.TEXTURE_2D, 0)
    tm.stores[url] = texture
    return




}

func NewTextureManager()(tm *TextureManager)  {
    tm = new(TextureManager)
    tm.init()
    return
}