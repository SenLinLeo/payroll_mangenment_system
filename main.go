package main

import (
	"geek-nebula/models"
	"geek-nebula/utils"
	"github.com/astaxie/beego"
	"github.com/patrickmn/go-cache"
	_ "geek-nebula/routers"
	"time"
)

func main() {
	models.Init()
	utils.Che = cache.New(60*time.Minute, 120*time.Minute)
	beego.Run()
}
