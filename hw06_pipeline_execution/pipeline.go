package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in

	for _, stage := range stages {
		out = handleStage(done, stage, out)
	}

	return out
}

func handleStage(done In, stage Stage, in In) Out {
	outCh := make(Bi)

	go func() {
		defer func() {
			close(outCh)
			<-in // for TestAllStageStop
		}()

		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}

				select {
				case <-done:
					return
				case outCh <- v:
				}
			}
		}
	}()

	return stage(outCh)
}
