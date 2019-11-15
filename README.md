# hotso

[![Build Status](https://travis-ci.org/mjrao/hotso.svg?branch=master)](https://travis-ci.org/mjrao/hotso)

## 介绍
定时抓取百度和微博热搜数据,并对数据定时进行分词后, 提供数据接口服务


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

## 数据接口样例

* 微博热搜10条数据
    http://host:port/hotso/v1/hotso/weibo/json/10

* 百度热搜10条数据
    http://host:port/hotso/v1/hotso/baidu/json/10

* 知乎热搜10条数据
    http://host:port/hotso/v1/hotso/zhihu/json/10


* 微博热搜分词后热词10条
    http://host:port/hotso/v1/hotword/weibo/json/10

* 百度热搜分词后热词10条
    http://host:port/hotso/v1/hotword/baidu/json/10

如果需要20条数据，将后面10改为20即可

热搜数据最多提供50条

热词最多提供100条

## 样例

* [weibo 10条](http://121.41.23.201:8806/hotso/v1/hotso/weibo/json/10)

* [baidu 10条](http://121.41.23.201:8806/hotso/v1/hotso/baidu/json/10)

* [zhihu 10条](http://121.41.23.201:8806/hotso/v1/hotso/zhihu/json/10) 
    
    由于知乎热搜需要登陆才能获取，目前做法是人工通过浏览器登陆之后保存cookie,并将cookie存储在云盘里，程序通过webdav方式获取保存的cookie，无任何商业用途，仅供学习和上班时间看知乎热搜方便，希望知乎大佬不要封了我的IP和账户

* [hotwords 10条](http://121.41.23.201:8806/hotso/v1/hotword/weibo/json/10)


Chrome浏览器 + JSON Formatter 扩展体验程序体验效果更佳

![hot_weibo.png](https://i.loli.net/2019/09/24/XxbJaI8n59u4mM2.png "微博热搜")

![hot_baidu.png](https://i.loli.net/2019/09/24/4o89aSig1WfmGhl.png "百度实时热搜")

![hot_zhihu.png](https://i.loli.net/2019/09/24/TwLYEqAm7duDB41.png "知乎热搜")

![hotword_weibo.png](https://i.loli.net/2019/09/24/tyEFzrcdkmHYTlp.png "微博热词")

## 声明

抓取数据仅供个人学习，以及本人在工作时间关注实时热点提供方便，无任何商业用途，如有侵权，联系删除

## 反馈

roading@pm.me