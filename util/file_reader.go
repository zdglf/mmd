package util

import (
    "os"
    "errors"
    "encoding/binary"
    "math"
    "golang.org/x/text/encoding"
    "golang.org/x/text/transform"
    "bytes"
    "io/ioutil"
)

type FileReader struct  {
    f *os.File

}

func (fr *FileReader)Set(f *os.File){
    fr.f = f
}

func (fr *FileReader)GetUIntLittle(size int)(data int,err error){
    data = 0
    if(fr.f!=nil){
        bits := make([]byte, size)
        _, err =fr.f.Read(bits)
        switch size {
        case 1:
            data = int(uint8(bits[0]))
        case 2:
            data = int(binary.LittleEndian.Uint16(bits))
        case 4:
            data = int(binary.LittleEndian.Uint32(bits))
        }

    }else {
        err = errors.New("file not set")
    }
    return data, err

}

func (fr *FileReader)GetIntLittle(size int)(data int,err error) {
    data = 0
    if(fr.f!=nil){
        bits := make([]byte, size)
        _, err =fr.f.Read(bits)
        switch size {
        case 1:
            data = int(int8(bits[0]))
        case 2:
            data = int(int16(binary.LittleEndian.Uint16(bits)))
        case 4:
            data = int(int32(binary.LittleEndian.Uint32(bits)))
        }

    }else {
        err = errors.New("file not set")
    }
    return data, err

}
func (fr *FileReader)GetInt8Little()(data int,err error){
    return fr.GetIntLittle(1)
}
func (fr *FileReader)GetInt16Little()(data int,err error){
    return fr.GetIntLittle(2)
}
func (fr *FileReader)GetInt32Little()(data int,err error){
    return fr.GetIntLittle(4)
}
func (fr *FileReader)GetUInt8Little()(data int,err error){
    return fr.GetUIntLittle(1)
}
func (fr *FileReader)GetUInt16Little()(data int, err error){
    return fr.GetUIntLittle(2)
}
func (fr *FileReader)GetUInt32Little()(data int, err error){
    return fr.GetUIntLittle(4)
}
func (fr *FileReader)GetFloatLittle()(data float32,err error){
    var n int
    n, err = fr.GetUInt32Little()
    data = math.Float32frombits(uint32(n))
    return data, err
}
func (fr *FileReader)GetString(size int, decoder *encoding.Decoder)(data string,err error){
    data = ""
    if(fr.f!=nil){
        bits := make([]byte, size)
        _, err = fr.f.Read(bits)
        rInUTF8 := transform.NewReader(bytes.NewReader(bits), decoder)
        decBytes, _ := ioutil.ReadAll(rInUTF8)
        data = string(decBytes)

    }else{
        err = errors.New("file not set")
    }
    return data, err

}

func (fr *FileReader)GetStringTrim(size int, decoder *encoding.Decoder)(data string,err error){
    data = ""
    if(fr.f!=nil){
        bits := make([]byte, size)
        _, err = fr.f.Read(bits)
        index := bytes.IndexByte(bits, 0)
        rInUTF8 := transform.NewReader(bytes.NewReader(bits[:index]), decoder)
        decBytes, _ := ioutil.ReadAll(rInUTF8)
        data = string(decBytes)

    }else{
        err = errors.New("file not set")
    }
    return data, err

}

func (fr *FileReader)GetStringUTF8(size int)(data string,err error){
    data = ""
    if(fr.f!=nil){
        bits := make([]byte, size)
        _, err = fr.f.Read(bits)
        data = string(bits)

    }else{
        err = errors.New("file not set")
    }
    return data, err

}

func (fr *FileReader)GetStringUTF8Trim(size int)(data string,err error){
    data = ""
    if(fr.f!=nil){
        bits := make([]byte, size)
        _, err = fr.f.Read(bits)
        index := bytes.IndexByte(bits, 0)
        data = string(bits[:index])

    }else{
        err = errors.New("file not set")
    }
    return data, err

}

func (fr *FileReader)GetBytes(size int)(data []byte,err error){
    if(fr.f!=nil){
        data = make([]byte, size)
        _, err = fr.f.Read(data)

    }else{
        err = errors.New("file not set")
    }
    return data, err
}


func (fr *FileReader)GetUIntBig(size int)(data int,err error){
    data = 0
    if(fr.f!=nil){
        bits := make([]byte, size)
        _, err =fr.f.Read(bits)
        switch size {
        case 1:
            data = int(uint8(bits[0]))
        case 2:
            data = int(binary.BigEndian.Uint16(bits))
        case 4:
            data = int(binary.BigEndian.Uint32(bits))
        }

    }else {
        err = errors.New("file not set")
    }
    return data, err
}

func (fr *FileReader)GetIntBig(size int)(data int,err error) {
    data = 0
    if(fr.f!=nil){
        bits := make([]byte, size)
        _, err =fr.f.Read(bits)
        switch size {
        case 1:
            data = int(int8(bits[0]))
        case 2:
            data = int(int16(binary.BigEndian.Uint16(bits)))
        case 4:
            data = int(int32(binary.BigEndian.Uint32(bits)))
        }

    }else {
        err = errors.New("file not set")
    }
    return data, err
}
func (fr *FileReader)GetInt8Big()(data int,err error){
    return fr.GetIntBig(1)
}
func (fr *FileReader)GetInt16Big()(data int,err error){
    return fr.GetIntBig(2)
}
func (fr *FileReader)GetInt32Big()(data int,err error){
    return fr.GetIntBig(4)
}
func (fr *FileReader)GetUInt8Big()(data int,err error){
    return fr.GetUIntBig(1)
}
func (fr *FileReader)GetUInt16Big()(data int, err error){
    return fr.GetUIntBig(2)
}
func (fr *FileReader)GetUInt32Big()(data int, err error){
    return fr.GetUIntBig(4)
}
func (fr *FileReader)GetFloatBig()(data float32,err error){
    var n int
    n, err = fr.GetUInt32Big()
    data = math.Float32frombits(uint32(n))
    return data, err
}

func isBitSet(data byte, count int) bool {
    if count<0 ||count>7{
        return false
    }
    return (data&(1<<uint(count))) != 0
}