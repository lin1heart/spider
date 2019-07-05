# coding=utf-8
import sys
import ssl
import string
import urllib2
import pymysql
import logging
import operator
import threading
from bs4 import BeautifulSoup


#设置编码
reload(sys)
sys.setdefaultencoding('utf-8')
try:
    _create_unverified_https_context = ssl._create_unverified_context
except AttributeError:
    pass
else:
    ssl._create_default_https_context = _create_unverified_https_context
#log
logger = logging.getLogger("dingdian_spider")
logger.setLevel(logging.DEBUG)
fh = logging.FileHandler("spider_error.log")
fh.setLevel(logging.DEBUG)
ch = logging.StreamHandler()
ch.setLevel(logging.DEBUG)
formatter = logging.Formatter("%(asctime)s - %(name)s - %(levelname)s : %(message)s")
fh.setFormatter(formatter)
ch.setFormatter(formatter)
logger.addHandler(fh)
logger.addHandler(ch)

#根据传入参数设置从哪里开始下载
starturl = "https://www.dingdiann.com"
searchurl = "https://www.dingdiann.com/searchbook.php?keyword="
db = pymysql.connect("39.104.226.149", "root", "root", "spider", charset='utf8')


#获取章节内容
def spiderContent(url,id, name):
    try:
        response = urllib2.urlopen(url, timeout=60)
        the_page = response.read()
        soup = BeautifulSoup(the_page, "html.parser")
        bookName = soup.select("div[class='bookname'] > h1")[0].text
        bookContent = soup.select("div[id='content']")[0]
        nextPage = soup.select("div[class='bottem1'] > a")[3]["href"]
        li_plants = bookContent.script
        if li_plants:
            li_plants.clear()
        data = str(bookContent).replace("\\","").replace("<br/><br/>", "\n").replace('<script></script>', "").replace('</div>', "").replace('<div id="content">', "").replace('\'', '\\\'').strip()
        checkdata = "正在手打中，请稍等片刻，内容更新后，需要重新刷新页面，才能获取最新更新！"
        if data == checkdata and nextPage.endswith(".html"):
            logger.info("本章节无内容 | " + name)
            insert("",id,bookName,url,starturl+nextPage)
        elif data != checkdata and "" != data and not nextPage.endswith(".html"):
            logger.info("最新章 "+bookName + " | " + name)
            insert(data,id,bookName,url,"")
        elif data != checkdata and nextPage.endswith(".html"):
            logger.debug("正常章节 "+bookName + " | " + name)
            insert(data,id,bookName,url,starturl+nextPage)
        else:
            logger.info("本章节未更新或者获取章节异常 "+bookName +" | "+ name)
            return
    except Exception, e:
        logger.error(e)
        logger.error(data)




#从目录爬取小说,获取第一章
def spiderM(url, id, name):
    response = urllib2.urlopen(url)
    the_page = response.read()
    soup = BeautifulSoup(the_page, "html.parser")
    books = soup.select("dd > a")
    length = len(books)
    count = 0
    hreflist = []
    textlist = []
    if length > 20:
        for a in books:
            if count > 20:
                break
            count += 1
            hreflist.append(int(a["href"].split('/')[2].split('.')[0]))
            textlist.append(a.text)
    else:
        for a in books:
            hreflist.append(int(a["href"].split('/')[2].split('.')[0]))
            textlist.append(a.text)

    min_index, min_number = min(enumerate(hreflist), key=operator.itemgetter(1))
    updateF(url+str(min_number)+".html", id)
    spiderContent(url+str(min_number)+".html", id, name)


#通过网站搜索功能 先搜索小说，再爬取
def searchNovel(name, id):
    response = urllib2.urlopen(searchurl+urllib2.quote(str(name)))
    the_page = response.read()
    soup = BeautifulSoup(the_page, "html.parser")
    book = soup.select("span[class='s2'] > a")
    if not book:
        return
    for boo in book:
        if boo.text == name:
            spiderM(starturl + boo["href"], id, name)



# 判断是否需要查询小说，有crawl_url就不需要查询
def spiderStart(data):
    for da in data:
        ossId = da[0]
        name = da[1]
        curl = da[2]
        if not name and not ossId:
            continue
        if curl and string.find(curl, starturl) != -1:
            nurl = search(ossId)
            if nurl:
                spiderContent(nurl, ossId, name)
            else:
                spiderContent(curl, ossId, name)
        elif curl:
            continue
        else:
            searchNovel(name, ossId)

def search(id):
    scursor = db.cursor()
    scursor.execute("SELECT crawl_url,next_crawl_url FROM novel WHERE oss_id=%d ORDER BY chapter_index DESC LIMIT 1" %(id))
    w = scursor.fetchone()
    if not w:
        return ''
    elif w[1]:
        return w[1]
    else:
        return w[0]


def insert(data,id,name,url,next):
    icursor = db.cursor()
    sqlInsert = "INSERT INTO novel(chapter_index,chapter_title,oss_id,content,crawl_url,next_crawl_url)VALUES ('%d','%s','%d','%s','%s','%s')"
    icursor.execute("SELECT chapter_index,chapter_title,next_crawl_url FROM novel WHERE oss_id=%d ORDER BY chapter_index DESC LIMIT 1" %(id))
    sdata = icursor.fetchone()
    if not sdata:
        index = 1
        icursor.execute(sqlInsert %(index,name,id,data,url,next))
        db.commit()
    if sdata and sdata[1] != name:
        index = sdata[0] + 1
        icursor.execute(sqlInsert %(index,name,id,data,url,next))
        db.commit()
    elif sdata and not sdata[2]:
        c_index = sdata[0]
        icursor.execute("UPDATE novel SET next_crawl_url='%s' WHERE oss_id=%d AND chapter_index=%d" %(next, id, c_index))
        db.commit()

def updateF(curl,id):
    ucursor = db.cursor()
    ucursor.execute("UPDATE oss SET crawl_url='%s' WHERE id=%d" %(curl,id))
    db.commit()

def select():
    ocursor = db.cursor()
    ocursor.execute("SELECT id,name,crawl_url,url FROM oss WHERE type='NOVEL' AND complete=FALSE")
    data = ocursor.fetchall()
    return data

def main():
    list = select()
    spiderStart(list)
    # db.close()
    threading.Timer(1.5, main).start()

if __name__ == '__main__':
    main()
