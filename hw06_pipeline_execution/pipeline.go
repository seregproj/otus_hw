package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func stageWrapper(in In, done In, stage Stage) Out {
	myIn := make(Bi)
	myOut := make(Bi)
	closeInCh := false

	gracefulShutdown := func(in In, out In) {
		if !closeInCh {
			close(myIn)
		}

		for range in {
		}
		for range out {
		}
	}

	go func() {
		defer close(myOut)
		out := stage(myIn)

		for {
			select {
			case <-done:
				gracefulShutdown(in, out)

				return
			default:
			}

			select {
			case val, ok := <-in:
				if !ok {
					closeInCh = true
					close(myIn)

					break
				}

				myIn <- val

			case <-done:
				gracefulShutdown(in, out)

				return
			}

			select {
			case val, ok := <-out:
				if !ok {
					return
				}

				myOut <- val
			case <-done:
				gracefulShutdown(in, out)

				return
			}
		}
	}()

	return myOut
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	myCh := in

	for _, stage := range stages {
		myCh = stageWrapper(myCh, done, stage)
	}

	return myCh
}
