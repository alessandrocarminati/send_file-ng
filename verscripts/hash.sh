#!/bin/sh
cache_file=hash.cache
if git rev-parse --is-inside-work-tree >/dev/null; then
	git log --pretty=oneline| head -n1 |cut -d" " -f1| tee ${cache_file} && exit 0
fi
[ -f ${cache_file} ] && cat ${cache_file}  && exit 0
echo unknown
