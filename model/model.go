package model


type GLBuf struct {
    size int
    buffer int32
}

type AttrArrBuf struct {
    size int
    array []float32
    attribute string
}