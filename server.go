package main

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/kardianos/service"
)

type (
	conf struct {
		Title   string
		Service service.Config
		HttpIp  HttpIp
	}
	Use_function struct {
		Dns     bool
		Wol     bool
		Http_ip bool
	}
	HttpIp struct {
		Ufw_default_insert_postion int
	}
)

func main() {

	fmt.Println("盈盈秋水，淡淡春山")

	configByte, err := os.ReadFile("./config.toml")
	var conf conf
	if err != nil {

		fmt.Println("未找到配置文件")
		return

	} else {

		toml.Decode(string(configByte), &conf)

	}

	svcConfig := &conf.Service
	prg := &programlinux{
		conf: &conf,
	}

	//prg := &program{}

	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)

	}
	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	//err = s.Run()
	if err != nil {
		logger.Error(err)
	}

	if len(os.Args) < 2 {
		err = s.Run()
		if err != nil {
			logger.Error(err)
		}
		return
	}

	cmd := os.Args[1]

	if cmd == "install" {
		err = s.Install()
		if err != nil {
			fmt.Println("安装失败")
			logger.Error(err)
			return
		}
		err = s.Start()
		if err != nil {
			fmt.Println("安装后运行失败")
			logger.Error(err)
			return
		}
		fmt.Println("安装成功并运行了")
		return
	}
	if cmd == "uninstall" {
		err = s.Uninstall()
		if err != nil {
			logger.Error(err)
			return
		}
		fmt.Println("卸载成功")
		return
	}

	if cmd == "start" {
		err = s.Start()
		if err != nil {
			logger.Error(err)
			return
		}
		fmt.Println("服务启动成功")
		return
	}

	if cmd == "stop" {
		err = s.Stop()
		if err != nil {
			logger.Error(err)
			return
		}
		fmt.Println("服务关闭成功")
		return
	}

	if cmd == "restart" {
		err = s.Stop()
		if err != nil {
			logger.Error(err)
			return
		}
		fmt.Println("服务重启成功")
		return
	}

	if cmd == "status" {
		status, err := s.Status()
		if err != nil {
			logger.Error(err)
			return
		}
		switch status {
		case 1:
			fmt.Println("服务正在运行")
		case 2:
			fmt.Println("服务已经停止")
		default:
			fmt.Println("未知状态")
		}
		return
	}
	if cmd == "test" {
		// go 更新运行()
		// for {
		//go prg.Start()
		// }
		// go 更新运行()
		// go gogoproxy.GMain()

		for {

		}

	}
	if cmd == "testF" {

		// fmt.Println("测试防火墙")
		// 激活远程访问("1.1.1.1")
		// //激活远程访问("0.0.0.0")

	}

	if cmd == "testL" {

		// fmt.Println("测试linux")
		// http运行_linux()
	}
}

var logger service.Logger
