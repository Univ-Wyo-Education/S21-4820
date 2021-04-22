#!/bin/bash

XXX=$(pwd)
for i in ?? ; do
	if [ -d $i ] ; then
		echo $i
		cd $i
		make
		pandoc Lect-$i.html -t latex -o Lect-$i.pdf
	fi
	cd $XXX
done

