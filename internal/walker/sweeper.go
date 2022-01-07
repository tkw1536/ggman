package walker

import "io/fs"

// Sweep recursively sweeps a directory tree, and returns all nodes that are empty or contain only empty directories
// When an error occurs, may continue sweeping until all units have exited and returns nil, err.
//
// This function is a convenience alternative to:
//
//  scanner := &Walker{Visit: Visit, Params: Params}
//  err := scanner.Walk();
//	results := scanner.Results()
func Sweep(Visit SweepProcess, Params Params) ([]string, error) {
	scanner := &Walker{
		Process: Visit,
		Params:  Params,
	}

	err := scanner.Walk()

	// we can safely access results directly
	// because now the walker becomes inaccessible!
	return scanner.results, err
}

// SweepProcess is a function that is called once for each directory that is being sweeped.
// It returns a boolean stop.
//
// stop should indicate if the scan should continue recursively, or stop and treat the appropriate directory as non-empty.
//
// Visit may be nil.
// In such a case, it is assumed to return the pair false for every indication.
//
// SweepProcess implements Process and can be used with Walk.
type SweepProcess func(path string, root FS, depth int) (stop bool)

func (v SweepProcess) Visit(context WalkContext) (shouldRecurse bool, err error) {
	var shouldStop bool
	if v != nil {
		shouldStop = v(context.NodePath(), context.Root(), context.Depth())
	}
	if shouldStop {
		return false, nil
	}
	context.Snapshot(func(snapshot interface{}) (value interface{}) { return true })
	return true, nil // we should recurse!
}
func (SweepProcess) VisitChild(child fs.DirEntry, valid bool, context WalkContext) (action Step, err error) {
	context.Snapshot(func(snapshot interface{}) (value interface{}) {
		isEmpty := snapshot.(bool)
		if !valid {
			// non-directory => we are not empty!
			isEmpty = false
			action = DoNothing
		} else if isEmpty {
			// we have an empty directory, so we need to keep checking the rest in sync!
			action = DoSync
		} else {
			// directory is not empty, so it doesn't matter
			action = DoConcurrent
		}

		return isEmpty
	})
	return action, nil
}

func (SweepProcess) AfterVisitChild(child fs.DirEntry, resultValue interface{}, resultOK bool, context WalkContext) (err error) {
	context.Snapshot(func(snapshot interface{}) interface{} {
		isEmpty := snapshot.(bool)

		// this directory remains empty iff all child directories are empty
		if !(resultOK && resultValue.(bool)) {
			isEmpty = false
		}

		return isEmpty
	})
	return nil
}

func (SweepProcess) AfterVisit(context WalkContext) (err error) {
	context.Snapshot(func(snapshot interface{}) interface{} {
		isEmpty := snapshot.(bool)
		if isEmpty {
			context.Mark(float64(context.Depth()))
		}
		return isEmpty
	})
	return nil
}
