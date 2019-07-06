# coding=utf-8
import re
import json
import pymysql
import requests
import threading

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

def main():
    da = seloss()
    selKeyValue(da)
    selNovel()
    threading.Timer(5, main).start()


# 查询所有小说
def seloss():
    ossCur = db.cursor()
    ossCur.execute("SELECT id,crawl_url FROM oss WHERE type='NOVEL'")
    db.commit()
    data = ossCur.fetchall()
    return data

# 查询key_value表中最新章
def selKeyValue(da):
    for oss_id in da:
        if oss_id[1]:
            keyValue(oss_id[0])

# 更新内存中的key_value
def keyValue(ids):
    keyCur = db.cursor()
    keyCur.execute("SELECT value FROM key_value WHERE `key`='%s'" %("novel." + str(ids)))
    db.commit()
    val = keyCur.fetchone()
    if val:
        keyVa.update({ids: val[0]})
        print keyVa

# 查询小说是否更新
def selNovel():
    novCur = db.cursor()
    for id in keyVa:
        index = int(keyVa[id])
        print index
        novCur.execute("SELECT chapter_index,chapter_title,content FROM novel WHERE oss_id=%d AND chapter_index>%d ORDER BY chapter_index DESC LIMIT 1" %(id, index))
        db.commit()
        co = novCur.fetchone()
        print co
        if co:
            updKeyValue(id, co[0])
            keyVa.update({id: co[0]})
            print co[1]
            # 查书名
            name = bookname(id)
            # 查订阅表
            mails = usermail(id)
            for m in mails:
                sendmail(m,name,id,co[2],co[1])

# 查询书名
def bookname(id):
    nameCur = db.cursor()
    nameCur.execute("SELECT name FROM oss WHERE id='%d'" %(id))
    db.commit()
    na = nameCur.fetchone()[0]
    print na
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
            print mail
            # 判断邮箱和合理性
            pattern = r'^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$'
            if mail and re.match(pattern,mail) is not None:
                mails.append(mail)
            else:
                print '邮箱不正确'
    return mails

# 发邮件
def sendmail(mail,name,id,content,title):
    # 拿到最新章的标题和内容
    body["subject"] = "《" + name.encode('UTF-8') + "》 小说出最新章啦"
    body["to"] = mail
    # body["text"] = "新的章节是 " + co[1].encode('UTF-8')
    body["html"] = "<h1>" + title.encode('UTF-8') + "</h1><div style='white-space: pre-wrap;font-size:15px;'>" + content.encode('UTF-8') + "</div><a href='"+ weburl + str(id) +"'>最新章节</a >"
    response = requests.post(url, data = json.dumps(body), headers = headers)
    print response

# 更新key_value表
def updKeyValue(key,value):
    print "update"
    upCur = db.cursor()
    upCur.execute("UPDATE key_value SET value='%s' WHERE `key`='%s'" %(str(value), "novel."+str(key)))
    db.commit()

if __name__ == '__main__':
    main()
