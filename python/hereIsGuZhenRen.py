# coding=utf-8
import urllib2
import sys
import threading
from bs4 import BeautifulSoup


#设置编码
reload(sys)
sys.setdefaultencoding('utf-8')


#获取一个章节的内容
def getChapterContent(file,name,content):
    try:
        file.write("\n" + str(name.text) + "\n");
        data = str(content).replace("<br/><br/>", "\n").replace('<script></script>', "").replace('</div>', "").replace('<div id="content">', "").strip()
        if data=="正在手打中，请稍等片刻，内容更新后，需要重新刷新页面，才能获取最新更新！":
            # print(name)
            threading.Timer(600, getCurrentUrlBooks(starturl + url)).start()
        else:
            file.write(data)
    except BaseException:
        #如果出错了，就重新运行一遍
        print(BaseException.message)
        getChapterContent(file, name, content)
    else:
        print(name)

#获取全部内容
def getCurrentUrlBooks(url):
    response = urllib2.urlopen(url)
    the_page = response.read()
    soup = BeautifulSoup(the_page, "html.parser")
    bookName = soup.select("div[class='bookname'] > h1")[0]
    # 先创建.txt文件，然后获取文本内容写入
    bookContent = soup.select("div[id='content']")[0]
    li_plants=bookContent.script
    li_plants.clear()
    bookFile = open("所以这里是蛊真人.txt", "a+")
    getChapterContent(bookFile, bookName, bookContent)
    bookFile.close()
    nextPage = soup.select("div[class='bottem1'] > a")[3]
    # print nextPage['href']
    return nextPage["href"]

#根据传入参数设置从哪里开始下载
starturl = "https://www.dingdiann.com"
url = "/ddk182237/9683466.html"

# getCurrentUrlBooks(starturl + url)

#死循环 直到没有下一章
while True:
    if url.endswith(".html"):
        url = getCurrentUrlBooks(starturl + url)
    else:
        break;
