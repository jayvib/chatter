package thesaurus

type Thesaurus interface {
	Synonyms(term string) (syns []string, error error)
}
