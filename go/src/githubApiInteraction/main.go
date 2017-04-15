package main

import (
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "time"
  "os"
  "regexp"
  // "net/url"
)

var Committers = make(map[string]int)

type repositoryCommit []struct {
  URL string `json:"url"`
  Sha string `json:"sha"`
  HTMLURL string `json:"html_url"`
  CommentsURL string `json:"comments_url"`
  Commit struct {
    URL string `json:"url"`
    Author struct {
      Name string `json:"name"`
      Email string `json:"email"`
      Date time.Time `json:"date"`
    } `json:"author"`
    Committer struct {
      Name string `json:"name"`
      Email string `json:"email"`
      Date time.Time `json:"date"`
    } `json:"committer"`
    Message string `json:"message"`
    Tree struct {
      URL string `json:"url"`
      Sha string `json:"sha"`
    } `json:"tree"`
    CommentCount int `json:"comment_count"`
    Verification struct {
      Verified bool `json:"verified"`
      Reason string `json:"reason"`
      Signature string `json:"signature"`
      Payload string `json:"payload"`
    } `json:"verification"`
  } `json:"commit"`
  Author struct {
    Login string `json:"login"`
    ID int `json:"id"`
    AvatarURL string `json:"avatar_url"`
    GravatarID string `json:"gravatar_id"`
    URL string `json:"url"`
    HTMLURL string `json:"html_url"`
    FollowersURL string `json:"followers_url"`
    FollowingURL string `json:"following_url"`
    GistsURL string `json:"gists_url"`
    StarredURL string `json:"starred_url"`
    SubscriptionsURL string `json:"subscriptions_url"`
    OrganizationsURL string `json:"organizations_url"`
    ReposURL string `json:"repos_url"`
    EventsURL string `json:"events_url"`
    ReceivedEventsURL string `json:"received_events_url"`
    Type string `json:"type"`
    SiteAdmin bool `json:"site_admin"`
  } `json:"author"`
  Committer struct {
    Login string `json:"login"`
    ID int `json:"id"`
    AvatarURL string `json:"avatar_url"`
    GravatarID string `json:"gravatar_id"`
    URL string `json:"url"`
    HTMLURL string `json:"html_url"`
    FollowersURL string `json:"followers_url"`
    FollowingURL string `json:"following_url"`
    GistsURL string `json:"gists_url"`
    StarredURL string `json:"starred_url"`
    SubscriptionsURL string `json:"subscriptions_url"`
    OrganizationsURL string `json:"organizations_url"`
    ReposURL string `json:"repos_url"`
    EventsURL string `json:"events_url"`
    ReceivedEventsURL string `json:"received_events_url"`
    Type string `json:"type"`
    SiteAdmin bool `json:"site_admin"`
  } `json:"committer"`
  Parents []struct {
    URL string `json:"url"`
    Sha string `json:"sha"`
  } `json:"parents"`
}

func callApi(token, queryString string) (int, error) {
  base := "https://api.github.com/"
  api := "repos/fastlane/fastlane/commits"
  url := fmt.Sprintf("%s%s%s", base, api, queryString)

  //// re := regexp.MustCompile("(?:\\?page=)([\\d]*)(?:&)")
  re := regexp.MustCompile("(?:\\?page=)\\d*")
  initialMatch := re.FindStringSubmatch(queryString)
  fmt.Printf("Current Page: %s\n", initialMatch[0])


  /// me := regexp.MustCompile("(?:\\d*)")
  /// page := me.FindStringSubmatch(initialMatch[0])
  /// fmt.Printf("Current Page: %s\n", page)

  req, err := http.NewRequest("GET", url, nil)
  if err != nil {
    log.Fatal("New request: ", err)
    return 0, err
  }

  client := &http.Client{}

  resp, err := client.Do(req)
  if err != nil {
    log.Fatal("Do: ", err)
    return 0, err
  }
  defer resp.Body.Close()
  var record repositoryCommit

  if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
    log.Println(err)
  }

  fmt.Printf("Records returned on this page: %d\n", len(record))
  if len(record) > 0 {
    for _, commit := range record {
      Committers[commit.Author.Login] += 1
    }
  }
  return len(record), nil
}

func getToken() string {
  token := os.Getenv("GHTOKEN")
  if len(token) > 0 {
    return token
  } else {
    log.Printf("You need to set your Github Access Token")
    log.Printf("export GHTOKEN=<your gh access token>")
    os.Exit(4)
    return "error"
  }
}

func main() {
  token := getToken()
  totalRecords := 0
  additionalRecords := 1
  for i :=1; additionalRecords > 0; i++ {
    queryString := fmt.Sprintf("?page=%d&per_page=100&access_token=%s", i, token)
    if responseCount, err := callApi(token, queryString); err != nil {
      log.Println(err)
    } else {
      additionalRecords = responseCount
      totalRecords += additionalRecords
    }
  }
  fmt.Println(Committers)
  fmt.Printf("Total Records: %d\n", totalRecords)
  fmt.Printf("Commits for aln787: %d\n", Committers["aln787"])
}
