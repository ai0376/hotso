# hotso

[![Build Status](https://travis-ci.org/ai0376/hotso.svg?branch=master)](https://travis-ci.org/ai0376/hotso)

Demo site:  [**hotso.top**](http://hotso.top)

## 介绍
定时抓取热搜数据,并对数据定时进行分词后, 提供数据接口服务

## 依赖存储服务

* redis

    存储分词以及年度榜单数据

* mongodb

    存储站点抓取的资讯信息

## 工程信息

* cmd/hotso

    hotso 进程抓取网络数据并存储到mongodb 中
    ```
    cd cmd/hotso
    go build
    ````
* cmd/hottop

    hottop 进程每周从mongodb 中取得数据 排序到redis中，得出年度榜单
    ```
    cd cmd/hottop
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

* web/*

    一个简单的index.html页面，部署在nginx的document root 目录（粗糙的前端展示）


## 进程部署

定时任务的进程和脚本统一用linux 的crond 服务

`crontab -e`

```
0 0 * * * /usr/local/python3/bin/python3 /opt/hotso/hotword/hotword.py
0 8-22/1 * * * /opt/hotso/app/hotso
0 2 * * 1 /opt/hotso/hottop/hottop
0 3 * * 1 /opt/hotso/webdavcli/webdavcli
0 1 * * 1 /opt/hotso/mongobackup/mongobackup.sh
```

`service crond reload`

启动service 服务

`
./start_service.sh
`

## 数据样例

[**hotso.top**](http://hotso.top)

* [微博热搜 10条](http://hotso.top/hotso/v1/hotso/weibo/10)

* [百度搜索热点 10条](http://hotso.top/hotso/v1/hotso/baidu/10)

* [知乎热榜 10条](http://hotso.top/hotso/v1/hotso/zhihu/10) 
    
* [水木10大 10条](http://hotso.top/hotso/v1/hotso/shuimu/10)

* [天涯热帖 10条](http://hotso.top/hotso/v1/hotso/tianya/10)

* [V2EX最热 10条](http://hotso.top/hotso/v1/hotso/v2ex/10)

* [微博热搜分词 10条](http://hotso.top/hotso/v1/hotword/weibo/2019/10)

* [查询微博某天热搜数据](http://hotso.top/hotso/v1/query/weibo/2019-12-01/10)

* [查询知乎某天热搜数据](http://hotso.top/hotso/v1/query/zhihu/2019-12-01/10)

* 查询历史某天热搜数据接口
    `http://hotso.top/hotso/v1/query/:type/:day/:num`

    可变字段取值

    :type   weibo/baidu/zhihu/v2ex/shuimu/tianya

    :day    YYYY-MM-DD

    :num    获取数据条数

## 声明

抓取数据仅供个人学习，以及本人在工作时间关注实时热点提供方便，无任何商业用途，如有侵权，联系删除

所有数据均来自第三方(微博，百度，知乎等)，禁止将数据用于任何商业用途

## 反馈

ai0376 at 126.com

## 关于我

http://ownwiki.cn/about/
