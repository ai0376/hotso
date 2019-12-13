# hotso

[![Build Status](https://travis-ci.org/mjrao/hotso.svg?branch=master)](https://travis-ci.org/mjrao/hotso)

## 介绍
定时抓取热搜数据,并对数据定时进行分词后, 提供数据接口服务


## 工程信息

* cmd/hotso

    hotso 进程抓取网络数据并存储到mongodb 中
    ```
    cd cmd/hotso
    go build
    ````
* cmd/hotword

    hotword.py 脚本将mongodb 中的热搜数据定时分词到redis中

* cmd/service

    service 进程提供对外的http数据访问接口服务
    ```
    cd cmd/service
    go build
    ```

* cmd/mongobackup

    mongobackup.sh 脚本用于定时备份mongodb数据库

* cmd/webdavcli

    webdavcli 进程是定时将备份的mongodb 数据文件上传到云盘(目前程序使用支持webdav的坚果云，坚果云支持历史查询)
    ```
    cd cmd/webdavcli
    go build
    ```


* config/config.json 

    config.json  是所有进程的配置文件


## 进程部署

定时任务的进程和脚本统一用linux 的crond 服务

`crontab -e`

```
0 8-22/1 * * * /opt/hotso/app/hotso

0 0 * * * /usr/local/python3/bin/python3 /opt/hotso/hotword/hotword.py

0 1 * * 1 /opt/hotso/mongobackup/mongobackup.sh

0 3 * * 1 /opt/hotso/webdavcli/webdavcli
```

`service crond reload`

启动service 服务

`
./start_service.sh
`

## 数据样例

[**hotso.top**](hotso.top)

* [微博热搜 10条](http://hotso.top/hotso/v1/hotso/weibo/json/10)

* [百度搜索热点 10条](http://hotso.top/hotso/v1/hotso/baidu/json/10)

* [知乎热榜 10条](http://hotso.top/hotso/v1/hotso/zhihu/json/10) 
    
* [水木10大 10条](http://hotso.top/hotso/v1/hotso/shuimu/json/10)

* [天涯热帖 10条](http://hotso.top/hotso/v1/hotso/tianya/json/10)

* [V2EX最热 10条](http://hotso.top/hotso/v1/hotso/v2ex/json/10)

* [微博热搜分词 10条](http://hotso.top/hotso/v1/hotword/weibo/json/10)

## 声明

抓取数据仅供个人学习，以及本人在工作时间关注实时热点提供方便，无任何商业用途，如有侵权，联系删除

## 反馈

roading@pm.me