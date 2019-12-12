#!/usr/bin/python3
# encoding=utf-8

#结巴分词处理数据
# 由于gojieba 有问题,所以直接采用python处理
import jieba
import pymongo
import redis
import os
import sys
import platform
import json
from enum import Enum
import time

class HotType(Enum):
    Weibo=0
    Baidu=1

hot_collection = {
    HotType.Weibo:"weibo",
    HotType.Baidu:"baidu"
}

#加载配置文件
def load_config_file():
    f = os.path.abspath(__file__)
    if platform.system() == "Windows":
        f = f[:f.rfind('\\')]
        f = f + "\\" + "..\\config\\config.json"
    else:  # Linux or Mac
        f = f[:f.rfind('/')]
        f = f + "/" + "../config/config.json"
    text = ""
    with open(f, "r", encoding='UTF-8') as fp:
        text = fp.read()
    return json.loads(text)

# 用户字典文件
def user_dict_file():
    f = os.path.abspath(__file__)
    if platform.system() == "Windows":
        f = f[:f.rfind('\\')]
        f = f + "\\" + "..\\config\\user_dict.txt"
    else:  # Linux or Mac
        f = f[:f.rfind('/')]
        f = f + "/" + "../config/user_dict.txt"
    return f

# 创建停用词list  
def stopwordslist():
    f = os.path.abspath(__file__)
    if platform.system() == "Windows":
        f = f[:f.rfind('\\')]
        f = f + "\\" + "..\\config\\stopword.txt"
    else:  # Linux or Mac
        f = f[:f.rfind('/')]
        f = f + "/" + "../config/stopword.txt"
    stopwords = [line.strip() for line in open(f, 'r', encoding='utf-8').readlines()]  
    return stopwords

#取redis key 参数为集合名字
def get_redis_key(collection):
    tm = time.localtime(time.time())
    return "hotword:"+collection+":"+str(tm.tm_year)

#任务处理时间线
def get_redis_key_intime():
    return "tasktimeline"

#切割词
def cut_word(word):
    seg_list = jieba.cut(word, cut_all=False)
    outstr=''
    for w in seg_list:
        if w not in stopwords:
            s = w.strip()
            if s != '\t' and s !='\n' and s != '':
                outstr += s
                outstr +='/'
    if len(outstr) > 0:
        outstr = outstr[:len(outstr)-1]
    return outstr.split("/")

#将处理的数据时间写入库
def set_time_in_redis(ntime):
    return redisCli.set(get_redis_key_intime(), ntime)

#读取数据处理时间点
def get_time_in_redis():
    return redisCli.get(get_redis_key_intime())

#任务处理
def do_task():
    t = get_time_in_redis()
    begin_unix = 0
    if t == None:
        begin_unix = cfg["hotword"]["begin_unix"]
    else:
        begin_unix = int(t)
    duration_sec = cfg["hotword"]["duration_sec"]
    end_unix = 0
    if begin_unix <= 0:
        begin_unix = 0
    if duration_sec == 0:
        end_unix = int(time.time())
    else:
        end_unix = begin_unix + duration_sec

    binsert = False
    mdb = mongoCli["hotso"]
    for v in hot_collection.values():
        cli = mdb[v]
        # data = cli.find().sort([("intime",-1)]).limit(1)
        data = cli.find({"intime":{"$gte":begin_unix, "$lte":end_unix}})
        for d in data:
            for one_data in d["data"]:
                #print(v)
                if one_data["state"] == "荐":
                    #广告
                    pass
                elif one_data["reading"] == "": 
                    #我党相关
                    pass
                else:
                    words = cut_word(one_data["title"])
                    reading = int(one_data["reading"])
                    for w in words:
                        redisCli.zincrby(get_redis_key(v), reading, w)
                        binsert = True
                pass
    if binsert == True:
        set_time_in_redis(end_unix)

cfg = load_config_file()
mongoCli = pymongo.MongoClient(cfg["mongodb"]["host"])
redisCli = redis.Redis(host=cfg["redis"]["host"], port=cfg["redis"]["port"],db=0)
stopwords = stopwordslist()


#定时任务代码，每天执行一次
if __name__ == "__main__":
    jieba.load_userdict(user_dict_file())# 加载字典
    #print(user_dict_file())
    do_task()