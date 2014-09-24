#!/bin/bash

cd `dirname $0`

if [ -f /home/paas/paas/env.sh ]; then
  source /home/paas/paas/env.sh
fi

./grepprocessinfo

