package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type Expand struct {
	Expand     string `json:"expand"`
	MaxResults int    `json:maxResults`
	Total      int    `json:total`
	//Issues     []Issue
}
type Issues struct {
	Issues []Issue
}

type Issue struct {
	Id     string `json:id`
	Key    string `json:key`
	Self   string `json:self`
	Fields Fields
	Names  Names
}

type Fields struct {
	Summary  string `json:summary`
	Priority Priority
	Assignee Assignee
	Created  SpecialDate `json:created`
	Updated  SpecialDate `json:updated`
	Status   Status
}

type Status struct {
	Name string `json:name`
}

type SpecialDate struct {
	time.Time
}

type Priority struct {
	Self string `json:self`
	Name string `json:name`
}

type Assignee struct {
	Self string `json:self`
	Name string `json:name`
}

type Names struct {
	AdditionalProperties string `json:Additional Properties`
}

func (sd *SpecialDate) UnmarshalJSON(input []byte) error {
	strInput := string(input)
	strInput = strings.Trim(strInput, `"`)
	newTime, err := time.Parse("2006-01-02T15:04:05.999999999Z0700", strInput)
	if err != nil {
		return err
	}

	sd.Time = newTime
	return nil
}

/*
References
http://polyglot.ninja/golang-making-http-requests/
https://stackoverflow.com/questions/19253469/make-a-url-encoded-post-request-using-http-newrequest
https://appdividend.com/2019/12/02/golang-http-example-get-post-http-requests-in-golang/
https://girishjoshi.io/post/implementing-http-basic-authentication-in-golang/
https://www.sohamkamani.com/blog/2017/10/18/parsing-json-in-golang/

Jira:
https://developer.atlassian.com/cloud/jira/platform/rest/v3/#api-rest-api-3-filter-defaultShareScope-put
https://godoc.org/github.com/andygrunwald/go-jira#FilterSearchOptions
https://godoc.org/github.com/andygrunwald/go-jira#pkg-examples

https://developer.atlassian.com/cloud/jira/platform/rest/v3/#api-rest-54: Search for issues using JQL (GET)

https://tanaikech.github.io/2017/09/15/spreadsheets.values.batchupdate-using-golang/

https://developer.atlassian.com/cloud/jira/platform/rest/v3/#api-rest-api-3-search-get


*/

func queryJira(username string, password string, query []string, dataToAdd []string, apiURL string) []byte {

	username = username
	passwd := password

	//function to call Jira, allow different queries
	//apiURL := "https://jira.points.com/rest/api/2/search"
	resource := "/rest/api/2/search"
	data := url.Values{}
	data.Set(query[0], query[1])

	for _, value := range dataToAdd {
		data.Add("fields", value)
	}

	//building URL
	u, _ := url.ParseRequestURI(apiURL)
	u.Path = resource
	u.RawQuery = data.Encode()
	urlStr := fmt.Sprintf("%v", u)
	fmt.Println("URL TO COPY PASTE FOR CHECKING IN BROWSER\n", urlStr, "\n\n")

	client2 := &http.Client{}
	r, _ := http.NewRequest("GET", urlStr, nil)
	//setting basicAuth credentials
	r.SetBasicAuth(username, passwd)
	//setting head for JSON
	r.Header.Add("Accept", "application/json")
	resp, err := client2.Do(r)
	if err != nil {
		fmt.Println("error connecting: ", err)
	}
	defer resp.Body.Close()

	//reading body response
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading data after being sent", err)
	}

	return bodyText

}

func remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	// We do not need to put s[i] at the end, as it will be discarded anyway
	return s[:len(s)-1]
}

func Find(a []string, x string) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}
	return len(a)
}

func main() {

	// == THIS IS FOR ENTERING CREDENTIALS ON THE COMMAND LINE
	if len(os.Args) != 6 {
		log.Fatalln("Usage: main <username> <password> <domain> <ticket status>")
	}
	var username = os.Args[1]

	username = "haydn.johnson@points.com" // testing truffle

	var password = os.Args[2]
	var ApiURL = "https://" + os.Args[3] + "/rest/api/latest/search"
	var jiraProject = os.Args[4]
	var ticketStatus = os.Args[5]
	var dataToAdd []string

	reader := bufio.NewReader(os.Stdin)

	fieldsToChoose := []string{"created", "updated", "summary", "status", "priority", "assignee"}
	//Foor Loop for options.
	for x := 0; x < 2; {
		fmt.Println("Please Enter Field # to print out to screen")
		for key, value := range fieldsToChoose {
			fmt.Print("[", key, "]", value, "\n")
		}
		//fmt.Println(fieldsToChoose)
		value, _ := reader.ReadString('\n')
		value = strings.TrimSpace(value)
		if value == "" {
			break
		}

		valueInt, err := strconv.Atoi(value)
		if err != nil {
			panic(err)
		}
		//Find(fieldsToChoose, value)
		dataToAdd = append(dataToAdd, fieldsToChoose[valueInt])

		fieldsToChoose = remove(fieldsToChoose, valueInt)

	}

	query := []string{"jql", "project=" + jiraProject, "and status = " + ticketStatus, "  AND  issuetype not in subtaskIssueTypes()"}
	fmt.Println("Will run this query in JIRA:: ", query[1])
	jiraResponse := queryJira(username, password, query, dataToAdd, ApiURL)

	var jiraIssues Issues
	json.Unmarshal([]byte(jiraResponse), &jiraIssues)
	//fmt.Println(issues)
	for _, issue := range jiraIssues.Issues {

		fmt.Print(issue.Key, " ")
		for _, value := range dataToAdd {
			//switch statements for adding selected fields
			// KEY value is printed to screen by defaultefault
			switch value {
			case "created":
				formatCreated := fmt.Sprintf("%d-%02d-%02d", issue.Fields.Created.Year(), issue.Fields.Created.Month(), issue.Fields.Created.Day())
				fmt.Print(formatCreated, " ")
			case "updated":
				formatUpdated := fmt.Sprintf("%d-%02d-%02d", issue.Fields.Updated.Year(), issue.Fields.Updated.Month(), issue.Fields.Updated.Day())
				fmt.Print(formatUpdated, " ")
			case "summary":
				fmt.Print(issue.Fields.Summary, " ")
			case "status":
				fmt.Print(issue.Fields.Status.Name, " ")
			case "priority":
				fmt.Print(issue.Fields.Priority.Name, " ")
			case "assignee":
				fmt.Print(issue.Fields.Assignee.Name, " ")
			}

		}
		fmt.Println()

	}
}
