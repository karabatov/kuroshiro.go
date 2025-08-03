package analyzer

import (
	"errors"
	"strings"

	gmecab "github.com/shogo82148/go-mecab"
)

type MecabAnalyzer struct {
	tagger  *gmecab.MeCab
	lattice *gmecab.Lattice
}

var MecabInitErr = errors.New("failed to initialize mecab")
var MecabParseErr = errors.New("failed to parse text")

func NewMecab() (*MecabAnalyzer, error) {
	tagger, err := gmecab.New(map[string]string{"output-format-type": "wakati"})
	if err != nil {
		return nil, MecabInitErr
	}

	lattice, err := gmecab.NewLattice()
	if err != nil {
		return nil, MecabInitErr
	}

	return &MecabAnalyzer{
		tagger:  &tagger,
		lattice: &lattice,
	}, nil
}

func (m *MecabAnalyzer) Destroy() {
	m.tagger.Destroy()
	m.lattice.Destroy()
}

func (m *MecabAnalyzer) ParseRaw(text string) (RawResult, error) {
	m.lattice.SetSentence(text)

	err := m.tagger.ParseLattice(*m.lattice)
	if err != nil {
		return nil, MecabParseErr
	}

	mecabResult := m.lattice.String()
	var rawResult RawResult

	lines := strings.Split(mecabResult, "\n")
	for _, line := range lines {
		arr := strings.Split(line, "\t")

		if len(arr) == 1 {
			rawResult = append(rawResult, []string{line})
			break
		}

		vals := []string{arr[0]}
		arr = strings.Split(arr[1], ",")

		vals = append(vals, arr...)

		rawResult = append(rawResult, vals)
	}

	return rawResult[0 : len(rawResult)-1], nil
}

func (m *MecabAnalyzer) Parse(text string) ([]Token, error) {
	rawResult, err := m.ParseRaw(text)
	if err != nil {
		return nil, err
	}

	/*
		# F[0]: pos1
		# F[1]: pos2
		# F[2]: pos3
		# F[3]: pos4
		# F[4]: cType
		# F[5]: cForm
		# F[6]: lForm
		# F[7]: lemma
		# F[8]: orth
		# F[9]: orthBase
		# F[10]: pron
		# F[11]: pronBase
		# F[12]: goshu
		# F[13]: aType
		# F[14]: aConType
		# F[15]: aModType
	*/

	var result []Token
	for _, r := range rawResult {
		length := len(r)

		basicForm := NUL
		pronunciation := NUL

		if length > 8 {
			basicForm = r[8]
		}

		if length > 10 {
			pronunciation = r[10]
		}

		result = append(result, Token{
			SurfaceForm:    r[0],
			Pos:            r[1],
			PosDetail1:     r[2],
			PosDetail2:     r[3],
			PosDetail3:     NUL,
			ConjugatedType: r[4],
			ConjugatedForm: r[5],
			BasicForm:      basicForm,
			Reading:        pronunciation,
			Pronunciation:  pronunciation,
		})
	}

	return result, nil
}
