package main

import (
	"errors"
	"fmt"
	"github.com/mcoo/OPQBot"
	"github.com/mcoo/OPQBot/qzone"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

var (
	picArray []string
)

type shuoshuo struct {
	gorm.Model
	Owner    string
	Content  string
	Pic      string
	Admin_qq string
	Status   string
}

func main() {

	opqBot := OPQBot.NewBotManager(3386013275, "")

	err := opqBot.Start()
	if err != nil {
		log.Println(err.Error())
	}
	db, err := gorm.Open(
		sqlite.Open("data.db"),
		&gorm.Config{},
	)
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&shuoshuo{})
	defer opqBot.Stop()

	//纯文字
	_, err = opqBot.AddEvent(OPQBot.EventNameOnFriendMessage, func(botQQ int64, packet *OPQBot.FriendMsgPack) {
		log.Printf("%+v", packet)
		s := opqBot.Session.SessionStart(packet.FromUin)
		last, _ := s.GetString("last")
		if packet.Content == "投稿" {
			opqBot.Send(OPQBot.SendMsgPack{
				SendToType: OPQBot.SendToTypeFriend,
				ToUserUid:  packet.FromUin,
				Content:    OPQBot.SendTypeTextMsgContent{Content: OPQBot.MacroAt([]int64{packet.FromUin}) + "请将你的投稿内容以【文字】的形式【一次性】发送给我:)"},
			})
			s.Set("last", "投稿")
		} else if last == "投稿" && packet.MsgType == "TextMsg" {
			time_t := time.Now().Format("2006-01-02 15:04:05")
			fmt.Println("Curret Time: ", time_t)

			db.Create(&shuoshuo{Owner: strconv.FormatInt(packet.FromUin, 10), Content: packet.Content, Pic: "", Admin_qq: "", Status: "0"})
			opqBot.Send(OPQBot.SendMsgPack{
				SendToType: OPQBot.SendToTypeFriend,
				ToUserUid:  packet.FromUin,
				Content:    OPQBot.SendTypeTextMsgContent{Content: OPQBot.Face_献吻 + "投稿完成，等待管理员审核:)"},
			})
			s.Set("last", "")

			var count int64
			db.Model(&shuoshuo{}).Where("status = ?", "0").Count(&count)
			opqBot.Send(OPQBot.SendMsgPack{
				SendToType: OPQBot.SendToTypeGroup,
				ToUserUid:  753471899,
				Content:    OPQBot.SendTypeTextMsgContent{Content: "有人投稿，请及时处理~\n当前未审核：" + strconv.FormatInt(count, 10)},
			})
			//ck, _ := opqBot.GetUserCookie()
			//qz := qzone.NewQzoneManager(opqBot.QQ, ck)
			//res, err := qz.SendShuoShuo(packet.Content)
			//if err != nil {
			//	log.Println("发空间失败：", err)
			//}
			//log.Println("发空间反馈：", res)
		}
	})
	//审核
	_, err = opqBot.AddEvent(OPQBot.EventNameOnGroupMessage, func(botQQ int64, packet *OPQBot.GroupMsgPack) {
		if packet.FromGroupID == 753471899 {
			s := opqBot.Session.SessionStart(packet.FromGroupID)
			t, _ := s.GetString("ing?")
			if packet.Content == "上班" {
				if t == "yes" {
					opqBot.Send(OPQBot.SendMsgPack{
						SendToType: OPQBot.SendToTypeGroup,
						ToUserUid:  packet.FromGroupID,
						Content:    OPQBot.SendTypeTextMsgContent{Content: OPQBot.MacroAt([]int64{packet.FromUserID}) + "已经有人在上班了~"},
					})
				} else {
					s.Set("ing?", "yes")
					opqBot.Send(OPQBot.SendMsgPack{
						SendToType: OPQBot.SendToTypeGroup,
						ToUserUid:  packet.FromGroupID,
						Content:    OPQBot.SendTypeTextMsgContent{Content: OPQBot.MacroAt([]int64{packet.FromUserID}) + "开始上班了:)请发：“来”获取待审核的稿件~"},
					})
				}
			} else if packet.Content == "来" && t == "yes" {
				fmt.Println("shangban:", packet.Content)
				var tt shuoshuo
				err := db.First(&tt, "status = ?", "0").Error
				if errors.Is(err, gorm.ErrRecordNotFound) {
					s.Set("ing?", "no")
					opqBot.Send(OPQBot.SendMsgPack{
						SendToType: OPQBot.SendToTypeGroup,
						ToUserUid:  packet.FromGroupID,
						Content:    OPQBot.SendTypeTextMsgContent{Content: OPQBot.MacroAt([]int64{packet.FromUserID}) + "已经没有待审核的稿件了，提前下班~"},
					})
				} else {
					fmt.Println("数据库：", tt)
					s.Set("shuoshuo", strconv.FormatInt(int64(tt.ID), 10))
					opqBot.Send(OPQBot.SendMsgPack{
						SendToType: OPQBot.SendToTypeGroup,
						ToUserUid:  packet.FromGroupID,
						Content: OPQBot.SendTypeTextMsgContent{Content: "[序号]" + strconv.FormatInt(int64(tt.ID), 10) + "\n[来自]" + tt.Owner +
							"\n[内容]" + tt.Content + "\n[时间]" + tt.CreatedAt.Format("2006-01-02 15:04:05") + "\n1通过，0拒绝"},
					})
				}
			} else if packet.Content == "1" && t == "yes" {
				shuo, _ := s.GetString("shuoshuo")
				if shuo != "" {
					opqBot.Send(OPQBot.SendMsgPack{
						SendToType: OPQBot.SendToTypeGroup,
						ToUserUid:  packet.FromGroupID,
						Content:    OPQBot.SendTypeTextMsgContent{Content: "[序号]" + shuo + "\n[状态]通过"},
					})
					var tt shuoshuo
					err := db.First(&tt, shuo).Error
					if !errors.Is(err, gorm.ErrRecordNotFound) {
						tt.Admin_qq = strconv.FormatInt(packet.FromUserID, 10)
						tt.Status = "1"
						db.Model(&tt).Select("admin_qq", "status").Updates(&tt)
					}
					s.Set("shuoshuo", "")
					ck, _ := opqBot.GetUserCookie()
					qz := qzone.NewQzoneManager(opqBot.QQ, ck)
					res, err := qz.SendShuoShuo(tt.Content)
					if err != nil {
						log.Println("发空间失败：", err)
					}
					log.Println("发空间反馈：", res)
				} else {
					opqBot.Send(OPQBot.SendMsgPack{
						SendToType: OPQBot.SendToTypeGroup,
						ToUserUid:  packet.FromGroupID,
						Content:    OPQBot.SendTypeTextMsgContent{Content: "当前没有查看任何稿件~"},
					})
				}
			} else if packet.Content == "0" && t == "yes" {
				shuo, _ := s.GetString("shuoshuo")
				if shuo != "" {
					opqBot.Send(OPQBot.SendMsgPack{
						SendToType: OPQBot.SendToTypeGroup,
						ToUserUid:  packet.FromGroupID,
						Content:    OPQBot.SendTypeTextMsgContent{Content: "[序号]" + shuo + "\n[状态]不通过"},
					})
					var tt shuoshuo
					err := db.First(&tt, shuo).Error
					if !errors.Is(err, gorm.ErrRecordNotFound) {
						tt.Admin_qq = strconv.FormatInt(packet.FromUserID, 10)
						tt.Status = "-1"
						db.Model(&tt).Select("admin_qq", "status").Updates(&tt)
					}
					s.Set("shuoshuo", "")
				} else {
					opqBot.Send(OPQBot.SendMsgPack{
						SendToType: OPQBot.SendToTypeGroup,
						ToUserUid:  packet.FromGroupID,
						Content:    OPQBot.SendTypeTextMsgContent{Content: "当前没有查看任何稿件~"},
					})
				}
			}
		}
	})
	opqBot.Wait()
}
