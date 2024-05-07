package repository

import (
	"sync"

	"k8s.io/klog/v2"
)

type Pool struct {
	WorkList []Job
	JobQueue chan Job

	Wait       sync.WaitGroup
	NumWorkers int
	BachSize   int
}

func NewPool(num int, bachSize int) *Pool {
	return &Pool{
		JobQueue:   make(chan Job, bachSize),
		NumWorkers: num,
	}
}

func (p *Pool) Listener() {
	for i := 0; i < p.NumWorkers; i++ {
		go func(i int) {
			klog.Infof("Worker %d started", i)
			for job := range p.JobQueue {
                job.Run()
                p.Wait.Done()
			}
		}(i)
	}
}

func (p *Pool) Start() {
	for _, job := range p.WorkList {
		p.JobQueue <- job
	}
}

func (p *Pool) AddJob(job Job) {
	p.WorkList = append(p.WorkList, job)
}
