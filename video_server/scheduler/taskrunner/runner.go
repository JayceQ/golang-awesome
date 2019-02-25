package taskrunner

type Runner struct {
	Controller controlChan
	Error      controlChan
	Data       dataChan
	dataSize   int
	longLived  bool
	Dispatcher fn
	Executor   fn
}

func NewRunner(dataSize int, longLived bool, dispatcher fn, executor fn) *Runner {
	return &Runner{
		Controller: make(chan string,1),
		Error: make(chan string,1),
		Data: make(chan interface{},dataSize),
		dataSize: dataSize,
		longLived: longLived,
		Dispatcher: dispatcher,
		Executor: executor}
}

func (r *Runner) startDispatch() {
	defer func() {
		if !r.longLived {
			close(r.Controller)
			close(r.Data)
			close(r.Error)
		}
	}()

	for {
		select {
		case c := <-r.Controller:
			if c == READY_TO_DISPATCH{
				err := r.Dispatcher(r.Data)
				if err != nil {
					r.Error <- CLOSE
				}else {
					r.Controller <- READY_TO_EXECUTE
				}
			}
			if c == READY_TO_EXECUTE{
				err := r.Executor(r.Data)
				if err != nil {
					r.Error <- CLOSE
				}else {
					r.Controller <- READY_TO_DISPATCH
				}
			}
		case e := <-r.Error:
			if e == CLOSE{
				return
			}
		}
	}
}

func (r *Runner) StartAll(){
	r.Controller <- READY_TO_DISPATCH
	r.startDispatch()
}