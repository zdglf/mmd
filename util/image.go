package util

import (
    "bytes"
    "os"
)

const(
    TYPE_PNG = "PNG"
    TYPE_GIF = "GIF"
    TYPE_JPG = "JPG"
    TYPE_TGA = "TGA"
    TYPE_BMP = "BMP"
    TYPE_UNKNOWN = "UNKNOWN"

)

/**
BMP first 2 bytes should be

    BM – Windows 3.1x, 95, NT, ... etc.
    BA – OS/2 struct bitmap array
    CI – OS/2 struct color icon
    CP – OS/2 const color pointer
    IC – OS/2 struct icon
    PT – OS/2 pointer
PNG first 8 bytes should be

    137 80 78 71 13 10 26 10
        P  N  G
GIF first 6 bytes should be

    3 bytes  "GIF"
    3 bytes  "87a" or "89a"

JPG, JPEG first 2 bytes should be

    0xFF, 0xD8

TGA, end with 18 bytes should be
    "TRUEVISION-XFILE.\x00"
 */
func GetImageType(url string)(t string, err error)  {

    TGA := []byte("TRUEVISION-XFILE.\x00")
    JPG := []byte("\xFF\xD8")
    GIF := []byte("GIF")
    PNG := []byte("\x89PNG")
    BMP := []byte("BM")
    var HEADER_SIZE int64 = 4
    var FOOTER_SIZE int64 = 18
    var fi os.FileInfo
    if fi, err = os.Stat(url); err != nil{
        return
    }
    filesize := fi.Size()
    var f *os.File
    if f, err = os.Open(url); err != nil{
        return
    }

    defer f.Close()

    header := make([]byte, HEADER_SIZE)
    footer := make([]byte, FOOTER_SIZE)
    if _, err = f.Read(header);err != nil{
        return
    }
    if bytes.Equal(PNG, header){
        t =  TYPE_PNG
        return
    }
    if bytes.Equal(BMP, header[:2]){
        t = TYPE_BMP
        return
    }
    if bytes.Equal(GIF, header[:3]){
        t = TYPE_GIF
        return
    }
    if bytes.Equal(JPG, header[:2]){
        t = TYPE_JPG
        return
    }
    if _, err = f.ReadAt(footer, filesize-FOOTER_SIZE);err != nil{
        return
    }
    if bytes.Equal(TGA, footer){
        t = TYPE_TGA
        return
    }
    t = TYPE_UNKNOWN
    return

}
