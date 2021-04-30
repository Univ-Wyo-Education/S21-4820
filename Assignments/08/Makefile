
DIR=../../../F20-1015
PP=`pwd`
PY= 
IMG=

MD= Assignment-08.md
HTML= Assignment-08.html

all: pre $(MD) $(HTML) 

pre:
	cp /dev/null open.1
.PHONEY: pre

%.md: %.raw.md $(PY) $(IMG)
	m4 -P $< >$@

%.html: %.md
	blackfriday-tool ./$< $@
	echo cat ./${DIR}/md.css $@ >/tmp/$@
	cat ./${DIR}/css/pre ./${DIR}/css/markdown.css ./${DIR}/css/post ./${DIR}/md.css ./${DIR}/css/hpre $@ ./${DIR}/css/hpost >/tmp/$@
	mv /tmp/$@ ./$@
	echo "file://${PP}/$@" >>open.1

