package main

import (
	"github.com/mcoo/OPQBot"
	"github.com/mcoo/OPQBot/qzone"
	"log"
)

var (
	picArray []string
)

func main() {
	opqBot := OPQBot.NewBotManager(, "")

	err := opqBot.Start()
	if err != nil {
		log.Println(err.Error())
	}
	defer opqBot.Stop()

	//纯文字
	_, err = opqBot.AddEvent(OPQBot.EventNameOnFriendMessage, func(botQQ int64, packet *OPQBot.FriendMsgPack) {
		log.Printf("%+v", packet)
		//voi := db.Voice{}
		//json.Unmarshal([]byte(packet.Content), &voi)
		//log.Println(voi.URL)

		s := opqBot.Session.SessionStart(packet.FromUin)
		last, _ := s.GetString("last")
		if packet.Content == "#开始投稿" {
			opqBot.Send(OPQBot.SendMsgPack{
				SendToType: OPQBot.SendToTypeFriend,
				ToUserUid:  packet.FromUin,
				Content:    OPQBot.SendTypeTextMsgContent{Content: OPQBot.MacroAt([]int64{packet.FromUin}) + "请将你的投稿内容以【文字】的形式【一次性】发送给我:)"},
			})
			//欢迎投稿
			opqBot.Send(OPQBot.SendMsgPack{
				SendToType: OPQBot.SendToTypeFriend,
				ToUserUid:  packet.FromUin,
				Content:    OPQBot.SendTypeVoiceByUrlContent{VoiceUrl: "http://203.205.248.81:80/?ver=2&rkey=3062020101045b30590201010201000204c9d2765b042431306a336c46796c4e67316e45485f784c63754178716c5334763733665461387079717802046404066a041f0000000866696c6574797065000000013000000005636f64656300000001310400&voice_codec=1&filetype=0"},
			})

			s.Set("last", "#开始投稿")
		} else if last == "#开始投稿" {
			opqBot.Send(OPQBot.SendMsgPack{
				SendToType: OPQBot.SendToTypeFriend,
				ToUserUid:  packet.FromUin,
				Content:    OPQBot.SendTypeTextMsgContent{Content: OPQBot.Face_献吻 + "投稿完成:)"},
			})
			s.Set("last", "")
			ck, _ := opqBot.GetUserCookie()
			qz := qzone.NewQzoneManager(opqBot.QQ, ck)
			res, err := qz.SendShuoShuo(packet.Content)
			if err != nil {
				log.Println("发空间失败：", err)
			}
			log.Println("发空间反馈：", res)
			opqBot.Send(OPQBot.SendMsgPack{
				SendToType: OPQBot.SendToTypeFriend,
				ToUserUid:  packet.FromUin,
				Content:    OPQBot.SendTypeVoiceByUrlContent{VoiceUrl: "http://203.205.248.47:80/?ver=2&rkey=3062020101045b30590201010201000204c9d2765b042431306a334845796c4e67316e45485f334a75783334733849306566667166774e383147520204640406ae041f0000000866696c6574797065000000013000000005636f64656300000001310400&voice_codec=1&filetype=0"},
			})
		}

	})
	opqBot.Wait()
}
