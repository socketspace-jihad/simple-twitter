package worker

import "errors"

var workers map[string]Worker = make(map[string]Worker)

type Worker interface {
	Consume()
}

func RegisterWorker(name string, val Worker) {
	workers[name] = val
}

func GetWorker(name string) (Worker, error) {
	w, ok := workers[name]
	if !ok {
		return nil, errors.New("worker name not found")
	}
	return w, nil
}

func RunAll() {
	for _, w := range workers {
		go w.Consume()
	}
}
