package program

// Keywords are special "commands" that manipulate arguments before execution.
//
// Keywords can not be stopped by calls to universal flags; they are expanded once before aliases and command expansion takes place.
type Keyword func(args *Arguments) error

// RegisterKeyword registers a new keyword.
// See also Keyword.
//
// If an keyword already exists, RegisterKeyword calls panic().
func (p *Program[Runtime, Parameters, Requirements]) RegisterKeyword(name string, keyword Keyword) {
	if p.keywords == nil {
		p.keywords = make(map[string]Keyword)
	}

	if _, ok := p.keywords[name]; ok {
		panic("RegisterKeyword(): Keyword already registered")
	}

	p.keywords[name] = keyword
}
