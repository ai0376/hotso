# hotso

## 介绍
定时抓取百度和微博热搜数据,并对数据定时进行分词后, 提供数据接口服务


## 工程信息

* cmd/hotso

    hotso 负责进程抓取网络数据并存储到mongodb 中
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

## 数据接口样例

* 微博热搜10条数据
    http://host:port/hotso/v1/hotso/weibo/json/10

* 百度热搜10条数据
    http://host:port/hotso/v1/hotso/baidu/json/10

* 微博热搜分词后热词10条
    http://host:port/hotso/v1/hotword/weibo/json/10

* 百度热搜分词后热词10条
    http://host:port/hotso/v1/hotword/baidu/json/10

如果需要20条数据，将后面10改为20即可

热搜数据最多提供50条

热词最多提供100条

## 反馈

roading@pm.me