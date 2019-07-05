# coding=utf-8
from wxpy import *
import json
import requests
import sys
reload(sys)
sys.setdefaultencoding('utf-8')

bot=Bot(cache_path=True)

# friend = bot.friends().search('amanooo')[0]
# groups = bot.groups()
me = bot.friends().search('lin~')[0]
# group = groups.search(unicode("八小时html快速入门", 'utf-8'))[0]
# print(friend)
print(me)
# print(group)
tuling = Tuling(api_key='1d2172280e104c378ea8f92e6c5c8ebf')
@bot.register(chats=me) # 接收从指定好友发来的消息，发送者即recv_msg.sender为指定好友friend
def recv_send_msg(recv_msg):
    print('收到的消息：',recv_msg.text) # recv_msg.text取得文本
    if recv_msg.sender == me:
        print 111
        bot.file_helper.send('逗比留言： '+ recv_msg.text)
        recv_msg.forward(bot.file_helper,prefix='逗比留言： ') #在文件传输助手里留一份，方便自己忙完了回头查看
        me.send('hello world!') #
        me.reply('lalallaa')

me.send("hello world!1")
# group.send("开启聊天机器人模式!")

# @bot.register(chats=group) #接收从指定群发来的消息，发送者即recv_msg.sender为组
# def recv_send(recv_msg):
#     print('收到的消息：',recv_msg.text)
#     tuling.do_reply(recv_msg)
#     auto_reply(recv_msg.text)
#     if recv_msg.member == friend:
#         recv_msg.forward(bot.file_helper,prefix='老板发言: ')


# 调用图灵机器人API，发送消息并获得机器人的回复
def auto_reply(text):
    url = "http://openapi.tuling123.com/openapi/api/v2"
    payload = {
        "reqType":0,
        "perception": {
            "inputText": {
                "text": text
            }
        },
        "userInfo": {
            "apiKey": "1d2172280e104c378ea8f92e6c5c8ebf",
            "userId": "1"
        }
    }
    r = requests.post(url, data=json.dumps(payload))
    result = json.loads(r.content)
    group.send(result["results"][0]["values"]["text"])

embed()