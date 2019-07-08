# coding=utf-8
import re
import json
import pymysql
import requests
import threading
from LogUtils import Log

url = 'http://util.online/spider/api/mail'
weburl = 'https://util.online/spider/novel/'
body = {
    "to": "",
    "subject":"XX小说出新章节咯",
    "text": "新的章节是 http://www.baidu.com",
    "html":"<h1>Welcome</h1><p>That was easy!</p ><a href=' '>新的章节</a >"
}
headers = {'content-type': "application/json"}
db = pymysql.connect("39.104.226.149", "root", "root", "spider", charset='utf8')
keyVa = {}
log = Log()


def main():
    # 查询所有小说
    da = seloss()
    # 查询key_value 表中的最新章
    selKeyValue(da)
    # 查询小说是否更新
    selNovel()
    threading.Timer(60, main).start()


# 查询所有小说
def seloss():
    ossCur = db.cursor()
    ossCur.execute("SELECT id,crawl_url FROM oss WHERE type='NOVEL' AND complete=0")
    db.commit()
    data = ossCur.fetchall()
    return data


# 查询key_value表中最新章
def selKeyValue(da):
    for oss_id in da:
        # 如果没爬过（新添加的小说，或者无法爬取的）忽略
        if oss_id[1]:
            # selNewnovel(oss_id[0])
            upkeyValue(oss_id[0])


# 更新内存中的key_value
def upkeyValue(ids):
    keyCur = db.cursor()
    keyCur.execute("SELECT value FROM key_value WHERE `key`='%s'" %("novel." + str(ids)))
    db.commit()
    val = keyCur.fetchone()
    if val:
        keyVa.update({ids: val[0]})
    selNewnovel(ids)


# 获取小说最新章节
def selNewnovel(id):
    novelCur = db.cursor()
    novelCur.execute("SELECT chapter_index FROM novel WHERE oss_id=%d ORDER BY chapter_index DESC LIMIT 1" %(id))
    db.commit()
    da = novelCur.fetchone()
    if da:
        if keyVa.has_key(id):
            if da[0] > int(keyVa[id]):
                # 更新key_value
                updKeyValue(id, da[0])
                # keyVa.update({id: da[0]})
        else:
            # 插入key_value表 加入内存
            insertKeyValue(id, da[0])
            keyVa.update({id: da[0]})


# 查询小说是否更新
def selNovel():
    novCur = db.cursor()
    for id in keyVa:
        index = int(keyVa[id])
        novCur.execute("SELECT chapter_index,chapter_title,content FROM novel WHERE oss_id=%d AND chapter_index>%d ORDER BY chapter_index" %(id, index))
        db.commit()
        cos = novCur.fetchall()
        if cos:
            flag = 0
            end = len(cos)
            for co in cos:
                flag += 1
                if flag >= end:
                    updKeyValue(id, co[0])
                    keyVa.update({id: co[0]})
                log.info(co[1])
                # 查书名
                name = bookname(id)
                # 查订阅表
                mails = usermail(id)
                for m in mails:
                    sendmail(m, name, id, co[2], co[1])


# 查询书名
def bookname(id):
    nameCur = db.cursor()
    nameCur.execute("SELECT name FROM oss WHERE id='%d'" %(id))
    db.commit()
    na = nameCur.fetchone()[0]
    log.info(na)
    return na


# 查订阅用户邮箱
def usermail(id):
    mails = []
    rsscur = db.cursor()
    rsscur.execute("SELECT user_id FROM rss WHERE oss_id='%d'" %(id))
    db.commit()
    userIds = rsscur.fetchall()
    for userid in userIds:
        rsscur.execute("SELECT mail,nick FROM user WHERE id='%d'" %(userid))
        db.commit()
        data = rsscur.fetchone()
        if data:
            mail = data[0]
            # 判断邮箱和合理性
            pattern = r'^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$'
            if mail and re.match(pattern,mail) is not None:
                mails.append(mail)
            else:
                log.info('邮箱不正确')
    return mails


# 发邮件
def sendmail(mail,name,id,content,title):
    try:
        # 拿到最新章的标题和内容
        body["subject"] = "《" + name.encode('UTF-8') + "》 小说出最新章啦"
        body["to"] = mail
        # body["text"] = "新的章节是 " + co[1].encode('UTF-8')
        body["html"] = "<h1>" + title.encode('UTF-8') + "</h1><div style='white-space: pre-wrap;font-size:15px;'>" + content.encode('UTF-8') + "</div><a href='"+ weburl + str(id) +"'>最新章节</a >"
        response = requests.post(url, data = json.dumps(body), headers = headers)
    except requests.exceptions, e:
        log.error("mail send err:"+e)


# 更新key_value表
def updKeyValue(key,value):
    try:
        log.info("update")
        upCur = db.cursor()
        upCur.execute("UPDATE key_value SET value='%s' WHERE `key`='%s'" %(str(value), "novel."+str(key)))
        db.commit()
    except Exception, e:
        log.error("sql err:" + e)
        db.rollback()


# 插入key_value表
def insertKeyValue(key,value):
    try:
        log.info("insert")
        incur = db.cursor()
        incur.execute("INSERT INTO key_value (`key`,value)VALUES ('%s','%s')" %("novel." + str(key), str(value)))
        db.commit()
    except Exception, e:
        log.error("sql err:" + e)
        db.rollback()


if __name__ == '__main__':
    main()
