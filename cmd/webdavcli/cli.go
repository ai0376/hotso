package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mjrao/hotso/config"
	"github.com/mjrao/hotso/internal/cloud"
)

var (
	client *cloud.Client
)

//将目录下的文件按时间排序
func sortByTime(pl []os.FileInfo) []os.FileInfo {
	sort.Slice(pl, func(i, j int) bool {
		flag := false
		if pl[i].ModTime().After(pl[j].ModTime()) {
			flag = true
		} else if pl[i].ModTime().Equal(pl[j].ModTime()) {
			if pl[i].Name() < pl[j].Name() {
				flag = true
			}
		}
		return flag
	})
	return pl
}

//上传指定的多个文件
func upload() {
	var err error
	webdavCfg := config.GetConfig().WebDav

	client, err = cloud.Dial(webdavCfg.Host, webdavCfg.User, webdavCfg.Password)
	if err != nil {
		panic(err.Error())
	}

	testDir := webdavCfg.Path
	remoteDir := strings.Replace(webdavCfg.RemoteDir, "\\", "/", -1)
	if remoteDir[len(remoteDir)-1:] != "/" {
		remoteDir = remoteDir + "/"
	}

	for _, v := range webdavCfg.Files {
		if src, e := ioutil.ReadFile(filepath.Join(testDir, v)); e != nil {
			panic(e.Error())
		} else {
			err = client.Upload(src, remoteDir+v)
			if err != nil {
				panic(err.Error())
			}
		}
	}
}

//上传一个指定目录下的最新文件
func uploadOneFileByNewTime() {
	var err error
	webdavCfg := config.GetConfig().WebDav
	client, err = cloud.Dial(webdavCfg.Host, webdavCfg.User, webdavCfg.Password)
	if err != nil {
		panic(err.Error())
	}
	testDir := webdavCfg.Path
	remoteDir := strings.Replace(webdavCfg.RemoteDir, "\\", "/", -1)
	if remoteDir[len(remoteDir)-1:] != "/" {
		remoteDir = remoteDir + "/"
	}

	readerInfos, err := ioutil.ReadDir(testDir)
	if err != nil {
		panic(err.Error())
	}
	readerInfosByTime := sortByTime(readerInfos)
	var newFile = ""
	for _, info := range readerInfosByTime {
		if info.IsDir() {
			continue
		} else {
			newFile = info.Name()
			break
		}
	}
	if newFile != "" {
		if src, e := ioutil.ReadFile(filepath.Join(testDir, newFile)); e != nil {
			panic(e.Error())
		} else {
			err = client.Upload(src, remoteDir+newFile)
			if err != nil {
				panic(err.Error())
			}
		}
	}
}
func main() {
	uploadOneFileByNewTime()
}
