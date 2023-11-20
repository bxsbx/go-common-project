package main

import (
	"StandardProject/common/local"
	"fmt"
	"strconv"
)

var Map = make(map[string]string)
var lockMap = local.NewLockMap(5, 2)

func M(k string) {
	for i := 0; i < 10; i++ {
		lock := lockMap.GetLock(strconv.Itoa(i))
		lock.Lock()
		fmt.Println(k, i, lock)
		Map[fmt.Sprintf("%v", lock)] = "lock"
		lock.Unlock()
	}
}

func main() {
	a := make([][]int, 2)
	a[0] = make([]int, 0)
	fmt.Println(a)
}
