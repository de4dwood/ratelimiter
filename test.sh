#!/bin/bash
count=$1
for i in $(seq $count); do
 curl -XGET "http://localhost:9000/request?url=https://www.alibaba.ir/$2&user=test2" -I
done
