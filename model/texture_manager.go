package model

import (
    "github.com/zdglf/mmd/gles2"
    "image/gif"
    "image"
    "image/jpeg"
    "image/png"
    "golang.org/x/image/bmp"
    "github.com/zdglf/mmd/tga"

    "github.com/zdglf/mmd/util"
    "os"
    "image/draw"
    "unsafe"
)

type TextureManager struct {
    stores map[string]int32
}

func (tm *TextureManager)init()  {
    tm.stores = make(map[string]int32)
}

func (tm *TextureManager)Get(t string, url string)(texture int32, err error){
    var ok bool
    if texture, ok = tm.stores[url];ok{
        return
    }
    textures := make([]int32, 1)
    gles2.GenTextures(1, textures)
    texture = textures[0]
    gles2.BindTexture(gles2.TEXTURE_2D, texture)
    var imageType string
    if imageType, err = util.GetImageType(url); err != nil{
        return
    }
    var f *os.File
    if f, err = os.Open(url); err != nil{
        return
    }
    defer f.Close()
    var textureImage image.Image
    switch imageType {
    case util.TYPE_JPG:
        if textureImage, err = jpeg.Decode(f); err != nil{
            return
        }
    case util.TYPE_PNG:
        if textureImage, err = png.Decode(f); err != nil{
            return
        }
    case util.TYPE_BMP:
        if textureImage, err = bmp.Decode(f); err != nil{
            return
        }
    case util.TYPE_TGA:
        if textureImage, err = tga.Decode(f); err != nil{
            return
        }
    case util.TYPE_GIF:
        if textureImage, err = gif.Decode(f); err != nil{
            return
        }
    default:
        if textureImage, err = png.Decode(f); err != nil{
            return
        }
    }

    rect := textureImage.Bounds()
    rgba := image.NewRGBA(rect)
    draw.Draw(rgba, rect, textureImage, rect.Min, draw.Src)
    gles2.TexImage2D(gles2.TEXTURE_2D, 0, gles2.RGBA,int32(rect.Max.X-rect.Min.X), int32(rect.Max.Y-rect.Min.Y), 0,  gles2.RGBA, gles2.UNSIGNED_BYTE, unsafe.Pointer(&rgba.Pix[0]))

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