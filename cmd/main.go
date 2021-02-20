package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/slack-go/slack"

    "github.com/fr123k/confluence-slackbot/pkg/config"

	confluence "github.com/fr123k/confluence-slackbot/pkg/confluence-cli"
	// natural language processing
	prose "github.com/jdkato/prose/v2"
)

func confluenceCli(cfg *config.Config) (*confluence.ConfluenceClient) {
    var confluenceCfg = confluence.ConfluenceConfig{}
    confluenceCfg.URL = cfg.ConfigConfluence.URL
    confluenceCfg.Username = cfg.ConfigConfluence.Username
    confluenceCfg.Password = cfg.ConfigConfluence.Token
    confluenceCfg.Debug = cfg.Debug
    return confluence.Client(&confluenceCfg)
}

func main() {

    cfg, err := config.Configuration()
    if err != nil {
        fmt.Printf("Could not read configuration: %v", err)
        os.Exit(1)
    }

    api := slack.New(
        cfg.ConfigSlack.Token,
        slack.OptionDebug(cfg.Debug),
        slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
    )

    rtm := api.NewRTM()
    go rtm.ManageConnection()

    http.HandleFunc(cfg.Server.ActionURL, actionHandler)
    go http.ListenAndServe(fmt.Sprintf(":%d",cfg.Server.Port), nil)

    for msg := range rtm.IncomingEvents {
        fmt.Print("Event Received: ")
        switch ev := msg.Data.(type) {
        case *slack.HelloEvent:
            // Ignore hello

        case *slack.ConnectedEvent:
            fmt.Println("Infos:", ev.Info)
            fmt.Println("Connection counter:", ev.ConnectionCount)
            rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", ""))

        case *slack.MessageEvent:
            msg := ev.Msg
            if msg.SubType != "" {
                break // We're only handling normal messages.
            }

            if msg.ThreadTimestamp != "" {
                break
            }

            fmt.Printf("Message: %s\n", msg.Text)

            // Create a new document with the default configuration:
            doc, err := prose.NewDocument(msg.Text)
            if err != nil {
                log.Fatal(err)
            }
            // var entities = doc.Entities()
            var words []string
            var query []string

            for _, tok := range doc.Tokens() {
                // fmt.Println(tok.Text, tok.Tag, tok.Label)
                if ((tok.Tag == "NNP" || tok.Tag == "NNS" || tok.Tag == "NNPS" || tok.Tag == "NN" )&& len(tok.Text) > 2) {
                    query = append(query, fmt.Sprintf(" title~\"%s\" ", tok.Text))
                    words = append(words, tok.Text)
                }
            }
            // Iterate over the doc's named-entities:
            // for _, ent := range entities {
            //     fmt.Println(ent.Text, ent.Label)
            // }

            var text = strings.Join(query," or ")
            fmt.Println(words)
            var result *confluence.ConfluencePagesSearch = nil
            if len(words) <= 0 {
                result = &confluence.ConfluencePagesSearch {
                    Size: 0,
                }
            } else {
                confluenceCli := confluenceCli(cfg)
                result = confluenceCli.CQLSearchPagesBy(fmt.Sprintf("label = \"kb-how-to-article\" and (type=page and %s)", text))
            }

            if result.Size > 0 {
                fmt.Println("Page Found!")
                fmt.Println()
            } else {
                // value is passed to message handler when request is approved.
                attachment := slack.Attachment{
                    Text:       "No page with matching title found. Proceed with full text search ?",
                    Color:      "#f9a41b",
                    CallbackID: "search",
                    Actions: []slack.AttachmentAction{
                        {
                            Name: "search",
                            Text: "Yes",
                            Type: "button",
                            Value: strings.Join(words,","),
                        },
                        {
                            Name:  "cancel",
                            Text:  "No",
                            Type:  "button",
                            Style: "danger",
                            Value: "cancel",
                        },
                    },
                }

                message := slack.MsgOptionAttachments(attachment)

                replyMsg := slack.MsgOptionCompose(threadOptionMessage(msg), message)

                channelID, timestamp, err := api.PostMessage(msg.Channel, slack.MsgOptionText("", false), replyMsg)
                fmt.Printf("Message with buttons sucessfully sent to channel %s at %s", channelID, timestamp)
                if err != nil {
                    fmt.Printf("Could not send message: %v", err)
                }
            }

            replyMsg := slack.MsgOptionCompose(threadOptionMessage(msg), searchResult(result, cfg))

            // Create a response object.
            channelID, timestamp, err := rtm.PostMessage(msg.Channel, replyMsg)
            fmt.Printf("Message with buttons sucessfully sent to channel %s at %s", channelID, timestamp)
            if err != nil {
                fmt.Printf("Could not send message: %v", err)
            }

        case *slack.PresenceChangeEvent:
            fmt.Printf("Presence Change: %v\n", ev)

        case *slack.LatencyReport:
            fmt.Printf("Current latency: %v\n", ev.Value)

        case *slack.DesktopNotificationEvent:
            fmt.Printf("Desktop Notification: %v\n", ev)

        case *slack.RTMError:
            fmt.Printf("Error: %s\n", ev.Error())

        case *slack.InvalidAuthEvent:
            fmt.Printf("Invalid credentials")

        default:

            // Ignore other events..
            // fmt.Printf("Unexpected: %v\n", msg.Data)
        }
    }
}

