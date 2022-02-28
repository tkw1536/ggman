package goprogram

// Keywords are special "commands" that manipulate arguments before execution.
//
// Keywords can not be stopped by calls to universal flags; they are expanded once before aliases and command expansion takes place.
type Keyword[Flags any] func(args *Arguments[Flags]) error

// RegisterKeyword registers a new keyword.
// See also Keyword.
//
// If an keyword already exists, RegisterKeyword calls panic().
func (p *Program[Runtime, Parameters, Flags, Requirements]) RegisterKeyword(name string, keyword Keyword[Flags]) {
	if p.keywords == nil {
		p.keywords = make(map[string]Keyword[Flags])
	}

	if _, ok := p.keywords[name]; ok {
		panic("RegisterKeyword(): Keyword already registered")
	}

	p.keywords[name] = keyword
}
