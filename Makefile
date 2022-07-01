include .env

build:
	go build .

run:
	./do-ddns $(DOMAIN) $(SUBDOMAIN) $(APIKEY) 
