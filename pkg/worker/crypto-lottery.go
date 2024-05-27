package worker

import (
	"context"
	"crypto-lottery-srv/pkg/consumer"
	"crypto-lottery-srv/pkg/generator"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

const (
	NumOfWorkers = 16 // 协程数量

)

// 设置Etherscan API密钥
var etherscanAPIKeys = []string{
	"",
} // 替换为你的 API keys

type AppData struct {
	Logger     *log.Logger
	GoodLogger *log.Logger
}

var App = AppData{}

func Lottery() {
	// 创建工作协程
	var wg sync.WaitGroup
	wg.Add(NumOfWorkers)
	// Create a context for cancellation
	ctx, _ := context.WithCancel(context.Background())
	for i := 0; i < NumOfWorkers; i++ {
		go func(workerID int) {
			defer wg.Done()
			runWorker(ctx, workerID, etherscanAPIKeys[workerID%len(etherscanAPIKeys)])
		}(i)
	}
	wg.Wait()
	App.Log(0, "程序执行完毕")
}

func BuyMoreLotteryInAws() {
	// 创建工作协程
	var wg sync.WaitGroup
	wg.Add(NumOfWorkers)
	for i := 0; i < 16; i++ {
		go func(workerID int) {
			master := NewAwsMaster(workerID)
			defer master.Stop()
			for {
				privateKey, err := generator.GeneratePrivateKey(fmt.Sprintf("%x", workerID))
				if err != nil {
					//fmt.Println(fmt.Sprintf("生成私钥失败：%v", err))
					App.Log(workerID, fmt.Sprintf("生成私钥失败：%v", err))
				}
				address := generator.PrivateKeyToAddress(privateKey)
				master.Upload(privateKey, address)
			}
		}(i)
	}
	wg.Wait()
	App.Log(0, "程序执行完毕")
}

func runWorker(ctx context.Context, workerID int, apiKey string) {
	addressList := make([]string, 0)
	keyList := make([]string, 0)
	for i := 0; ; i++ {
		select {
		case <-ctx.Done(): // If context is cancelled, exit the goroutine
			return
		default:
			privateKey, err := generator.GeneratePrivateKey(fmt.Sprintf("%x", workerID))
			if err != nil {
				App.Log(workerID, fmt.Sprintf("生成私钥失败：%v", err))
				continue
			}
			//key := start + strconv.FormatInt(int64(i), 16)
			//privateKey, err := generator.GeneratePrivateKey(fmt.Sprintf("%x", workerID))
			if err != nil {
				App.Log(workerID, fmt.Sprintf("生成私钥失败：%v", err))
				continue
			}
			address := generator.PrivateKeyToAddress(privateKey)
			keyList = append(keyList, privateKey)
			addressList = append(addressList, address)
			if len(addressList) < 20 {
				continue
			}
			balanceResults, err := consumer.GetBalance(addressList, apiKey)
			if err != nil {
				fmt.Println(fmt.Sprintf("查询地址余额失败：%v", err))
				App.Log(workerID, fmt.Sprintf("查询地址余额失败：%v", err))
				continue
			}
			App.Log(workerID, fmt.Sprintf("[Round %d]", i))
			App.Log(workerID, fmt.Sprintf("[Keys][%s]", strings.Join(keyList, ",")))
			for _, result := range balanceResults {
				App.Log(workerID, fmt.Sprintf("[%s](%s ETH)", result.Account, result.Balance))
				if result.Balance != "0" {
					App.GoodLog(workerID, fmt.Sprintf("找到有资产的地址：%s (余额：%s ETH)", result.Account, result.Balance))
					App.GoodLog(workerID, "Oh god look at me !!!!!!!")
					App.GoodLog(workerID, fmt.Sprintf("[Keys are here][%s]", strings.Join(keyList, ",")))
					App.GoodLog(workerID, fmt.Sprintf("[%s](%s ETH)", result.Account, result.Balance))
					break
				}
			}
			addressList = addressList[:0]
			keyList = keyList[:0]
		}
	}
}

func (a AppData) Log(workerId int, content string) {
	a.Logger.Println(fmt.Sprintf("[%d]:%s", workerId, content))
}

func (a AppData) GoodLog(workerId int, content string) {
	a.Logger.Println(fmt.Sprintf("[%d]:%s", workerId, content))
}

func InitApp(recordF, goodF *os.File) {
	App.Logger = log.New(recordF, "[CL]", log.LstdFlags)
	App.GoodLogger = log.New(goodF, "[CL]", log.LstdFlags)
}
