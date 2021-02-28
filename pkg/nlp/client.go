package nlp

import (
    // natural language processing
    prose "github.com/jdkato/prose/v2"
)

type NLP struct {
    tokens map[string][]string
    entities map[string][]string
}
func parse(text string) (*NLP, error) {
    doc, err := prose.NewDocument(text)
    if err != nil {
        return nil, err
    }
    var words []string
    // var query []string

    for _, tok := range doc.Tokens() {
        // fmt.Println(tok.Text, tok.Tag, tok.Label)
        if ((tok.Tag == "NNP" || tok.Tag == "NNS" || tok.Tag == "NNPS" || tok.Tag == "NN" )&& len(tok.Text) > 2) {
            // query = append(query, fmt.Sprintf(" title~\"%s\" ", tok.Text))
            words = append(words, tok.Text)
        }
    }
    // Iterate over the doc's named-entities:
    // for _, ent := range doc.Entities() {
    //     fmt.Println(ent.Text, ent.Label)
    // }
	return nil, err
}
