package client

//ConfluenceSpace stores the space information
type ConfluenceSpace struct {
    ID   int64  `json:"id,omitempty"`
    Key  string `json:"key,omitempty"`
    Name string `json:"name,omitempty"`
}

//ConfluencePageBodyStorage holds the storage objects of the body
type ConfluencePageBodyStorage struct {
    Value          string `json:"value,omitempty"`
    Representation string `json:"representation,omitempty"`
}

//ConfluencePageBody holds the page contents itself
type ConfluencePageBody struct {
    Storage *ConfluencePageBodyStorage `json:"storage,omitempty"`
}

//ConfluencePageVersion holds the version information for a page
type ConfluencePageVersion struct {
    Number int64 `json:"number,omitempty"`
}

//ConfluencePageAncestor holds the ID of a specific ancestor to a confluence page
type ConfluencePageAncestor struct {
    ID int64 `json:"id,omitempty"`
}

//ConfluencePage stores the base page object
type ConfluencePage struct {
    Title     string                   `json:"title,omitempty"`
    Type      string                   `json:"type,omitempty"`
    ID        string                   `json:"id,omitempty"`
    Ancestors []ConfluencePageAncestor `json:"ancestors,omitempty"`
    Space     *ConfluenceSpace         `json:"space,omitempty"`
    Body      *ConfluencePageBody      `json:"body,omitempty"`
    Version   *ConfluencePageVersion   `json:"version,omitempty"`
}

//ConfluencePageSerachResult stores the base search result object
type ConfluencePages struct {
    Title       string                   `json:"title,omitempty"`
    Excerpt     string                   `json:"excerpt,omitempty"`
    URL         string                   `json:"url,omitempty"`
    LastUpdated string                   `json:"friendlyLastModified,omitempty"`
    Type        string                   `json:"entityType,omitempty"`
    Content     ConfluenceContent        `json:"content,omitempty"`
    Container   ConfluenceContainer          `json:"resultGlobalContainer,omitempty"`
}

//ConfluenceSpace stores the space information
type ConfluenceContainer struct {
    Title       string  `json:"title,omitempty"`
    DisplayURL  string `json:"displayUrl,omitempty"`
}

//ConfluenceContent stores the base search result content object
type ConfluenceContent struct {
    ID        string                   `json:"id,omitempty"`
    Type      string                   `json:"type,omitempty"`
    Status    string                   `json:"status,omitempty"`
    Title     string                   `json:"title,omitempty"`
    Excerpt   string                   `json:"excerpt,omitempty"`
}

//ConfluencePageSearch stores search results for checking existing pages
type ConfluencePageSearch struct {
    Results []ConfluencePage `json:"results,omitempty"`
    Start   int64            `json:"start,omitempty"`
    Limit   int64            `json:"limit,omitempty"`
    Size    int64            `json:"size,omitempty"`
}

//ConfluencePagesSearch stores search results for checking existing pages
type ConfluencePagesSearch struct {
    Results     []ConfluencePages `json:"results,omitempty"`
    Query       string           `json:"cqlQuery,omitempty"`
    Start       int64            `json:"start,omitempty"`
    Limit       int64            `json:"limit,omitempty"`
    Size        int64            `json:"size,omitempty"`
    TotalSize   int64            `json:"totalSize,omitempty"`
    Duration    int64            `json:"searchDuration,omitempty"`
}

//Label the label base object
type Label struct {
}

//ConfluenceConvert is used to store the conversion request or result for a convert command
type ConfluenceConvert struct {
    Value          string `json:"value,omitempty"`
    Representation string `json:"representation,omitempty"`
}

//TinyMceRequest is used for the undocumented TineMCE API
type TinyMceRequest struct {
    EntityID string `json:"entityId,omitempty"`
    SpaceKey string `json:"spaceKey,omitempty"`
    Wiki     string `json:"wiki,omitempty"`
}

func newPage(title, spaceKey string) *ConfluencePage {
    return &ConfluencePage{
        Title: title,
        Type:  "page",
        Space: &ConfluenceSpace{Key: spaceKey},
        Body: &ConfluencePageBody{
            Storage: &ConfluencePageBodyStorage{
                Value:          "",
                Representation: "storage",
            },
        },
    }
}
