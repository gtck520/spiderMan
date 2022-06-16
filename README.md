# spiderMan
go学习实战cli程序之--爬虫小工具。使用cobra构建cli应用，使用colly库爬取数据
# 使用方法
下载代码之后 go build -o spider.exe main.go（windows）<br>
然后执行 spider add  url -name NAME 生成需要爬取网站的配置文件<br>
在目录siteconfig 找到NAME文件命名的json配置文件<br>
{<br>
    "name": "sina",   //名称<br>
    "tbname": "sina",   //表名<br>
    "url": "https://feed.sina.com.cn/api/roll/get?pageid=121&lid=1356&num=20&versionNumber=1.2.4&page=2&encode=utf-8&callback=feedCardJsonpCallback&_=1632886493594", //爬取链接<br>
    "rules": [ //爬取规则<br>
        {<br>
            "type": 2, // 1 html 2正则<br>
            "name": "列表",<br>
            "field": "list",<br>
            "match": "\"url\":\"(.*?)\"",  // 提取链接的规则，如果type为 1: 此处为jquery选择器语法  为2: 正则表达式<br>
            "pagerule": { //分页提取规则<br>
                "type": 2, // 1 html 2正则<br>
                "page": 0, //指定爬取页面 为0 爬取所有<br>
                "num": 5, //指定爬取页数  为0 爬取所有 <br>
                "match": "https://feed.sina.com.cn/api/roll/get?pageid=121&lid=1356&num=20&versionNumber=1.2.4&page={page}&encode=utf-8&callback=feedCardJsonpCallback&_=1632886493594" //分页链接提取规则，type为1：jQuery选择器，为2：为固定的分页提取链接 将分页用变量{page}替换 如：示例<br>
            },<br>
            "submatch": "body", //内容提取的位置 jquery选择器语法<br>
            "subrules": [//具体内容页的提取规则<br>
                {<br>
                    "type": 1,  // 1 html 2正则<br>
                    "name": "标题", //名称<br>
                    "field": "title",//字段名  对应数据库或者csv 的字段<br>
                    "match": ".main-title" //jquery选择器语法<br>
                },<br>
                {<br>
                    "type": 1,<br>
                    "name": "内容",<br>
                    "field": "content",<br>
                    "match": ".article"<br>
                }<br>
            ]<br>
        }<br>
    ],<br>
    "out": 1 //输出方式 1 csv，2 mysql<br>
}<br>
以上配置 仅用于普通的简单的爬虫，如果有更复杂的情况，则需要扩展单独针对着写
