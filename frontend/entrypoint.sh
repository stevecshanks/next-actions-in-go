#!/bin/bash

envsubst '$$SERVER_NAME' < /etc/nginx/conf.d/default.conf.template > /etc/nginx/conf.d/default.conf

certbot --nginx -n --email ${EMAIL} --agree-tos --no-eff-email -d ${SERVER_NAME}

cron

# Certbot will run Nginx in background, so we need to prevent container exiting as soon as it starts. We also cannot
# have this log file as a symbolic link to /dev/stdout (or /proc/1/fd/1) as Certbot will crash with an "Illegal Seek"
# error when it tries to determine if the log file requires rotation. So, tail the log file instead.
tail -F /var/log/letsencrypt/letsencrypt.log
