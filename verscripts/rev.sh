#!/bin/sh
cache_file=rev.cache
state="unknown"
if git rev-parse --is-inside-work-tree >/dev/null; then
	state="clean"
	git status --untracked-files=no --porcelain 2>/dev/null|  grep -q M && hash=$(git diff HEAD -- . ':(exclude)*.cache' 2>/dev/null| md5sum | cut -d ' ' -f1) && state="rev-$(echo "$hash" | cut -c1-5)"
	echo $state | tee ${cache_file} && exit 0
fi
[ -f ${cache_file} ] && cat ${cache_file}  && exit 0
echo $state
