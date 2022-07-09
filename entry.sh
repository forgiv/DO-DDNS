#!/bin/sh

# run dns update once as soon as image start
/do-ddns $DOMAIN $SUBDOMAIN $APIKEY >> /var/log/do-ddns.log

# start cron
/usr/sbin/crond -f -l 8
