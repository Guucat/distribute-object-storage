package main

import "time"

func main() {
	//var m map[string]string
	//m = make(map[string]string)
	//m["a"] = "1"
	//fmt.Println(m["a"])
	names := []string{"a", "b", "c"}
	for _, name := range names {
		go func(name *string) {
			println(*name)
		}(&name)
	}
	time.Sleep(1 * time.Second)
}
