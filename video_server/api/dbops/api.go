package dbops

import (
	"database/sql"
	"golang-awesome/video_server/api/defs"
	"golang-awesome/video_server/api/utils"
	"log"
	"time"
)

func AddUserCredential(loginName string, pwd string) error {
	stmtIns, err := dbConn.Prepare("INSERT INTO users(login_name,pwd) values (?,?)")
	if err != nil {
		return err
	}

	_, err = stmtIns.Exec(loginName, pwd)
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	return nil
}

func GetUserCredential(loginName string) (string, error) {
	stmtOut, err := dbConn.Prepare("select pwd from users where login_name = ?")
	if err != nil {
		log.Printf("%s", err)
		return "", err
	}

	var pwd string
	err = stmtOut.QueryRow(loginName).Scan(&pwd)
	if err != nil {
		return "", nil
	}

	defer stmtOut.Close()
	return pwd, nil
}

func DeleteUser(loginName, pwd string) (err error) {
	stmtDel, err := dbConn.Prepare("DELETE  FROM users WHERE  login_name = ? and pwd = ?")
	if err != nil {
		log.Printf("Delete user error: %s", err)
		return
	}

	_, err = stmtDel.Exec(loginName, pwd)
	if err != nil {
		log.Printf("Delete user error: %s", err)
		return
	}

	defer stmtDel.Close()
	return
}

func GetUser(loginName string) (user *defs.User, err error) {
	stmtOut, err := dbConn.Prepare("select id, pwd from users where login_name = ?")
	if err != nil {
		log.Printf("%s", err)
		return
	}

	var id int
	var pwd string
	err = stmtOut.QueryRow(loginName).Scan(&id, &pwd)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("%s", err)
		return
	}
	if err == sql.ErrNoRows {
		return
	}

	user = &defs.User{Id: id, LoginName: loginName, Pwd: pwd}
	defer stmtOut.Close()
	return
}

func AddNewVideo(aid int, name string) (video *defs.VideoInfo, err error) {
	vid, err := utils.NewUUID()
	if err != nil {
		return
	}
	t := time.Now()
	ctime := t.Format("Jan 02 2006, 15:04:05")
	stmtIns, err := dbConn.Prepare(`insert into video_info 
			(id, author_id, name, display_ctime) values (?,?,?,?)`)
	if err != nil {
		log.Printf("%s", err)
		return
	}
	_, err = stmtIns.Exec(vid, aid, name, ctime)
	if err != nil {
		return
	}
	video = &defs.VideoInfo{Id: vid, AuthorId: aid, Name: name, DisplayTime: ctime}
	defer stmtIns.Close()
	return
}

func GetVideoInfo(vid string) (video *defs.VideoInfo, err error) {
	stmtOut, err := dbConn.Prepare("select author_id, name, display_ctime from video_info where vid = ?")
	if err != nil {
		return
	}

	var aid int
	var dct, name string
	err = stmtOut.QueryRow(vid).Scan(&aid, &name, &dct)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}
	defer stmtOut.Close()
	video = &defs.VideoInfo{
		Id:          vid,
		AuthorId:    aid,
		Name:        name,
		DisplayTime: dct,
	}
	return
}

func ListVideoInfo(uname string, from, to int)(videos []*defs.VideoInfo, err error){
	stmtOut, err := dbConn.Prepare(`select v.id, v.author_id, v.name,v.display_ctime from video_info v 
				inner join users u on v.author_id = u.id where u.login_name = ? and v.create_time > FROM_UNIXTIME(?) 
				and v.create_time <= FROM_UNIXTIME(?) order by v.create_time desc `)
	if err != nil {
		log.Printf("%s",err)
		return 
	}

	rows, err := stmtOut.Query(uname, from, to)
	if err != nil {
		log.Printf("%s",err)
		return
	}

	for rows.Next() {
		var id, name, ctime string
		var aid int
		if err := rows.Scan(&id, &aid, &name, &ctime); err != nil {
			log.Printf("%s",err)
			return
		}
		video := &defs.VideoInfo{
			Id:id,
			AuthorId:aid,
			Name:name,
			DisplayTime:ctime,
		}
		videos = append(videos,video)
	}
	defer stmtOut.Close()
	return
}

func DeleteVideoInfo(vid string) (err error){
	stmtDel, err := dbConn.Prepare("delete from video_info where id = ?")
	if err != nil {
		log.Printf("%s",err)
		return
	}

	_, err  = stmtDel.Exec(vid)
	if err != nil {
		log.Printf("%s",err)
		return
	}

	defer stmtDel.Close()
	return
}

func AddNewComments(vid string, aid int,content string) (err error){
	id, err := utils.NewUUID()
	if err != nil {
		log.Printf("%s",err)
		return
	}

	stmtIns, err := dbConn.Prepare("insert into comments(id, video_id, author_id, content) values (?,?,?,?)")
	if err != nil {
		log.Printf("%s",err)
		return
	}

	_, err = stmtIns.Exec(id,vid,aid,content)
	if err != nil {
		log.Printf("%s",err)
		return
	}
	defer stmtIns.Close()
	return
}

func ListComments(vid string, from ,to int)(comments []*defs.Comment, err error){
	stmtOut, err := dbConn.Prepare(`select c.id, u.login_name, c.content from comments c
				inner join users u on c.author_id = u.id where c.video_id = ? and c.time >
				FROM_UNIXTIME(?) and c.time <= FROM_UNIXTIME(?) order by c.time desc`)

	rows, err := stmtOut.Query(vid, from, to)
	if err != nil {
		log.Printf("%s",err)
		return
	}

	for rows.Next() {
		var id, name, content string
		if err := rows.Scan(&id, &name, &content); err != nil {
			return
		}
		comment := &defs.Comment{
			Id:id,
			VideoId:vid,
			Author:name,
			Content:content,
		}
		comments = append(comments,comment)
	}
	defer stmtOut.Close()
	return
}