#!/usr/bin/env bash
git remote add dokku dokku@dokku:goali
ssh -t root@dokku << EOF
    dokku apps:create goali
    # dokku mysql:create goali-db
    # dokku mysql:link goali-db goali
EOF