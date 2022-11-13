![GitHub Workflow Status](https://img.shields.io/github/workflow/status/sanyarise/hezzlTest/Go)
![GitHub top language](https://img.shields.io/github/languages/top/sanyarise/hezzlTest)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/sanyarise/hezzlTest)
![GitHub repo file count (file type)](https://img.shields.io/github/directory-file-count/sanyarise/hezzlTest)
![GitHub commit activity](https://img.shields.io/github/commit-activity/m/sanyarise/hezzlTest)
![GitHub contributors](https://img.shields.io/github/contributors/sanyarise/hezzlTest)
![GitHub last commit](https://img.shields.io/github/last-commit/sanyarise/hezzlTest)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/sanyarise/hezzlTest)

<img align="right" width="50%" src="./images/image.jpeg">

# hezzlTest
 
## Task description

1. Describe the proto file with the service from 3 methods: add user, delete user, list of users

2. Implement a gRPC service based on a proto file on Go

3. Use PostgreSQL for data storage

4. upon request for a list of users, the data will be cached in redis for a minute and taken from radish

5. When adding a user, make a log in ClickHouse

6. Add logs to ClickHouse via Kafka queue

## Solution notes

- :book: standard Go project layout (well, more or less :blush:)
- :cd: github CI/CD + docker compose + Makefile included
- :white_check_mark: tests with mocks included

## HOWTO

- run service with tests and client `make run`
- in client type kind of operation (create (then the number of users to create), delete (to delete user type id), all (show all saving users), ch (show logs from clickhouse), exit)
- run service without client `make up`
- run tests `make test`


