#!/bin/bash
for file in dist/*
do

 if [ "$file" == *"win"* ]; then
	new=${file:0:-4}
	new+=".zip"
	echo "win"
 else
	new=$file
	new+=".tar.gz"
 fi;
 echo "$new $file"
 archiver make "$new" $file ui config.json
done;

mkdir packages
mv *.zip *.gz packages/