package walker

import "io/fs"

// Scan recursively scans a directory tree, and returns all nodes matching the Visit function.
// Nodes returned are first sorted descending by score, then by lexicographical order.
// When an error occurs, may continue scanning until all units have exited and returns nil, err.
//
// This function is a convenience alternative to:
//
//	 scanner := Walker{Visit: Visit, Params: Params}
//	 err := scanner.Walk();
//		results := scanner.Results()
func Scan(Visit ScanProcess, Params Params) ([]string, error) {
	scanner := Walker[struct{}]{
		Process: Visit,
		Params:  Params,
	}

	err := scanner.Walk()

	// we can safely access results directly
	// because now the walker becomes inaccessible!
	return scanner.rpaths, err
}

// ScanProcess is a function that is called once for each directory that is being walked.
// It returns a triple of float64 score, bool continue and err error.
//
// match indicates that what score the path received.
// A non-negative score indicates a match, and will be returned in the array from Scan().
// cont indicates if Scan() should continue scanning recursively.
// err != nil indicates that an error has occurred, and the entire process should be aborted.
//
// ScanProcess may be nil.
// In such a case, it is assumed to return (0, true, nil) for every invocation.
//
// ScanProcess implements Process and can be used with Walk
type ScanProcess func(path string, root FS, depth int) (score float64, cont bool, err error)

func (v ScanProcess) Visit(context WalkContext[struct{}]) (shouldVisitChildren bool, err error) {

	// implement v == nil case
	if v == nil {
		context.Mark(0)
		return true, nil
	}

	// call the match function (safely)
	match, cont, err := v(context.NodePath(), context.Root(), context.Depth())
	if err != nil {
		return false, err
	}
	if match >= 0 {
		context.Mark(match)
	}
	return cont, nil
}

func (ScanProcess) VisitChild(child fs.DirEntry, valid bool, context WalkContext[struct{}]) (action Step, err error) {
	return DoConcurrent, nil
}

func (ScanProcess) AfterVisitChild(child fs.DirEntry, resultValue any, resultOK bool, context WalkContext[struct{}]) (err error) {
	return nil
}

func (ScanProcess) AfterVisit(context WalkContext[struct{}]) (err error) {
	return nil
}

// ScanMatch can be used to implement a boolean scan process.
// When value is true, it returns 1, when it is false, it returns -1.
func ScanMatch(value bool) float64 {
	if value {
		return 1
	}
	return -1
}
