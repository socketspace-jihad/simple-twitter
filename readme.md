# Simple Twitter
This is a minimalist social ecosystem that focuses on the core of digital interaction: the message. Designed for speed and simplicity, it allows users to broadcast thoughts, engage with a global feed, and maintain full control over their digital presence through a robust CRUD-based architecture.
- User Authentication
- Create a tweet
- Read a tweet
- Update a tweet
- Delete a tweet

# Pre-Requisite
```
  sudo apt upgrade
  sudo apt install git golang-go postgresql
```
# Setup postgresql
```
  sudo su
  sudo -u postgres psql postgres
  ALTER USER postgres WITH PASSWORD '<your new password>';
```
dont close the psql terminal yet
copy the value from internal/db/postgresql/migration/ddl.sql
paste it on your psql session

# How to Build the app
`go build -o app .`

# How to run the app
set the env, look at the .env.example env keys
```
./app

```
