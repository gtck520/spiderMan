# spiderMan
go学习实战cli程序之--爬虫小工具，使用colly库
# 使用方法
下载代码之后 go build -o spider.exe main.go（windows）
然后执行 spider add  url -name NAME 生成需要爬取网站的配置文件
在目录siteconfig 找到NAME文件命名的json配置文件
{
    "name": "sina",   //名称
    "url": "https://feed.sina.com.cn/api/roll/get?pageid=121&lid=1356&num=20&versionNumber=1.2.4&page=2&encode=utf-8&callback=feedCardJsonpCallback&_=1632886493594", //爬取链接
    "rules": [ //爬取规则
        {
            "type": 2, // 1 html 2正则
            "name": "列表",
            "field": "list",
            "match": "\"url\":\"(.*?)\"",  // 提取链接的规则，如果type为 1: 此处为jquery选择器语法  为2: 正则表达式
            "pagerule": { //分页提取规则
                "type": 2, // 1 html 2正则
                "page": 0, //指定爬取页面 为0 爬取所有
                "num": 5, //指定爬取页数  为0 爬取所有 
                "match": "https://feed.sina.com.cn/api/roll/get?pageid=121&lid=1356&num=20&versionNumber=1.2.4&page={page}&encode=utf-8&callback=feedCardJsonpCallback&_=1632886493594" //分页链接提取规则，type为1：jQuery选择器，为2：为固定的分页提取链接 将分页用变量{page}替换 如：示例
            },
            "submatch": "body", //内容提取的位置 jquery选择器语法
            "subrules": [//具体内容页的提取规则
                {
                    "type": 1,  // 1 html 2正则
                    "name": "标题", //名称
                    "field": "title",//字段名  对应数据库或者csv 的字段
                    "match": ".main-title" //jquery选择器语法
                },
                {
                    "type": 1,
                    "name": "内容",
                    "field": "content",
                    "match": ".article"
                }
            ]
        }
    ],
    "out": 1
}
以上配置 仅用于普通的简单的爬虫，如果有更复杂的情况，则需要扩展单独针对着写
