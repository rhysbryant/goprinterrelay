#!/bin/bash
mkdir release 
configFile=""
outputName=""
for file in dist/*
do
 if [[ "$file" == *"win"* ]]; then
	new=${file:0:-4}
	new+=".zip"
	configFile="windows-config-example.json"
	outputName="goprint.exe"
 else
	new=$file
	new+=".tar.gz"
	configFile="linux-config-example.json"
	outputName="goprint"
 fi;
 cp $configFile release/config.json
 cp $file release/$outputName
 archiver make "$new" ui tools release
 rm release/*
done;

mkdir packages
mv dist/*.zip dist/*.gz packages/