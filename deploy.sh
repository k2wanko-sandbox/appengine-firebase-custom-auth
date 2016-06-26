#!/bin/sh

cd $(dirname ${0})/backend

application=${APPLICATION:=k2wanko-sandbox-firebase-auth}
branch=${BRANCH:=$(git rev-parse --abbrev-ref HEAD)}
token=$(gcloud auth print-access-token 2> /dev/null)

appcfg.py update --oauth2_access_token $token --application $application --version $branch .
