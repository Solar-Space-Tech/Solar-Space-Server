# Solar-Space-Server  
Publish info and trade SSC.  
# Get Started  
Install golang(on 1.16) and setup `GOPATH` following this guide [https://golang.org/doc/install](https://golang.org/doc/install).  
```
$ git clone https://github.com/Solar-Space-Tech/Solar-Space-Server.git  
$ cd Solar-Space-Server  
```
# configure
Create `db.json`, `keystore.json` and `pin&client_secret.json` as `*_example.json` form and same path.  
# Complie
```
$ go build
$ ./Sollar-Space-Server
```

# Use Docker to Deploy (Recommended)
To Use Docker you need to install docker on your system
```
$ docker build -t sst .
$ docker run -itd --name sst -p 443:443 sst
```
