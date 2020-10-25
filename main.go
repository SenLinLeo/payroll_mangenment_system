package main

import (
	"geek-nebula/models"
	_ "geek-nebula/routers"
	"geek-nebula/utils"
	"github.com/astaxie/beego"
	"github.com/patrickmn/go-cache"
	"time"
)

func main() {
	models.Init()
	utils.Che = cache.New(60*time.Minute, 120*time.Minute)
	beego.Run()
}
