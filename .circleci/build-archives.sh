#!/bin/bash
mkdir release 
configFile=""
outputName=""
cp -R ui tools  release/
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
 cp $configFile release/config.json -f
 cp $file release/$outputName -f
 cd release/
 archiver make "../$new" ui tools config.json goprint*
 cd ..
done;

mkdir packages
mv dist/*.zip dist/*.gz packages/