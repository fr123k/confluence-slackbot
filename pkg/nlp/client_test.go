package nlp

import (
    "fmt"
    "strings"
    "testing"
)

const exampleText1 = `One question: what would be required if we want to expose a Docker container for testing purposes in the VPN?
Could we use the phdp-dev cluster in a dedicated namespace with an additional DNS?
Or would it be easier for you to put that on a separate VM - still with an additional DNS?
Context: for Smart4Health we are using a test-broker to ingest data according to the FHIR standard.
This component exposes an HTTPS endpoint and must speak to the PHDP.`

func assert(t *testing.T, err error) {
    if err != nil {
        t.Error(err.Error())
    }
}

func assertError(t *testing.T, err error, want string) {
    if err == nil {
        t.Errorf("Expected error was nil, got: nil, want: %s", want)
    }
    if !strings.Contains(err.Error(), want) {
        t.Errorf("The error message is wrong, got: %s, want: %s.", err.Error(), want)
    }
}

func TestParseNoText(t *testing.T) {
    _, err := Parse("")
    assert(t, err)
}

// TestParseSimpleText1 parse a simple text and check the parsed words.
func TestParseSimpleText1(t *testing.T) {
    nlp, err := Parse(exampleText1)
    assert(t, err)

    if nlp.Model != "en-v2.0.0" {
        t.Errorf("The expected model name is wrong, got: %s, want: %s.", nlp.Model, "en-v2.0.0")
    }

    if len(nlp.Sentences) != 5 {
        fmt.Printf("Sentences: %v\n", nlp.Sentences)

        t.Errorf("The found number of sentences is wrong, got: %d, want: %d.", len(nlp.Sentences), 5)
    }

    if len(nlp.Words) != 22 {
        for k, v := range nlp.Words {
            fmt.Printf("%v:%v\n", k, v)
        }

        t.Errorf("The found number of words is wrong, got: %d, want: %d.", len(nlp.Words), 22)
    }

    if len(nlp.Nouns().Entries) != 20 {
        fmt.Printf("Nouns: %v\n", nlp.Nouns())
        t.Errorf("The found number of nouns is wrong, got: %d, want: %d.", len(nlp.Nouns().Entries), 20)
    }

    if len(nlp.Verbs().Entries) != 16 {
        fmt.Printf("Verbs: %v\n", nlp.Verbs())
        t.Errorf("The found number of verbs is wrong, got: %d, want: %d.", len(nlp.Verbs().Entries), 16)
    }

    verbs := nlp.Verbs()
    summary := verbs.ForEachWithSort(func (entry Entry) (string, string) {
        return fmt.Sprintf("%s %d", entry.Value, entry.Count), ","
        
    }, func (i, j int) bool {
        return verbs.Entries[i].Value < verbs.Entries[j].Value
    })

    if summary != "according 1,according 1,are 1,be 1,expose 1,exposes 1,ingest 1,put 1,required 1,speak 1,testing 1,testing 1,use 1,using 1,using 1,want 1" {
        t.Errorf("The summary text is wrong, got: '%s', want: '%s'.", summary, "according 1,according 1,are 1,be 1,expose 1,exposes 1,ingest 1,put 1,required 1,speak 1,testing 1,testing 1,use 1,using 1,using 1,want 1")
    }
}
