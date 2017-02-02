package main

import "strings"

func main() {
    var a string
    a = "stad.png\\*dfasdf.gif\\*sdfas.tex"
    b := strings.Split(a, "\\*")
    println(b)

    println(a)



}
