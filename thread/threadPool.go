package thread

import "sync"


type Task func()

type ThreadPool struct{
    tasks       chan Task
    wg          sync.WaitGroup
    done        chan struct{}
}

// Create a New thread pool with workes to take submitted task.
// @param worker: Number of worker in the pool to do task async  
func NewThreadPool(threadCount int) *ThreadPool{
    pool := &ThreadPool{tasks: make(chan Task, 10), done: make(chan struct{})}

    for range(threadCount){
        go pool.worker()
    }
    return pool
}

func (tp *ThreadPool) worker(){
    defer tp.wg.Done()
    for {
        select{
        case task := <-tp.tasks:
            task()
            return
        case <- tp.done:
            return
        }
    }
}

// Submit a task. A task is a lambda function
// This will be blocking if  there are no workers avaiable 
func (p *ThreadPool) Submit(task Task) {
	p.wg.Add(1)
	p.tasks <- task
}

// Wait for workes to complete their work. This is a blocking code
func (p *ThreadPool) Wait() {
	p.wg.Wait()
}

// Closes the thread pool and wait for all workes to complete their work (Blocking)
func (p *ThreadPool) Shutdown() {
    p.wg.Wait()
	close(p.tasks)
	close(p.done)
}
