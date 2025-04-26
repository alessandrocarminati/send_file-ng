#!/bin/sh
cache_file=min.cache
if git rev-parse --is-inside-work-tree >/dev/null; then
	if ! git show-ref --tags >/dev/null 2>/dev/null; then
		num_commits=$(git rev-list --count HEAD 2>/dev/null)
	else
		first_tag=$(git show-ref --tags 2>/dev/null| tail -n1| cut -d" " -f1)
		last_commit=$(git log --pretty=oneline 2>/dev/null| head -n1| cut -d" " -f1)
		num_commits=$(git rev-list --count $first_tag..$last_commit 2>/dev/null)
	fi
	echo "${num_commits:-0}" | tee ${cache_file} && exit 0
	fi
[ -f ${cache_file} ] && cat ${cache_file}  && exit 0
echo unknown
