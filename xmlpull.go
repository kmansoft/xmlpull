package xmlpull

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"
)

/* ----- */

type Atoms interface {
	AddAtom(s string) int
}

type atomsImpl struct {
	atoms   map[string]int
	next_id int
}

func NewAtoms() Atoms {
	return &atomsImpl{
		atoms:   make(map[string]int),
		next_id: 1}
}

func (x *atomsImpl) AddAtom(s string) int {
	if len(s) == 0 {
		return 0
	}

	if id, ok := x.atoms[s]; ok {
		return id
	}

	id := x.next_id
	x.next_id++
	x.atoms[s] = id
	return id
}

/* ----- */

type Token interface {
}

type Parser interface {
	GetAtoms() Atoms
	NextToken() (Token, error)
}

type Tag struct {
	Parent  *Tag
	Name    xml.Name
	Space   int
	Local   int
	IsStart bool
	IsEnd   bool
}

func (t *Tag) IsTag(space, local int) bool {
	return t.Space == space && t.Local == local
}

func (t *Tag) IsParentTag(space, local int) bool {
	return t.Parent != nil && t.Parent.IsTag(space, local)
}

type Text struct {
	Text string
	Tag  *Tag
}

func (t *Text) IsTag(space, local int) bool {
	return t.Tag != nil && t.Tag.IsTag(space, local)
}

func (t *Text) IsParentTag(space, local int) bool {
	return t.Tag != nil && t.Tag.IsParentTag(space, local)
}

func (t *Text) AsBool() bool {
	return strings.EqualFold(t.Text, "true")
}

type parserImpl struct {
	atoms   Atoms
	decoder *xml.Decoder
	curr    *Tag
}

func NewParser(decoder *xml.Decoder) Parser {
	return &parserImpl{
		atoms:   NewAtoms(),
		decoder: decoder,
		curr:    nil}
}

func NewParserBytes(b []byte) Parser {
	return NewParser(xml.NewDecoder(bytes.NewReader(b)))
}

func (x *parserImpl) GetAtoms() Atoms {
	return x.atoms
}

func (x *parserImpl) NextToken() (Token, error) {
	for {
		t, err := x.decoder.Token()
		if err != nil && err != io.EOF {
			return nil, err
		} else if t == nil {
			return nil, nil
		}

		switch se := t.(type) {
		case xml.StartElement:
			se_space := x.atoms.AddAtom(se.Name.Space)
			se_local := x.atoms.AddAtom(se.Name.Local)

			tag := &Tag{
				Parent:  x.curr,
				Name:    se.Name,
				Space:   se_space,
				Local:   se_local,
				IsStart: true,
				IsEnd:   false}
			x.curr = tag
			return *tag, nil

		case xml.EndElement:
			se_space := x.atoms.AddAtom(se.Name.Space)
			se_local := x.atoms.AddAtom(se.Name.Local)

			if x.curr == nil || x.curr.Space != se_space || x.curr.Local != se_local {
				return nil, errors.New(fmt.Sprintf("Unmatched tag %s", se.Name))
			}

			tag := x.curr
			tag.IsStart = false
			tag.IsEnd = true
			x.curr = x.curr.Parent
			return *tag, nil

		case xml.CharData:
			return Text{
				Tag:  x.curr,
				Text: string(se)}, nil
		}
	}
}
