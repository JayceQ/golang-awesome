package tail

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/hpcloud/tail"
	"time"
)

type CollectConf struct {
	LogPath string
	Topic   string
}

type TailObj struct {
	tail *tail.Tail
	conf CollectConf
}


type TextMsg struct {
	Msg   string
	Topic string
}
type TailObjMgr struct {
	tailObjs []*TailObj
	msgChan  chan *TextMsg
}

var (
	tailObjMgr *TailObjMgr
)

func GetOneLine()(msg *TextMsg){
	msg = <-tailObjMgr.msgChan
	return
}

func InitTail(conf []CollectConf, chanSize int) (err error) {

	if len(conf) == 0 {
		err = fmt.Errorf("invalid config for log collect: conf:%v", conf)
		return
	}

	tailObjMgr = &TailObjMgr{
		msgChan: make(chan *TextMsg, chanSize),
	}

	for _, v := range conf {
		obj := &TailObj{
			conf: v,
		}

		tails, errTail := tail.TailFile(v.LogPath, tail.Config{
			ReOpen: true,
			Follow: true,
			//Location: &tail.SeekInfo{Offset:0,Whence:2,},
			MustExist: false,
			Poll:      true,
		})

		if errTail != nil {
			err = errTail
			return
		}

		obj.tail = tails

		tailObjMgr.tailObjs = append(tailObjMgr.tailObjs, obj)

		go readFromTail(obj)
	}
	return
}

func readFromTail(tailObj *TailObj) {
	for {
		lines, ok := <-tailObj.tail.Lines
		if !ok {
			logs.Warn("tail file close reopen ,filename:%s\n", tailObj.tail.Filename)
			time.Sleep(100 * time.Millisecond)
			continue
		}

		textMsg := &TextMsg{
			Msg:   lines.Text,
			Topic: tailObj.conf.Topic,
		}

		tailObjMgr.msgChan <- textMsg
	}
}
