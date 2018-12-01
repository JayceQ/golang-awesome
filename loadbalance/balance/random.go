package balance

import (
	"errors"
	"math/rand"
)

func init(){
	RegisterBalance("random",&RandomBalance{})
}

type RandomBalance struct {
}

func (p *RandomBalance) DoBalance(insts []*Instance,key ...string) (inst *Instance,err error){

	if len(insts) == 0{
		err = errors.New("no instance")
		return
	}

	lens := len(insts)
	index := rand.Intn(lens)
	inst = insts[index]
	return
}
