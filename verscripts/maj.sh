#!/bin/sh
cache_file=maj.cache
if git rev-parse --is-inside-work-tree > /dev/null; then
  maj=$(git show-ref --tags | tail -n1 | cut -d" " -f2 | cut -d/ -f3)
  echo ${maj:-0} | tee ${cache_file} && exit 0
else
 [ -f ${cache_file} ] && cat ${cache_file}  && exit 0
 echo "unknown"
fi


