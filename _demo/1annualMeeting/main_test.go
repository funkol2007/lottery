package main

import (
	"fmt"
	"sync"
	"testing"

	"github.com/kataras/iris/v12/httptest"
)

func TestMVC(t *testing.T) {
	e := httptest.New(t, newApp())
	var wg sync.WaitGroup

	e.GET("/").Expect().Status(httptest.StatusOK).Body().Equal("当前总共参与抽奖用户数：0\n")
	for i := 0; i != 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			e.POST("/import").WithFormField("users", fmt.Sprintf("用户_%d", i)).Expect().Status(httptest.StatusOK)
		}(i)
	}
	wg.Wait()

	e.GET("/").Expect().Status(httptest.StatusOK).Body().
		Equal("当前总共参与抽奖用户数：100\n")
	e.GET("/lucky").Expect().Status(httptest.StatusOK)
	e.GET("/").Expect().Status(httptest.StatusOK).Body().
		Equal("当前总共参与抽奖用户数：99\n")
}
