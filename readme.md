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
