package service

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"errors"
)

var sciFiAuthors = []string{
	"Isaac Asimov",
	"Arthur C. Clarke",
	"Philip K. Dick",
	"Ursula K. Le Guin",
	"Ray Bradbury",
	"Robert A. Heinlein",
	"Frank Herbert",
	"William Gibson",
	"Octavia E. Butler",
	"H. G. Wells",
	"Aldous Huxley",
	"Jules Verne",
	"Neal Stephenson",
	"Orson Scott Card",
	"Kim Stanley Robinson",
	"James S. A. Corey",
	"Anne McCaffrey",
	"Stanislaw Lem",
	"Frederik Pohl",
	"Brian W. Aldiss",
	"Greg Bear",
	"John Scalzi",
	"Vernor Vinge",
	"Cory Doctorow",
	"China Miéville",
	"Joe Haldeman",
	"Peter F. Hamilton",
	"Samuel R. Delany",
	"Alastair Reynolds",
	"Michael Moorcock",
	"Stephen Baxter",
	"Larry Niven",
	"David Brin",
	"Dan Simmons",
	"Lois McMaster Bujold",
	"Charles Stross",
	"Gene Wolfe",
	"Mary Shelley",
	"Kurt Vonnegut",
	"Yoon Ha Lee",
	"Becky Chambers",
	"Andy Weir",
	"N. K. Jemisin",
	"Ken Liu",
	"Adrian Tchaikovsky",
	"Ann Leckie",
	"Ramez Naam",
	"R. F. Kuang",
	"Naomi Novik",
	"Elizabeth Moon",
	"Ted Chiang",
	"Richard K. Morgan",
	"John Wyndham",
	"Roger Zelazny",
	"Poul Anderson",
	"Olaf Stapledon",
	"Philip José Farmer",
	"Michael Flynn",
	"Bruce Sterling",
	"Joan D. Vinge",
	"James Blish",
	"Connie Willis",
	"Gene Roddenberry",
	"David Weber",
	"Iain M. Banks",
	"Christopher Priest",
	"Hannu Rajaniemi",
	"Margaret Atwood",
	"Tamsyn Muir",
	"Malka Older",
	"Marissa Meyer",
	"Scott Westerfeld",
	"Douglas Adams",
	"Ernest Cline",
	"Patrick Ness",
	"Lev Grossman",
	"Carl Sagan",
	"Greg Egan",
	"Charlie Jane Anders",
	"Jeff VanderMeer",
	"Tim Powers",
	"Kage Baker",
	"David Zindell",
	"Paul McAuley",
	"Alexander Jablokov",
	"Rudy Rucker",
	"J. G. Ballard",
	"Jack Vance",
	"John Varley",
	"Tom Godwin",
	"Nancy Kress",
}

type NameGenerator interface {
	NameFor(any) (string, error)
}

func NewNameGenerator() NameGenerator {
	return &nameGenerator{}
}

type nameGenerator struct{}

func (n *nameGenerator) NameFor(key any) (string, error) {
	if key == nil {
		return "", errors.New("key cannot be nil")
	}

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(key); err != nil {
		return "", err
	}

	hash := sha256.Sum256(buf.Bytes())

	idx := int(binary.BigEndian.Uint64(hash[:8]) % uint64(len(sciFiAuthors)))
	return sciFiAuthors[idx], nil
}
