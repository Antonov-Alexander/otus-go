package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	proxy := func(in In, done In) Out {
		write := make(Bi)
		go func() {
			for value := range in {
				if done != nil {
					select {
					case _, ok := <-done:
						if !ok {
							close(write)
							return
						}
					default:
					}
				}

				write <- value
			}
			close(write)
		}()

		return write
	}

	if len(stages) == 0 {
		closedChan := make(Bi)
		close(closedChan)
		return proxy(in, closedChan)
	}

	for _, stage := range stages {
		in = stage(proxy(in, done))
	}

	return in
}
