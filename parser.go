package anitomy

import "sort"

type Parser struct {
	Tokens   []Token
	Elements []Element
}

func NewParser(tokens []Token) *Parser {
	return &Parser{Tokens: tokens}
}

func (p *Parser) Parse(options Options) {
	// File extension
	if options.ParseFileExtension {
		if e := ParseFileExtension(p.Tokens); e != nil {
			p.addElement(*e)
		}
	}

	// Keywords
	p.addElements(ParseKeywords(p.Tokens, options))

	// Checksum
	if options.ParseFileChecksum {
		if e := ParseFileChecksum(p.Tokens); e != nil {
			p.addElement(*e)
		}
	}

	// Video resolution
	if options.ParseVideoResolution {
		p.addElements(ParseVideoResolution(p.Tokens))
	}

	// Year
	if options.ParseYear {
		if e := ParseYear(p.Tokens); e != nil {
			p.addElement(*e)
		}
	}

	// Season
	if options.ParseSeason {
		p.addElements(ParseSeason(p.Tokens))
	}

	// Part
	if options.ParsePart {
		if e := ParsePart(p.Tokens); e != nil {
			p.addElement(*e)
		}
	}

	// Episode
	if options.ParseEpisode {
		p.addElements(ParseVolume(p.Tokens))
		p.addElements(ParseEpisode(p.Tokens))
	}

	// Title
	if options.ParseTitle {
		if e := ParseTitle(p.Tokens); e != nil {
			p.addElement(*e)
		}
	}

	// Release group
	if options.ParseReleaseGroup && !p.Contains(ReleaseGroup) {
		if e := ParseReleaseGroup(p.Tokens); e != nil {
			p.addElement(*e)
		}
	}

	// Episode title
	if options.ParseEpisodeTitle && p.Contains(Episode) {
		if e := ParseEpisodeTitle(p.Tokens); e != nil {
			p.addElement(*e)
		}
	}

	sort.Slice(p.Elements, func(i, j int) bool {
		return p.Elements[i].Position < p.Elements[j].Position
	})
}

func (p *Parser) addElement(element Element) {
	p.Elements = append(p.Elements, element)
}

func (p *Parser) addElements(elements []Element) {
	p.Elements = append(p.Elements, elements...)
}

func (p *Parser) Contains(kind ElementKind) bool {
	for _, e := range p.Elements {
		if e.Kind == kind {
			return true
		}
	}
	return false
}
