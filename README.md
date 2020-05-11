# Pull Jira

# How to Run:
```
go run ./main.go <username> <password> <domain> <project>
ex: go run /main.go user password domain.jira.com OPS
```

# Caution:
Proper Error handling not deployed yet. YMMV

## what it does
This program uses Basic Authentication to run an API Jira Search Query and print it to the command line.
It gives options for fields (not all) within a Jira API query.

It runs a JQL Query against the API with GET:
https://developer.atlassian.com/cloud/jira/platform/rest/v3/#api-rest-api-3-search-get


Basic Auth Jira Documentation here:
https://developer.atlassian.com/cloud/jira/platform/basic-auth-for-rest-apis/

