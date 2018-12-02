package model

import (
	"fmt"
	"time"

	"encoding/json"

	"github.com/garyburd/redigo/redis"
)

var (
	UserTable = "users"
)

type UserMgr struct {
	pool *redis.Pool
}

func NewUserMgr(pool *redis.Pool) (mgr *UserMgr) {

	mgr = &UserMgr{
		pool: pool,
	}
	return
}

func (p *UserMgr) getUser(conn redis.Conn, id int) (user *User, err error) {

	result, err := redis.String(conn.Do("HGet", UserTable, fmt.Sprintf("%d", id)))
	if err != nil {
		if err == redis.ErrNil {
			err = ErrUserNotExist
		}
		return
	}

	user = &User{}
	err = json.Unmarshal([]byte(result), user)
	if err != nil {
		return
	}
	return
}

func (p *UserMgr) Login(id int, passwd string) (user *User, err error) {

	conn := p.pool.Get()
	defer conn.Close()

	user, err = p.getUser(conn, id)
	if err != nil {
		return
	}

	if user.UserId != id || user.Passwd != passwd {
		err = ErrInvalidPasswd
	}

	user.Status = UserStatusOnline
	user.LastLogin = fmt.Sprintf("%v", time.Now())

	return
}

func (p *UserMgr) Register(user *User) (err error) {

	conn := p.pool.Get()
	defer conn.Close()

	if user == nil {
		fmt.Println("invalid user")
		err = ErrInvalidParams
		return
	}

	_, err = p.getUser(conn, user.UserId)
	if err == nil {
		err = ErrUserExist
		return
	}

	if err != ErrUserNotExist {
		return
	}

	data, err := json.Marshal(user)
	if err != nil {
		return
	}

	_, err = conn.Do("HSet", UserTable, fmt.Sprintf("%d", user.UserId), string(data))
	if err != nil {
		return
	}
	return
}
