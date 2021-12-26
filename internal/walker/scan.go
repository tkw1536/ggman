package walker

import "io/fs"

// Scan recursively scans a directory tree, and returns all nodes matching the Visit function.
// When an error occurs, may continue scanning until all units have exited and returns nil, err.
//
// This function is a convenience alternative to:
//
//  scanner := &Walker{Visit: Visit, Params: Params}
//  err := scanner.Walk();
//	results := scanner.Results()
func Scan(Visit ScanProcess, Params Params) ([]string, error) {
	scanner := &Walker{
		Process: Visit,
		Params:  Params,
	}

	err := scanner.Walk()
	results := scanner.Results()

	return results, err
}

// ScanProcess is a function that is called once for each directory that is being walked.
// It returns a pair of booleans match and bool.
//
// match indiciates that path should be returned in the array from Scan().
// cont indicates if Scan() should continue scanning recursively.
//
// Visit may be nil.
// In such a case, it is assumed to return the pair (true, true) for every invocation.
//
// ScanProcess implements Process and can be used with Walk
type ScanProcess func(path string, root FS, depth int) (match, cont bool)

func (v ScanProcess) Visit(context WalkContext) (shouldVisitChildren bool, err error) {
	match, cont := true, true
	if v != nil {
		match, cont = v(context.NodePath(), context.Root(), context.Depth())
	}
	if match {
		context.Mark(0)
	}
	return cont, nil
}

func (ScanProcess) VisitChild(child fs.DirEntry, valid bool, context WalkContext) (action Step, err error) {
	return DoConcurrent, nil
}

func (ScanProcess) AfterVisitChild(child fs.DirEntry, resultValue interface{}, resultOK bool, context WalkContext) (err error) {
	return nil
}

func (ScanProcess) AfterVisit(context WalkContext) (err error) {
	return nil
}
