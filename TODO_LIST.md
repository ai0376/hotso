
* mongodb 查询热度为"沸" 的最近三条微博数据

```
db.weibo.find({"data.state":"沸"},{"data":{$elemMatch:{"state":"沸"}}, "intime":1}).sort({"intime":-1}).limit(3)
```

输出：
```
/* 1 */
{
    "_id" : ObjectId("5dc817e1ebeae9d5302aa0dd"),
    "intime" : NumberLong(1573394401),
    "data" : [ 
        {
            "title" : "双十一晚会",
            "reading" : "2029322",
            "url" : "https://s.weibo.com/weibo?q=%E5%8F%8C%E5%8D%81%E4%B8%80%E6%99%9A%E4%BC%9A&Refer=top",
            "state" : "沸",
            "top" : "5"
        }
    ]
}

/* 2 */
{
    "_id" : ObjectId("5dc7fbc1ebeae9d5302a9df0"),
    "intime" : NumberLong(1573387201),
    "data" : [ 
        {
            "top" : "1",
            "title" : "冷冷冷冷冷冷冷冷冷",
            "reading" : "2794022",
            "url" : "https://s.weibo.com/weibo?q=%23%E5%86%B7%E5%86%B7%E5%86%B7%E5%86%B7%E5%86%B7%E5%86%B7%E5%86%B7%E5%86%B7%E5%86%B7%23&Refer=top",
            "state" : "沸"
        }
    ]
}

/* 3 */
{
    "_id" : ObjectId("5dc7dfa2ebeae9d5302a9b03"),
    "intime" : NumberLong(1573380002),
    "data" : [ 
        {
            "url" : "https://s.weibo.com/weibo?q=%23%E5%8F%8C%E5%8D%81%E4%B8%80%E6%99%9A%E4%BC%9A%E8%8A%82%E7%9B%AE%E5%8D%95%23&Refer=top",
            "state" : "沸",
            "top" : "4",
            "title" : "双十一晚会节目单",
            "reading" : "1265685"
        }
    ]
}
```

* 标题模糊查找

```
db.weibo.find({"data.title":{"$regex":"别墅"}},{"data":{$elemMatch:{"title":{"$regex":"别墅"}}}, "intime":1}).sort({"intime":-1})
```