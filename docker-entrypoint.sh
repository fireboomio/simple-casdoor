#!/bin/bash
if [ "${MYSQL_ROOT_PASSWORD}" = "" ] ;then MYSQL_ROOT_PASSWORD=285637zq ;fi

service mariadb start

mysqladmin -u root password ${MYSQL_ROOT_PASSWORD}

exec /server --createDatabase=true
