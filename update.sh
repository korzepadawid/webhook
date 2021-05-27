#!/bin/bash

path="$1";

cd $path;
git reset --hard;
git pull;

pm2pid="$2";

pm2 stop $pm2pid;
npm install;
npm run build;
pm2 start $pm2pid;