func searchResult(result *confluence.ConfluencePagesSearch, config *config.Config) (slack.MsgOption) {
    var blocks []slack.Block
    
    headerText := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("We found *%d Pages* in Confluence (in %d ms) with the following query `%s`.", result.TotalSize, result.Duration, result.Query), false, false)

    var overflowBlocks []*slack.OptionBlockObject

    for i := range result.Results {
        // Build Text Objects associated with each option
        overflowOptionTextOne := slack.NewTextBlockObject("plain_text", fmt.Sprintf("Option %d", i), false, false)
        overflowOptionOne := slack.NewOptionBlockObject(fmt.Sprintf("value-%d", i), overflowOptionTextOne, nil)
        overflowBlocks = append(overflowBlocks, overflowOptionOne)
    }

    // Build overflow section
    overflow := slack.NewOverflowBlockElement("", overflowBlocks...)

    // Create the header section
    headerSection := slack.NewSectionBlock(headerText, nil, slack.NewAccessory(overflow))
    divSection := slack.NewDividerBlock()

    blocks = append(blocks, headerSection)
    blocks = append(blocks, divSection)


    // Shared Objects
    locationPinImage := slack.NewImageBlockElement("https://api.slack.com/img/blocks/bkb_template_images/tripAgentLocationMarker.png", "Location Pin Icon")
    
    for _, element := range result.Results  {
        var url = fmt.Sprintf("%s%s",config.ConfigConfluence.URL,element.URL)

        //TODO render confluence labels into the search result
        hotelOneInfo := slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*<%s|%s>*\n%s\n", url, element.Content.Title, element.LastUpdated), false, false)
        hotelOneLoc := slack.NewTextBlockObject("plain_text", fmt.Sprintf("Location: %s", element.Container.Title), true, false)

        hotelOneSection := slack.NewSectionBlock(hotelOneInfo, nil, nil)
        hotelOneContext := slack.NewContextBlock("", []slack.MixedElement{locationPinImage, hotelOneLoc}...)

        blocks = append(blocks, hotelOneSection)
        blocks = append(blocks, hotelOneContext)
        blocks = append(blocks, divSection)
    }

    //TODO implements browsing the search results
    // Action button
    // btnTxt := slack.NewTextBlockObject("plain_text", "Next 2 Results", false, false)
    // nextBtn := slack.NewButtonBlockElement("next", "click_me_123", btnTxt)
    // actionBlock := slack.NewActionBlock("", nextBtn)
    // blocks = append(blocks, actionBlock)

    // Build Message with blocks created above
    msg := slack.MsgOptionBlocks(blocks...)
    return msg
}

func actionHandler(w http.ResponseWriter, r *http.Request) {
    //TODO validate request signature to identify the request is coming from slack
    api := slack.New(
        "xoxb-1764752828277-1768305869763-bxj7eCV8LRpYrdyFC1eVU3ZX",
        slack.OptionDebug(false),
        slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
    )
    var payload slack.InteractionCallback
    err := json.Unmarshal([]byte(r.FormValue("payload")), &payload)
    fmt.Printf("response JSON: %s", r.FormValue("payload"))
    if err != nil {
        fmt.Printf("Could not parse action response JSON: %v", err)
    }
    fmt.Printf("Message button pressed by user %s with value %s", payload.User.Name, payload.ActionCallback.AttachmentActions[0].Value)

    text := payload.ActionCallback.AttachmentActions[0].Value
    var queries []string
    words := strings.Split(text, ",")
    for _, word := range words {
        queries = append(queries, fmt.Sprintf(" text~\"%s\" ", word))
    }
    var query = strings.Join(queries," or ")
    fmt.Println(query)

    cfg, err := config.Configuration()
    if err != nil {
        fmt.Printf("Could not read configuration: %v", err)
        os.Exit(1)
    }

    confluenceCli := confluenceCli(cfg)

    result := confluenceCli.CQLSearchPagesBy(fmt.Sprintf("(type=page and (%s))", query))

    //TODO add handling of case no page found.
    fmt.Printf("Found %d with %s\n", result.Size, result.Query)

    message := slack.MsgOptionReplaceOriginal(payload.ResponseURL)
    replyMsg := slack.MsgOptionCompose(threadOptionPayload(payload), message, searchResult(result, cfg))

    channelID, timestamp, err := api.PostMessage(payload.Channel.ID, slack.MsgOptionText("", false), replyMsg)
                fmt.Printf("Message with buttons sucessfully sent to channel %s at %s", channelID, timestamp)
                if err != nil {
                    fmt.Printf("Could not send message: %v", err)
                }

}

func threadOptionMessage(msg slack.Msg) (slack.MsgOption) {
    // Respond in same thread if message came from a thread.
    if msg.ThreadTimestamp != "" {
        return slack.MsgOptionTS(msg.ThreadTimestamp)
    }

    if !strings.HasPrefix(msg.Channel, "D") {
        return slack.MsgOptionTS(msg.Timestamp)
    }
    return slack.MsgOptionTS(msg.Timestamp)
}

func threadOptionPayload(payload slack.InteractionCallback) (slack.MsgOption) {
    // Respond in same thread if message came from a thread.
    if payload.OriginalMessage.ThreadTimestamp != "" {
        return slack.MsgOptionTS(payload.OriginalMessage.ThreadTimestamp)
    }

    if !strings.HasPrefix(payload.OriginalMessage.Channel, "D") {
        return slack.MsgOptionTS(payload.MessageTs)
    }
    return slack.MsgOptionTS(payload.MessageTs)
}
