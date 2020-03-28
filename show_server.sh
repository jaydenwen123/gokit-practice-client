#!/bin/bash
client="myClient"
pids=`ps -ef |grep ./$client |grep -v "grep"   |awk '{print $2}'`
echo "the pid list:" $pids
