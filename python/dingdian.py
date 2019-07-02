# coding=utf-8
import sys
import ssl
import string
import urllib2
import pymysql
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

#根据传入参数设置从哪里开始下载
starturl = "https://www.dingdiann.com"
searchurl = "https://www.dingdiann.com/searchbook.php?keyword="
db = pymysql.connect("39.104.226.149", "root", "root", "spider", charset='utf8')


#获取章节内容
def spiderContent(url,id):
    response = urllib2.urlopen(url)
    the_page = response.read()
    soup = BeautifulSoup(the_page, "html.parser")
    bookName = soup.select("div[class='bookname'] > h1")[0].text
    bookContent = soup.select("div[id='content']")[0]
    nextPage = soup.select("div[class='bottem1'] > a")[3]["href"]
    li_plants = bookContent.script
    li_plants.clear()
    data = str(bookContent).replace("<br/><br/>", "\n").replace('<script></script>', "").replace('</div>', "").replace('<div id="content">', "").strip()
    if data == "正在手打中，请稍等片刻，内容更新后，需要重新刷新页面，才能获取最新更新！":
        return
    elif not nextPage.endswith(".html"):
        print "最新章"
        return
    else:
        insert(data,id,bookName,url,starturl+nextPage)
        # timer = threading.Timer(1.5, spiderContent(starturl+nextPage,id))
        # timer.start()



#从目录爬取小说
def spiderM(url):
    response = urllib2.urlopen(url)
    the_page = response.read()
    soup = BeautifulSoup(the_page, "html.parser")
    # print soup
    # book = soup.select("div[id='list']")
    book = soup.select("dd > a")
    books = soup.select("dd > a")[0].text
    print book

#通过网站搜索功能 先搜索小说，再爬取
def searchNovel(name):
    response = urllib2.urlopen(searchurl+urllib2.quote(str(name)))
    the_page = response.read()
    soup = BeautifulSoup(the_page, "html.parser")
    book = soup.select("span[class='s2'] > a")
    if not book:
        return
    if book[0].text == name:
        spiderM(starturl + book[0]["href"])



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
            print (nurl)
            if nurl:
                spiderContent(nurl, ossId)
            else:
                spiderContent(curl, ossId)
        elif curl:
            continue
        else:
            searchNovel(name)

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
    print id,name,url,next
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


def select():
    cursor = db.cursor()
    cursor.execute("SELECT id,name,crawl_url,url FROM oss WHERE type='NOVEL' AND complete=FALSE ")
    data = cursor.fetchall()
    return data

def main():

    list = select()
    spiderStart(list)
    db.close()

if __name__ == '__main__':
    main()
