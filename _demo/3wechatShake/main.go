/*
微信摇一摇
压力测试： wrk -t10 -c10 -d5 http://localhost:8080/lucky (t:线程数，c:tcp连接数 d:压测时间)
*/

package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

// 奖品类型 枚举类型
const (
	gitftTypeCoin     = iota //虚拟币
	giftTypeCoupon           // 不同券
	giftTypeCouponFix        // 相同的券
	giftTypeRealSmall        // 实物小奖
	giftTypeRealLarge        // 实物大奖
)
const (
	rateMax = 10000 // 最大中奖号码
)

var (
	logger   *log.Logger
	giftList []*gift
	mu       sync.Mutex
)

type gift struct {
	id       int      // 奖品ID
	name     string   // 奖品名称
	pic      string   // 奖品图片
	link     string   //奖品链接
	gtype    int      // 奖品类型
	data     string   // 奖品数据（特定配置信息）
	datalist []string // 奖品数据集合 （不同优惠券的编码）
	total    int      // 总数，0不限量
	left     int      // 剩余数量
	inuse    bool     // 是否使用中
	rate     int      // 中奖概率 万分之N 0-9999
	rateMin  int      // 大于等于最小中奖编码
	rateMax  int      // 小于中奖编码
}

type lotterController struct {
	Ctx iris.Context
}

func newApp() *iris.Application {
	app := iris.New()
	initLog()
	initGift()
	mvc.New(app.Party("/")).Handle(&lotterController{})
	return app
}

func initLog() {
	f, _ := os.Create("/Users/kongli/gocode/src/lottery/info.log")
	logger = log.New(f, "", log.Ldate|log.Lmicroseconds)
}

func initGift() {
	giftList = make([]*gift, 5)
	g1 := gift{
		id:       1,
		name:     "手机大奖",
		pic:      "",
		link:     "",
		gtype:    giftTypeRealLarge,
		data:     "",
		datalist: nil,
		total:    100000,
		left:     100000,
		inuse:    true,
		rate:     10000,
		rateMin:  0,
		rateMax:  0,
	}
	g2 := gift{
		id:       2,
		name:     "充电器",
		pic:      "",
		link:     "",
		gtype:    giftTypeRealSmall,
		data:     "",
		datalist: nil,
		total:    5,
		left:     5,
		inuse:    false,
		rate:     10,
		rateMin:  0,
		rateMax:  0,
	}
	g3 := gift{
		id:       3,
		name:     "优惠券满200减50",
		pic:      "",
		link:     "",
		gtype:    giftTypeCouponFix,
		data:     "mall-coupon-2018",
		datalist: nil,
		total:    5,
		left:     5,
		inuse:    false,
		rate:     5000,
		rateMin:  0,
		rateMax:  0,
	}
	g4 := gift{
		id:       4,
		name:     "直减50优惠券",
		pic:      "",
		link:     "",
		gtype:    giftTypeCoupon,
		data:     "",
		datalist: []string{"c01", "c02", "c03", "c04", "c05"},
		total:    5,
		left:     5,
		inuse:    false,
		rate:     5000,
		rateMin:  0,
		rateMax:  0,
	}
	g5 := gift{
		id:       5,
		name:     "比特币",
		pic:      "",
		link:     "",
		gtype:    gitftTypeCoin,
		data:     "0.01个比特币",
		datalist: nil,
		total:    5,
		left:     5,
		inuse:    false,
		rate:     5000,
		rateMin:  0,
		rateMax:  0,
	}
	giftList[0] = &g1
	giftList[1] = &g2
	giftList[2] = &g3
	giftList[3] = &g4
	giftList[4] = &g5

	// 数据整理，中奖区间数据

	rateStart := 0

	for _, data := range giftList {
		if !data.inuse {
			continue
		}
		data.rateMin = rateStart
		data.rateMax = rateStart + data.rate

		if data.rateMax >= rateMax {
			data.rateMax = rateMax
			rateStart = 0
		} else {
			rateStart += data.rate
		}
	}

}

func main() {
	app := newApp()
	app.Run(iris.Addr(":8080"))
}

func (c *lotterController) Get() string {
	count := 0
	total := 0
	for _, data := range giftList {
		if data.inuse && (data.total == 0) || (data.total > 0 && data.left > 0) {
			count++
			total += data.left
		}
	}
	return fmt.Sprintf("当前有效奖品种类：%d,限量奖品数量：%d", count, total)
}

// http://localhost:8080/lucky
func (c *lotterController) GetLucky() map[string]interface{} {
	mu.Lock()
	defer mu.Unlock()
	code := luckyCode()
	var ok bool
	result := make(map[string]interface{})
	result["success"] = ok
	for _, data := range giftList {
		if !data.inuse || (data.total > 0 && data.left <= 0) {
			continue
		}
		if data.rateMin <= int(code) && data.rateMax > int(code) {
			// 中奖，抽奖编码在奖品编码范围内
			// 开始发奖
			sendData := ""
			switch data.gtype {
			case gitftTypeCoin:
				ok, sendData = sendCoin(data)
			case giftTypeCoupon:
				ok, sendData = sendCoupon(data)
			case giftTypeCouponFix:
				ok, sendData = sendCouponFix(data)
			case giftTypeRealSmall:
				ok, sendData = sendRealSmall(data)
			case giftTypeRealLarge:
				ok, sendData = sendRealLarge(data)
			}
			if ok {
				// 记录中奖结果
				saveLuckData(code, data, sendData)
				result["success"] = ok
				result["id"] = data.id
				result["name"] = data.name
				result["link"] = data.link
				result["data"] = data.data
				break
			}
		}
	}
	return result
}

func sendCoin(data *gift) (bool, string) {
	if data.total == 0 {
		// 数量无限
		return true, data.data
	} else if data.left > 0 {
		data.left = data.left - 1
		return true, data.data
	} else {
		return false, "奖品已发完"
	}
}

func sendCoupon(data *gift) (bool, string) {
	if len(data.datalist) < data.left {
		data.left = len(data.datalist)
	}
	if data.left > 0 {
		left := data.left - 1
		data.left = left
		return true, data.datalist[left]
	} else {
		return false, "奖品已发完"
	}
}

func sendCouponFix(data *gift) (bool, string) {
	if data.total == 0 {
		// 数量无限
		return true, data.data
	} else if data.left > 0 {
		data.left = data.left - 1
		return true, data.data
	} else {
		return false, "奖品已发完"
	}
}
func sendRealSmall(data *gift) (bool, string) {
	if data.total == 0 {
		// 数量无限
		return true, data.data
	} else if data.left > 0 {
		data.left = data.left - 1
		return true, data.data
	} else {
		return false, "奖品已发完"
	}
}
func sendRealLarge(data *gift) (bool, string) {
	if data.total == 0 {
		// 数量无限
		return true, data.data
	} else if data.left > 0 {
		data.left = data.left - 1
		return true, data.data
	} else {
		return false, "奖品已发完"
	}
}

func luckyCode() int32 {
	seed := time.Now().UnixNano()
	code := rand.New(rand.NewSource(seed)).Int31n(int32(rateMax))
	return code
}

func saveLuckData(code int32, data *gift, sendData string) {
	logger.Printf("lucky,code=%d, gift=%d, name=%s, link=%s, data=%s, left=%d", code, data.id, data.name, data.link, sendData, data.left)
}
