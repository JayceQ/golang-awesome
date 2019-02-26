package taskrunner

import (
	"errors"
	"golang-awesome/video_server/scheduler/dbops"
	"log"
	"os"
	"path/filepath"
	"sync"
)

func deleteVideo(vid string) error {
	path, _ := filepath.Abs(VIDEO_PATH + vid)
	log.Panicln(path)
	err := os.Remove(VIDEO_PATH + vid)
	if err != nil {
		log.Printf("deleting video error: %v", err)
		return err
	}
	return nil
}

func VideoClearDispathcer(dc dataChan) error {
	res, err := dbops.ReadVideoDeletionRecord(3)
	if err != nil {
		log.Printf("video clear dispatcher error: %s", err)
	}
	if len(res) == 0 {
		return errors.New("all tasks finished")
	}

	for _, id := range res {
		dc <- id
	}
	return nil
}

func VideoClearExecutor(dc dataChan) error {
	errMap := &sync.Map{}
	var err error
loop:
	for {
		select {
		case vid := <-dc:
			go func(id interface{}) {
				if err := deleteVideo(id.(string)); err != nil {
					errMap.Store(id, err)
					return
				}
				if err := dbops.DelVideoDeletionRecord(id.(string)); err != nil {
					errMap.Store(id, err)
					return
				}
			}(vid)
		default:
			break loop
		}
	}
	errMap.Range(func(k, v interface{}) bool {
		err = v.(error)
		if err != nil {
			return false
		}
		return true
	})
	return err
}
