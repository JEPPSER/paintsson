package main

import (
	"strings"
	"strconv"
	"fmt"
)

func parse(str string) {
	parts := strings.Split(str, " ")
	if len(parts) == 1 {

	} else if len(parts) == 2 {
		if parts[0] == "brush-size" {
			size, err := strconv.ParseInt(parts[1], 10, 32)
			if err != nil {
				fmt.Println(err)
				return
			}
			b.rect.W = int32(size)
			b.rect.H = int32(size)
		} else if parts[0] == "brush-color" {
			b.color = colors[parts[1]]
		} else if parts[0] == "clear" {
			clearBuffer(buffer, colors[parts[1]])
		}
	}
}