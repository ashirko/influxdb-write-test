#!/bin/bash
IFS=$'\n'
COMMANDS=`cat scripts/commands.txt`
OPTIONS=`cat scripts/options.txt`
LOGFILE=log/result.log
for command in ${COMMANDS}; do
	for option in ${OPTIONS}; do
		${command} -c ${option} -m ${option}  >> ${LOGFILE} 2>&1
	done
done