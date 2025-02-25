# Circuitos Circuito de Nurburgring

This is a CIRCUITOS website

## Todo content

- Create /bin and /logs directory for logs and binary

## Compile binary with CGO enabled

```
# CGO_ENABLED if not works
sudo apt-get install build-essential

# Compile Linux
env GOOS=linux CGO_ENABLED=1 go build -o app ./cmd/web/*

# For MacOs
env GOOS=darwin GOARCH=amd64 go build -o bin/app ./cmd/api/*
```

## Supervisor Conf

````
cd /etc/supervisor/conf.d/

# Create supervisor file
sudo vim circuitos.conf
````

````
anton@gosha:/etc/supervisor/conf.d$ cat circuitos.conf
[program:circuitos]
command=/home/anton/www/circuitos/app
directory=/home/anton/www/circuitos
autorestart=true
autostart=true
startretries=10
stdout_logfile=/home/anton/www/circuitos/logs/supervisor.log
stderr_logfile=/home/anton/www/circuitos/logs/supervisor.err
````

Then

````
sudo supervisorctl
> status
> update
> status
> restart circuitos
````

## Nginx Configuration

````
server {
        listen 80;

        server_name example.com;

        location / {
                proxy_pass http://localhost:3335;
        }
}
````

## Import From CSV File

````
sqlite> .mode csv
sqlite> .import ./202304.csv transactions
````