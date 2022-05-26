package main

import (
	"log"
	"time"
)

var (
	numMaker    int        //寿司师傅总数
	numCustomer int        //当前顾客总数
	maxSushiNum int        //传送带上的寿司最大数
	maxMaterial int        //原材料最多能做的寿司数
	m           []maker    //所有师傅
	c           []customer //所有顾客
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

func makeSushi(m maker, sushiChan chan int, exitMakerChan chan bool, mChan chan bool, exitChan chan bool) {
	for i := 0; i < m.maxMake; i++ {
		_, ok := <-mChan //查看寿司材料是否用完
		if ok {          //如果材料未用完，则制作一个寿司
			time.Sleep(time.Second * time.Duration(m.timeMake))
			sushiChan <- m.timeMake
			log.Printf("%s制作出一个寿司\n", m.name)
		} else { //如果材料已用完，触发结束条件
			log.Println("寿司材料已用尽")
			exitChan <- true //向exitChan写入true，标志退出
			goto label1
		}
	}
label1:
	exitMakerChan <- true
}
func eatSushi(c customer, sushiChan chan int, exitCustomerChan chan bool, exitChan chan bool) {
	for i := 0; i < c.maxEat; i++ {
		_, ok := <-sushiChan //验证是否仍有寿司仍在生产，如果有则返回true并等待生产，没有则返回false
		if ok {              //等待生产寿司之后，则该顾客可以消费
			time.Sleep(time.Second * time.Duration(c.timeEat))
			log.Printf("%d号顾客消费一个寿司\n", c.id)
		} else { //如果没有寿司则代表所有师傅已下班
			log.Println("所有师傅已下班")
			exitChan <- true //标志退出
			goto label2
		}
	}
label2:
	exitCustomerChan <- true
}
func init() { //初始化
	numMaker = 5
	numCustomer = 5
	maxSushiNum = 10
	maxMaterial = 30
	//初始化师傅和顾客信息
	m = []maker{{"师傅1", 1, 10}, {"师傅2", 2, 9},
		{"师傅3", 3, 8}, {"师傅4", 4, 7},
		{"师傅5", 5, 6}}
	c = []customer{{1, 5, 1}, {2, 4, 2}, {3, 3, 3},
		{4, 2, 4}, {5, 1, 5}}
}

func main() {
	sushiChan := make(chan int, maxSushiNum)         //模拟传送带的channel
	exitMakerChan := make(chan bool, numMaker)       //每有一位师傅完成今日工作量会向该channel写入true
	exitCustomerChan := make(chan bool, numCustomer) //每有一位顾客吃饱会向该channel写入true
	materialChan := make(chan bool, maxMaterial)     //将该channel初始化为true，每制作一个寿司，读取一个true
	exitChan := make(chan bool, 1)                   //该channel用于监听是否有goroutine满足退出条件
	for i := 0; i < maxMaterial; i++ {
		materialChan <- true //初始化为true
	}
	close(materialChan)
	for i := 0; i < numMaker; i++ { //为每位师傅开一个goroutine生产寿司
		go makeSushi(m[i], sushiChan, exitMakerChan, materialChan, exitChan)
	}
	for i := 0; i < numCustomer; i++ { //为每个顾客开一个goroutine消费寿司
		go eatSushi(c[i], sushiChan, exitCustomerChan, exitChan)
	}
	go func() { //当所有师傅都完成工作量之后，关闭模拟传送带的channel，以便于顾客goroutine判断是否仍有寿司在生产
		for i := 0; i < numMaker; i++ {
			<-exitMakerChan
		}
		close(sushiChan)
	}()
	go func() { //当所有顾客都吃饱后离开，触发结束条件
		for i := 0; i < numCustomer; i++ {
			<-exitCustomerChan
		}
		log.Println("所有顾客已离去")
		exitChan <- true //向exitChan写入true，标志结束
		close(exitCustomerChan)
	}()
	for { //当exitChan被传入true时，代表已触发结束条件
		_, ok := <-exitChan
		if ok {
			break
		}
	}
}
