import json
import sqlite3
import time
from botoy import Action, FriendMsg, GroupMsg, EventMsg
from botoy import decorators as deco
from botoy.parser import group as gp # 群消息(GroupMsg)相关解析
from botoy.parser import friend as fp # 好友消息(FriendMsg)相关解析
from botoy.parser import event as ep # 事件(EevntMsg)相关解析
from botoy.session import ctx, session
from botoy.decorators import equal_content, ignore_botself
from botoy.session import SessionHandler, ctx, session


conn = sqlite3.connect('data.db')

@deco.ignore_botself
def receive_friend_msg(ctx: FriendMsg):
    action = Action(ctx.CurrentQQ)
    if ctx.Content == "投稿":
        action.sendFriendText(ctx.FromUin, "请将文字和图片合成一个消息发给我，目前只支持一条消息投稿:)")
        session.set('tougao',1)
    if session.get('tougao') == 1:
        text = json.loads(ctx.Content)
        content = '图片投稿，没有文字'
        if text.has_key('content'):
            content = text['content']
        pic_data = gp.pic(ctx)
        pic_list = []
        if pic_data is not None:
            for pic in pic_data.GroupPic:
                pic_list.append(pic.Url)
                print(pic.Url)
        pic_json = json.dumps(pic_list)
        conn = sqlite3.connect('data.db')
        c = conn.cursor()
        time_data = time.strftime("%Y-%m-%d %H:%M:%S", time.localtime()) 
        query = "INSERT INTO shuoshuo (owner,content,pic,status,admin,time) VALUES (?, ?, ?, ?, ?, ?)"
        value = (ctx.FromUin,content,pic_json ,0,0,time_data)
@deco.ignore_botself
@deco.from_these_groups(753471899)
def receive_group_msg(ctx: GroupMsg):
    action = Action(ctx.CurrentQQ)
    if ctx.Content == "投稿":
        pass
        session.set
    if ctx.MsgType == 'PicMsg':
        text = json.loads(ctx.Content)
        print(text['Content'])
        pic_data = gp.pic(ctx)
        if pic_data is not None:
            for pic in pic_data.GroupPic:
                print(pic.Url)
        # action.sendGroupMultiPic(ctx.FromGroupId, ctx.Content['Url'],"图")




def search(word):
    return f"查询到关于【{word}】的内容如下..."


receive_handler = SessionHandler(
    ignore_botself,
    equal_content("投稿"),
    single_user=True,
    expiration=1 * 60,
).receive_friend_msg()


@receive_handler.handle
def _():
    session.send_text("请将文字和图片合成一个消息发给我，目前只支持一条消息投稿:)")
    while True:
        word = session.pop("word", wait=True)
        if word == "退出":
            session.send_text("好哒~")
            receive_handler.finish()
        session.send_text(search(word))
