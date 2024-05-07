package repository

type Job struct {
	Do  func() (error)
	Err error
}

func (j *Job) Run() {
	j.Err = j.Do()
}

func NewJob(do func() error) *Job {
	Job := Job{
		Do: do,
	}
	return &Job
}
