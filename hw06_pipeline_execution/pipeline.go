package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	checkDone := func(done In) bool {
		if done != nil {
			select {
			case _, ok := <-done:
				if !ok {
					return false
				}
			default:
			}
		}

		return true
	}

	proxy := func(readIn In, write Bi, done In) {
		for value := range readIn {
			if !checkDone(done) {
				close(write)
				return
			}

			write <- value
		}
		close(write)
	}

	stageWrapper := func(in In, done In, stage Stage) Out {
		write := make(Bi)
		go proxy(in, write, done)
		return stage(write)
	}

	for _, stage := range stages {
		in = stageWrapper(in, done, stage)
	}

	return in
}
