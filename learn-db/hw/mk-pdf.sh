#!/bin/bash

for i in hw*.html ; do
	# strip_tags_validate
	sed -e '/^<h4>Tags/d' -e '/^<h4>Validate/d' $i >tmp.html
	BA=$( basename $i .html )
	pandoc tmp.html -t latex -o $BA.pdf
done

