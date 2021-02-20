package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"
)

//ConfluenceConfig holds the current client configuration
type ConfluenceConfig struct {
	Username string
	Password string
	URL      string
	Debug    bool
}


//ConfluenceClient is the primary client to the Confluence API
type ConfluenceClient struct {
	username string
	password string
	baseURL  string
	debug    bool
	client   *http.Client
}

//Client returns a new instance of the client
func Client(config *ConfluenceConfig) *ConfluenceClient {
	return &ConfluenceClient{
		username: config.Username,
		password: config.Password,
		baseURL:  config.URL,
		debug:    config.Debug,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (c *ConfluenceClient) doRequest(method, url string, content, responseContainer interface{}) []byte {
	b := new(bytes.Buffer)
	if content != nil {
		json.NewEncoder(b).Encode(content)
	}
	furl := c.baseURL + url
	if c.debug {
		log.Println("Full URL", furl)
		log.Println("JSON Content:", b.String())
	}
	request, err := http.NewRequest(method, furl, b)
	request.SetBasicAuth(c.username, c.password)
	request.Header.Add("Content-Type", "application/json; charset=utf-8")
	if err != nil {
		log.Fatal(err)
	}
	if c.debug {
		log.Println("Sending request to services...")
	}
	response, err := c.client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	if c.debug {
		log.Println("Response received, processing response...")
		log.Println("Response status code is", response.StatusCode)
		log.Println(response.Status)
	}
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	if c.debug {
		log.Println("Response from service...", string(contents))
	}
	if response.StatusCode != 200 {
		log.Fatal("Bad response code received from server: ", response.Status)
	}
	json.Unmarshal(contents, responseContainer)
	return contents
}

func (c *ConfluenceClient) uploadFile(method, url, content, filename string, responseContainer interface{}) []byte {
	b := new(bytes.Buffer)
	writer := multipart.NewWriter(b)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		log.Fatal(err)
	}
	part.Write([]byte(content))
	writer.WriteField("minorEdit", "true")
	//writer.WriteField("comment", "test")
	writer.Close()

	furl := c.baseURL + url
	if c.debug {
		log.Println("Full URL", furl)
	}
	request, err := http.NewRequest(method, furl, b)
	request.SetBasicAuth(c.username, c.password)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	request.Header.Add("X-Atlassian-Token", "nocheck")
	if err != nil {
		log.Fatal(err)
	}
	if c.debug {
		log.Println("Sending request to services...")
	}
	response, err := c.client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	if c.debug {
		log.Println("Response received, processing response...")
		log.Println("Response status code is", response.StatusCode)
	}
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode != 200 {
		log.Fatal("Bad response code received from server: ", response.Status)
	}
	json.Unmarshal(contents, responseContainer)
	return contents
}

//AddPage adds a new page to the space with the given title
func (c *ConfluenceClient) AddPage(title, spaceKey, body string, ancestor int64) {
	page := newPage(title, spaceKey)
	if ancestor > 0 {
		page.Ancestors = []ConfluencePageAncestor{
			ConfluencePageAncestor{ancestor},
		}
	}
	response := &ConfluencePage{}
	page.Body.Storage.Value = body
	//page.Body.Storage.Representation = "wiki"
	c.doRequest("POST", "/rest/api/content/", page, response)
	log.Println("Confluence page added with ID", response.ID, "and version", response.Version.Number)
}

//UpdatePage adds a new page to the space with the given title
func (c *ConfluenceClient) UpdatePage(title, spaceKey, body string, ID string, version, ancestor int64) {
	page := newPage(title, spaceKey)
	page.ID = ID
	page.Version = &ConfluencePageVersion{version}
	if ancestor > 0 {
		page.Ancestors = []ConfluencePageAncestor{
			ConfluencePageAncestor{ancestor},
		}
	}
	response := &ConfluencePage{}
	page.Body.Storage.Value = body
	//page.Body.Storage.Representation = "wiki"
	c.doRequest("PUT", "/rest/api/content/"+ID, page, response)
	log.Println("Confluence page updated with ID", response.ID, "and version", response.Version.Number)
}

//SearchPages searches for pages in the space that meet the specified criteria
func (c *ConfluenceClient) SearchPages(title, spaceKey string) (results *ConfluencePageSearch) {
	results = &ConfluencePageSearch{}
	c.doRequest("GET", "/rest/api/content?title="+url.QueryEscape(title)+"&spaceKey="+url.QueryEscape(spaceKey)+"&expand=version", nil, results)
	return results
}

//CQLSearchPages searches for pages in the space that meet the specified criteria
func (c *ConfluenceClient) CQLSearchPages(title string) (results *ConfluencePagesSearch) {
	results = &ConfluencePagesSearch{}
	c.doRequest("GET", "/rest/api/search?cql="+url.QueryEscape("(type=page and text~\"" +title+"\")")+"&expand=version", nil, results)
	return results
}

//CQLSearchPagesBy searches for pages in the space that meet the specified criteria
func (c *ConfluenceClient) CQLSearchPagesBy(cql string) (results *ConfluencePagesSearch) {
	results = &ConfluencePagesSearch{}
	c.doRequest("GET", "/rest/api/search?limit=5&cql="+url.QueryEscape(cql)+"", nil, results)
	return results
}

//AddAttachment adds an attachment to an existing page
func (c *ConfluenceClient) AddAttachment(content, pageID, filename string) {
	results := &ConfluencePageSearch{}
	c.uploadFile("PUT", "/rest/api/content/"+pageID+"/child/attachment", content, filename, &results)
}

//GetLabel searches for pages in the space that meet the specified criteria
func (c *ConfluenceClient) GetLabel(label, spaceKey string) (results *Label) {
	c.doRequest("GET", "/rest/api/label?name="+label, nil, &Label{})
	return results
}
