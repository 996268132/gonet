package data

import (
	"gonet/base"
	"gonet/common"
	"sync"
)

var (
	waitGroup sync.WaitGroup
)

//异步读取ata
func ansyReadData(res common.IBaseDataRes) {
	waitGroup.Add(1)
	go func() {
		res.Read()
		waitGroup.Done()
	}()
}

func InitRepository() {
	base.GLOG.Println("----read data begin-----")
	//ansyReadData(&BANDATA)
	waitGroup.Wait()
	base.GLOG.Println("----read data end-----")
}
