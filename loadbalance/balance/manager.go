package balance

import "fmt"

type BalanceMgr struct {
	allBalancer map[string]Balance
}



var mgr =  BalanceMgr{
	allBalancer:make(map[string]Balance),
}

func (p *BalanceMgr) registerBalance(name string,b Balance){
	p.allBalancer[name] = b
}

func RegisterBalance(name string,b Balance){
	mgr.registerBalance(name,b)
}

func DoBalance(name string,insts []*Instance) (inst *Instance,err error){
	balance, ok := mgr.allBalancer[name]

	if !ok {
		err = fmt.Errorf("not found %s balance",name)
		return
	}

	fmt.Printf("use %s balance \n",name)
	inst,err = balance.DoBalance(insts)
	return
}