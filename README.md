# Gator
Gator is a cli tool used to aggregate RSS feeds and display the most recent posts from them.

## prerequisites
To run gator `Postgres 15` or higher needs to be installed. Go version `1.23` or higher.

## Installing
To install the gator cli run the following command:
```shell
go install https://github.com/Flarenzy/blog-aggregator
```
Once installed gator expects an configuration json file located at `~/.gatorconfig.json` where ~ is the home dir of the 
current user. Inside the configuration you need to pass in the psql database link that the app can use.
```json
{"db_url":"postgres://USER:@localhost:5432/gator?sslmode=disable"}
```
The gator DB has to be created before in psql. If on mac use the following command to create the gator db inside of psql:
```sql
CREATE DATABASE gator;
```

## DB migrations
Copy the sql/schemas folder and run goose inside of it to migrate the DB to the latest state.
Goose can be installed with: `go install github.com/pressly/goose/v3/cmd/goose@latest`
```shell
goose postgres "postgres://USER:@localhost:5432/gator" up
```


## Commands
Once installed gator can be run with the following commands:
1. blog-aggregator register USER - registers an user to the app, if an username already exists it will fail. It also sets the current user to the registered user.
2. blog-aggregator login USER - attempts to login as the provided user and fails if no registered user exists.
3. blog-aggregator users - list all registered users and the current logged in user
4. blog-aggregator addfeed "FeedName" "FEED_URL" - adds an RSS feed with the provided name and URL
5. blog-aggregator feeds - list all feeds
6. blog-aggregator follow "FEED_URL" - if the feed url has been added with addfeed follow the feed with the current user
7. blog-aggregator following - list all the feeds an user is following
8. blog-aggregator unfollow "FEED_URL" - unfollows the specified feed
9. blog-aggregator agg - starts aggregating all of the feeds a user is following
10. blog-aggregator browse 5 - returns posts from the feeds an user is following, the default value for the number of posts is 2 and max is 99.

