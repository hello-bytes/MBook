package main

import (
	"app"
	"io"
	"log"
	"os"
	"strings"
)

/*
var (
	cmd string
)

func init() {
	flag.StringVar(&g, "cmd", "none", "set the ")
}*/

func printUsage() {
	log.Println("MBook version: mbook/1.0.3")
	log.Println("Usage: mbook init|console|daemon|help")
	log.Println("")
	log.Println("Options:")
	log.Println("  init:generate the config, use the default value")
	log.Println("  console:sync run on with console")
	log.Println("  daemon:run as a daemon service")
	log.Println("  help:print the usage")
}

func main() {

	//log.Println("the cmder count :", len(os.Args))
	if len(os.Args) != 2 {
		printUsage()
		return
	}

	cmd := os.Args[1]
	if strings.Compare(cmd, "help") == 0 {
		printUsage()
		return
	} else if strings.Compare(cmd, "console") == 0 {
		consoleRun()
	} else if strings.Compare(cmd, "daemon") == 0 {
		daemonRun()
	} else if strings.Compare(cmd, "init") == 0 {
		initAppConfig()
	}

	log.Println(cmd)
}

func copyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}

func initAppConfig() {
	configFile := app.Path("resource/config/env.config")
	sampleConfig := app.Path("resource/config/sample.env.config")
	if app.IsFileExist(configFile) {
		log.Println("file : ", configFile, " is exist, no need to init, please run and fix config by /backend/config")
		return
	}

	_, err := copyFile(configFile, sampleConfig)
	if err == nil {
		log.Println("Your config is inited, and you can access with path : ", configFile)
		log.Println("Try to run the web server and access by backend")
		log.Println("just type 'mbook daemon' and access http://127.0.0.1:8080")
	} else {
		log.Println("Your init config  is failed", err)
		log.Println("maybe you shoud check access and check file exist for ", sampleConfig)
	}
}

func consoleRun() {
	app.InitConfig()
	startWeb()
}

func daemonRun() {
	//app.InitConfig()
	//startWeb()
}
