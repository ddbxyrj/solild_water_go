package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/kardianos/service"
)

const savetxt = "./soildDate.txt"
const last_water_point_setting = "./last_water_point_setting.txt"
const logfile = "./log.txt"

var iflog = false

type soildsetting struct {
	setting  string
	mutex    sync.Mutex
	ifupdate string
}

func (s *soildsetting) set(date string) {
	s.mutex.Lock()
	s.setting = date
	s.mutex.Unlock()
}
func (s *soildsetting) setUpdateOpen() {
	s.mutex.Lock()
	s.ifupdate = "update"
	s.mutex.Unlock()
}
func (s *soildsetting) setUpdateClose() {
	s.mutex.Lock()
	s.ifupdate = "noupdate"
	s.mutex.Unlock()
}
func (s *soildsetting) getUpdatestatus() string {
	p := ""
	s.mutex.Lock()
	p = s.ifupdate
	s.mutex.Unlock()
	return p
}
func (s *soildsetting) get() string {
	pp := ""
	s.mutex.Lock()
	pp = s.setting
	s.mutex.Unlock()
	return pp
}

var soildsett soildsetting

func soild_date(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	// 土壤温度
	// 土壤水分
	// 土壤电导率
	// 土壤酸碱度
	// 土壤氮
	// 土壤磷
	// 土壤钾

	//土壤温度 := req.Form["0"]
	savingstirng := ""
	for i := 0; i < 7; i++ {
		savingstirng = savingstirng + req.Form[strconv.Itoa(i)][0] + " "
	}
	timeStr := time.Now().Format("2006-01-02 15:04:05") //当前时间的字符串，2006-01-02 15:04:05据说是golang的诞生时间，固定写法
	savingstirng = savingstirng + timeStr
	file, err := os.OpenFile(savetxt, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	write := bufio.NewWriter(file)
	write.WriteString(savingstirng + "\r\n")
	write.Flush()
	file.Close()
	// resstring := "dduduududu1~~~~\n"
	// fmt.Fprintf(w, resstring)
}

func msg(w http.ResponseWriter, req *http.Request) {
	savingstirng := ""
	timeStr := time.Now().Format("2006-01-02 15:04:05") //当前时间的字符串，2006-01-02 15:04:05据说是golang的诞生时间，固定写法
	file, err := os.OpenFile(savetxt, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}

	req.ParseForm()
	msg := req.Form["msg"][0]
	switch msg {
	case "0":
		savingstirng = "stop_water"
		break
	case "1":
		savingstirng = "open_water"
		break
	case "3": //获取参数
		// int water_point = 2500;           // 要浇水的湿度 0-10000对应 0-100%
		// int water_time = 10;              //每次浇水多少毫秒
		// unsigned long waiter_waiter = 60; //在完成浇水后 等待水渗透的时间毫秒
		//第四个 检测间隔
		//等待检测器r485通讯返回间隔
		savingstirng = ""
		resstring := soildsett.get() //最后一位要加@
		fmt.Fprintf(w, resstring)
		break
	case "4":
		savingstirng = "water_point_achived_stop_water"
		break
	case "5":
		savingstirng = "soild_secern_error_reboot_delay_15s"
		break
	case "6":
		savingstirng = "soild_secern_error_total_die"
		break
	case "7":
		//savingstirng = "checkUpdate"
		resstring := soildsett.getUpdatestatus()
		fmt.Fprintf(w, resstring)
		break
	case "8":
		if iflog {

			logtxt := req.Form["logmsg"][0]

			decodedBytes, err := base64.StdEncoding.DecodeString(logtxt)
			if err != nil {
				fmt.Println("Error decoding Base64:", err)
				return
			}

			logtxt = string(decodedBytes)

			if logtxt != "" {
				// 打开文件，如果文件不存在则创建它
				file, err := os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
				if err != nil {
					fmt.Println("无法创建文件:", err)
					return
				}
				defer file.Close()
				logtxt = logtxt + " " + timeStr + "\n"
				// 写入数据到文件

				//写入文件时，使用带缓存的 *Writer
				write := bufio.NewWriter(file)
				write.WriteString(logtxt)
				write.Flush()
				if err != nil {
					fmt.Println("无法写入文件:", err)
					return
				}
			}
		}
		break
	case "9":
		iflog = !iflog
		break
	}
	if savingstirng != "" {
		savingstirng = savingstirng + " " + timeStr
		write := bufio.NewWriter(file)
		write.WriteString(savingstirng + "\r\n")
		write.Flush()
		file.Close()
	}

}
func msgfromcontral(w http.ResponseWriter, req *http.Request) {
	savingstirng := ""
	timeStr := time.Now().Format("2006-01-02 15:04:05") //当前时间的字符串，2006-01-02 15:04:05据说是golang的诞生时间，固定写法
	file, err := os.OpenFile(savetxt, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}

	req.ParseForm()
	msg := req.Form["msg"][0]
	switch msg {
	case "0":
		soildsett.set(req.Form["setmsg"][0])
		savingstirng = "get_new_setting:" + req.Form["setmsg"][0]
		writeToFile(last_water_point_setting, req.Form["setmsg"][0])
		break
	case "1":
		break
	case "2":
		soildsett.setUpdateOpen()
		break
	case "3":
		soildsett.setUpdateClose()
		break

	}
	if savingstirng != "" {
		savingstirng = savingstirng + " " + timeStr
		write := bufio.NewWriter(file)
		write.WriteString(savingstirng + "\r\n")
		write.Flush()
		file.Close()
	}

}
func maint() {
	// 检查文件是否存在

	createFile(logfile)
	createFile(savetxt)
	createFile(last_water_point_setting)
	ss, err := readFileAsString(last_water_point_setting)
	if ss == "" || err != nil {
		soildsett.set("6000@5000@600000@60000@2000@")
	}
	http.HandleFunc("/soild", soild_date)
	http.HandleFunc("/info", msg)
	http.HandleFunc("/contral", msgfromcontral)
	http.ListenAndServe(":6449", nil)
}

type programlinux struct {
	conf *conf
}

func (s *programlinux) Start(c service.Service) error {
	_ = c
	maint()
	return nil
}

func (s *programlinux) Stop(c service.Service) error {

	return nil
}

func createFile(filename string) error {
	// 检查文件是否存在
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		// 文件不存在，创建文件
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		// 关闭文件
		defer file.Close()
	}
	return nil
}

func readFileAsString(filename string) (string, error) {
	// 读取文件内容为字节数组
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	// 将字节数组转换为字符串
	str := string(data)
	return str, nil
}

func writeToFile(filename string, content string) error {
	// 将文本写入文件
	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		return err
	}
	return nil
}
