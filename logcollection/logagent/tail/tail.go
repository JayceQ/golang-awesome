package tail

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/hpcloud/tail"
	"sync"
	"time"
)

const (
	StatusNormal = 1
	StatusDelete = 2
)

type CollectConf struct {
	LogPath string `json:"logpath"`
	Topic   string `json:"topic"`
}

type TailObj struct {
	tail     *tail.Tail
	conf     CollectConf
	status   int
	exitChan chan int
}

type TextMsg struct {
	Msg   string
	Topic string
}
type TailObjMgr struct {
	tailObjs []*TailObj
	msgChan  chan *TextMsg
	lock     sync.Mutex
}

var (
	tailObjMgr *TailObjMgr
)

func GetOneLine() (msg *TextMsg) {
	msg = <-tailObjMgr.msgChan
	return
}

func UpdateConfig(confs []CollectConf) (err error) {
	tailObjMgr.lock.Lock()
	defer tailObjMgr.lock.Unlock()

	for _, oneConf := range confs {
		var isRunning = false
		for _, obj := range tailObjMgr.tailObjs {
			if oneConf.LogPath == obj.conf.LogPath {
				isRunning = true
				break
			}
		}
		if isRunning {
			continue
		}

		createNewTask(oneConf)
	}

	var tailObjs []*TailObj
	for _, obj := range tailObjMgr.tailObjs {
		obj.status = StatusDelete
		for _, oneConf := range confs {
			if oneConf.LogPath == obj.conf.LogPath {
				obj.status = StatusNormal
				break
			}
		}
		if obj.status == StatusDelete {
			obj.exitChan <- 1
			continue
		}
		tailObjs = append(tailObjs, obj)
	}
	tailObjMgr.tailObjs = tailObjs
	return
}

func createNewTask(conf CollectConf) {

	obj := &TailObj{
		conf:     conf,
		exitChan: make(chan int, 1),
	}

	file, e := tail.TailFile(conf.LogPath, tail.Config{
		ReOpen: true,
		Follow: true,
		//Location: &tail.SeekInfo{Offset:0,Whence:2,},
		MustExist: false,
		Poll:      true,
	})

	if e != nil {
		logs.Error("collect filename[%s] failed ,err:%v", e)
		return
	}
	obj.tail = file
	tailObjMgr.tailObjs = append(tailObjMgr.tailObjs, obj)

	go readFromTail(obj)
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
