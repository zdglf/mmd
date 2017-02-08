package main

import "strings"

func main() {

    a := "\\fsdf\\sdfa\\dfasd\\fasdf\\dfsa"

    b := strings.Replace(a, "\\", "/", 0)
    println(b)



}
