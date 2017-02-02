package util

import (
    "testing"
    "path"
    "log"
)

func TestGetImageType(t *testing.T) {
    p := path.Join("test_data", "gif.gif")
    log.Println(p)
    if tp, err := GetImageType(p); err != nil {
        t.Error(err)
    } else {
        if tp == TYPE_GIF{
            log.Println(tp, "is gif")
        }else{
            t.Error(tp, "is not gif")
        }
    }

    p = path.Join("test_data", "png.png")
    log.Println(p)
    if tp, err := GetImageType(p); err != nil {
        t.Error(err)
    } else {
        if tp == TYPE_PNG{
            log.Println(tp, "is png")
        }else{
            t.Error(tp, "is not png")
        }
    }

    p = path.Join("test_data", "tga.tga")
    log.Println(p)
    if tp, err := GetImageType(p); err != nil {
        t.Error(err)
    } else {
        if tp == TYPE_TGA{
            log.Println(tp, "is tga")
        }else{
            t.Error(tp, "is not tga")
        }
    }


    p = path.Join("test_data", "jpg.jpg")
    log.Println(p)
    if tp, err := GetImageType(p); err != nil {
        t.Error(err)
    } else {
        if tp == TYPE_JPG{
            log.Println(tp, "is jpg")
        }else{
            t.Error(tp, "is not jpg")
        }
    }

    p = path.Join("test_data", "bmp.bmp")
    log.Println(p)
    if tp, err := GetImageType(p); err != nil {
        t.Error(err)
    } else {
        if tp == TYPE_BMP{
            log.Println(tp, "is bmp")
        }else{
            t.Error(tp, "is not bmp")
        }
    }


}
