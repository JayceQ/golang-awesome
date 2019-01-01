package service

import "sync"

type SecLimitMgr struct {
	UserLimitMap map[int]*Limit
	IpLimitMap map[string]Limit
	lock sync.Mutex
}

