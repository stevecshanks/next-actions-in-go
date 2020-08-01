#!/bin/bash

envsubst '$$SERVER_NAME' < /etc/nginx/conf.d/default.conf.template > /etc/nginx/conf.d/default.conf

certbot --nginx -n --email ${EMAIL} --agree-tos --no-eff-email -d ${SERVER_NAME}

cron

# certbot will run nginx in background, so prevent container exiting as soon as it starts
tail -f /dev/null
