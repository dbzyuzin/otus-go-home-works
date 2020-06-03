package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	var out = in
	for _, stage := range stages {
		out = MergeWithDone(out, done)
		out = stage(out)
	}

	return out
}

func MergeWithDone(ch Out, done In) In {
	out := make(Bi)
	go func() {
		defer close(out)
		for {
			select {
			case r, ok := <-ch:
				if !ok {
					return
				}
				out <- r
			case <-done:
				return
			}
		}
	}()

	return out
}
