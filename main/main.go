package main

import (
	"log"
	"time"
)

var (
	numMaker    int //寿司师傅总数
	numCustomer int //当前顾客总数
	maxSushiNum int //传送带上的寿司最大数
	maxMaterial int //原材料最多能做的寿司数
)

type maker struct {
	name     string //师傅的名字
	timeMake int    //师傅制作寿司的时间，单位分钟
	maxMake  int    //师傅每天最多可做的寿司数
}
type customer struct {
	id      int //顾客号
	timeEat int //顾客吃一个寿司所需时间，单位分钟
	maxEat  int //顾客最多可吃寿司数
}

func makeSushi(m maker, sushiChan chan int, exitChan chan bool) {
	for i := 0; i < m.maxMake; i++ {
		time.Sleep(time.Second * time.Duration(m.timeMake))
		sushiChan <- m.timeMake
		log.Printf("%s制作出一个寿司\n", m.name)
	}
	exitChan <- true
}
func eatSushi(c customer, sushiChan chan int, exitChan chan bool) {
	for i := 0; i < c.maxEat; i++ {
		<-sushiChan
		time.Sleep(time.Second * time.Duration(c.timeEat))
		log.Printf("%d号顾客消费一个寿司\n", c.id)
	}
	exitChan <- true
}
func init() {
	numMaker = 5
	numCustomer = 5
	maxSushiNum = 10
	maxMaterial = 30
}

func main() {
	sushiChan := make(chan int, maxSushiNum)         //模拟传送带的channel
	exitMakerChan := make(chan bool, numMaker)       //每有一位师傅完成今日工作量会向该channel写入true
	exitCustomerChan := make(chan bool, numCustomer) //每有一位顾客吃饱会向该channel写入true
	//初始化师傅和顾客信息
	m := []maker{{"师傅1", 1, 10}, {"师傅2", 2, 9},
		{"师傅3", 3, 8}, {"师傅4", 4, 7},
		{"师傅5", 5, 6}}
	c := []customer{{1, 5, 5}, {2, 5, 5}, {3, 5, 5},
		{4, 5, 5}, {5, 5, 5}}
	for i := 0; i < numMaker; i++ {
		go makeSushi(m[i], sushiChan, exitMakerChan)
	}
	for i := 0; i < numCustomer; i++ {
		go eatSushi(c[i], sushiChan, exitCustomerChan)
	}

}
