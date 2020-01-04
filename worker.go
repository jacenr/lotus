package lotus

import "sync"

type WorkerManager struct {
	Num int           // the number of workers
	Ch  chan func()   // func channel from which worker get work to run
	Cl  chan struct{} // to end all worker by close chan
	WG  sync.WaitGroup
	//Cx  context.Context // to concel all wokers by context
}

func NewWorker(n int) *WorkerManager {
	wm := &WorkerManager{}
	wm.Num = n
	//if ctx != nil {
	//	wm.Cx = ctx
	//}
	c := make(chan func(), 0)
	wm.Ch = c

	cl := make(chan struct{}, 0)
	wm.Cl = cl

	wg := sync.WaitGroup{}
	wm.WG = wg

	return wm
}

func (wm *WorkerManager) StartWork() {
	for i := 0; i < wm.Num; i++ {
		wm.WG.Add(1)
		go func() {
			defer wm.WG.Done()
			for {
				select {
				case f, ok := <-wm.Ch:
					if !ok {
						return
					}
					f()
				//case <-wm.Cx.Done():
				//	return
				case <-wm.Cl:
					return
				}
			}
		}()
	}
}

func (wm *WorkerManager) ForceEndWorker() {
	close(wm.Cl)
}

func (wm *WorkerManager) EndWorkerAndWait() {
	close(wm.Ch)
	wm.WG.Wait()
}
