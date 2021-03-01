package nlp

import (
    // natural language processing
    "sort"
    "strings"

    prose "github.com/jdkato/prose/v2"
)

type NLP struct {
    Words map[string]map[string]Entry
    Sentences []string
    Model string
}

type Entry struct {
    Value string
    Count uint
    Kind string
    Label string
}

type Result struct {
    Entries []Entry
    Words []string
}

type convertFn func(Entry) (string, string)
type sortFn func(int, int) bool

func concat(nlp NLP, tokens []string) ([]Entry, []string) {
    var entries []Entry
    var words []string
    for _, token := range tokens {
        for word, v := range nlp.Words[token] {
            entries = append(entries, v)
            words = append(words, word)
        }
    }
    return entries, words
}

func (result Result) ForEach(convert convertFn) string {
    var sb strings.Builder
    value, sep := convert(result.Entries[0])
    sb.WriteString(value)
    for _, entry := range result.Entries[1:] {
        value, sep = convert(entry)
        sb.WriteString(sep)
        sb.WriteString(value)
    }
    return sb.String()
}

func (result Result) ForEachWithSort(convert convertFn, order sortFn) string {
    sort.Slice(result.Entries, order)
    return result.ForEach(convert)
}

func (nlp NLP) Verbs() Result {
    entries, words := concat(nlp, []string{
        VERB_BASE_FORM, VERB_PAST_TENSE, VERB_GERUND_PARTICIPLE,
        VERB_PRESENT_PARTICIPLE, VERB_PAST_PARTICIPLE, 
        VERB_NON_3RD_PERSON_SINGULAR_PRESENT, VERB_3RD_PERSON_SINGULAR_PRESENT,
    })

    return Result {
        Entries: entries,
        Words: words,
    }
}

func (nlp NLP) Nouns() Result {
    entries, words := concat(nlp, []string{
        NOUN_SINGULAR, NOUN_PROPER_SINGULAR, NOUN_PROPER_PLURAL, NOUN_PLURAL,
    })

    return Result {
        Entries: entries,
        Words: words,
    }
}

func Parse(text string) (*NLP, error) {
    doc, err := prose.NewDocument(text)
    if err != nil {
        return nil, err
    }

    words := make(map[string]map[string]Entry)

    for _, tok := range doc.Tokens() {
        if words[tok.Tag] == nil {
            words[tok.Tag] = make(map[string]Entry)
        }
        entry, ok := words[tok.Tag][tok.Text]
        if ok == false  {
            entry = Entry{
                Value: tok.Text,
                Count: 1,
                Kind: tok.Tag,
                Label: tok.Label,
            }
            words[tok.Tag][tok.Text] = entry
        } else {
            entry.Count = +1
            words[tok.Tag][tok.Text] = entry
        }
    }

    var sentences []string

    for _, sentence := range doc.Sentences() {
        sentences = append(sentences, sentence.Text)
    }

    return &NLP {
        Model: doc.Model.Name,
        Words: words,
        Sentences: sentences,
    }, err
}
