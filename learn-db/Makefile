
# Copyright (C) Philip Schlump 2016.
# MIT Licensed
# Source pulled from PureImagination Demo Server

# /home/pschlump/go/src/github.com/Univ-Wyo-Education/S21-4280/learn-db
all: sql mk_doc
	check-json-syntax --ignore-tab-warning ./cfg.json
	check-json-syntax --ignore-tab-warning ./prod-cfg.json
	check-json-syntax --ignore-tab-warning ./db_cfg.json
	rm -f server learn-db
	~/go/src/gitlab.com/pschlump/PureImaginationServer/bin/cmp-local.sh  learn-db

start_watch:
	( cd mtp ; go build )
	./mtp/mtp \
		--comment .../Univ-Wyo-Education/S21-4280/learn-db \
		--watch \
		--cache-burst-override 108 \
		./www-src/index.html \
		&

sql:
	( cd data ; make sql )

mk_code:


# formerly: ValidateTables.go
gen:
	check-json-syntax --ignore-tab-warning cfg.json
	~/go/src/gitlab.com/pschlump/PureImaginationServer/bin/validate_script.sh \
		"t_output" \
		"t_ymux_2fa" \
		"t_ymux_2fa_dev_otk" \
		"t_ymux_2fa_otk" \
		"t_ymux_auth_token" \
		"t_ymux_config" \
		"t_ymux_documents" \
		"t_ymux_origin_ok" \
		"t_ymux_priv" \
		"t_ymux_registration_token" \
		"t_ymux_role" \
		"t_ymux_role_priv" \
		"t_ymux_user" \
		"t_ymux_user_log" \
		"ct_homework" "ct_login" "ct_homework_seen" "ct_homework_grade" "ct_homework_list" "ct_homework_validation" \
			"ct_file_list" "ct_val_homework"

run: all
	./learn-db --hostport ":7001"


deploy0: all
	tar -czf ../blob.tar.gz ./learn-db ./www ./hw ./data ./prod-key
	scp ../blob.tar.gz pschlump@nfc-auth.com:/tmp
	echo
	scp ./prod-cfg.json pschlump@nfc-auth.com:/home/pschlump/learn-db/cfg.json
	scp ./basic-run.sh pschlump@nfc-auth.com:/home/pschlump/learn-db/
	echo
	rsync -r -a -v -e "ssh -l pschlump"    --delete 		./data     			pschlump@nfc-auth.com:/home/pschlump/learn-db/
	rsync -r -a -v -e "ssh -l pschlump"    --delete 		./mt      			pschlump@nfc-auth.com:/home/pschlump/learn-db/
	scp ./db_cfg.json  pschlump@nfc-auth.com:/home/pschlump/learn-db/
	scp ./db_cfg.json  pschlump@nfc-auth.com:/home/pschlump/learn-db/
	-ssh pschlump@nfc-auth.com "mv learn-db/,learn-db learn-db/,learn-db.1 ; mv learn-db/learn-db learn-db/,learn-db"
	scp ./learn-db pschlump@nfc-auth.com:/home/pschlump/learn-db/learn-db
	rsync -r -a -v -e "ssh -l pschlump"    --delete 		./hw      			pschlump@nfc-auth.com:/home/pschlump/learn-db/

deploy:  
	rsync -r -a -v -e "ssh -l pschlump"    --delete 		./www      			pschlump@nfc-auth.com:/home/pschlump/learn-db/
.PHONY: deploy


#deploy_74: mk_doc
#	~/go/src/gitlab.com/pschlump/PureImaginationServer/bin/cmp-prod.sh q8s-service.linux
#	echo ssh pschlump@192.154.97.74 "mkdir -p ~q8s/q8s/server"
#	echo ssh pschlump@192.154.97.74 "mkdir -p ~q8s/q8s/server/www"
#	echo ssh pschlump@192.154.97.74 "mkdir -p ~q8s/q8s/server/Fonts"
#	-echo ssh pschlump@192.154.97.74 "mv ~q8s/q8s/server/q8s-service.linux ,aaaaa"
#	echo scp *.linux pschlump@192.154.97.74:/home/q8s/q8s/server
#	check-json-syntax prod-cfg.json
#	scp prod-cfg.json pschlump@192.154.97.74:/home/q8s/q8s/server
#	scp run-prod.sh pschlump@192.154.97.74:/home/q8s/q8s/server
#	scp prod-cfg.js pschlump@192.154.97.74:/home/q8s/q8s/server
#	scp *.sql pschlump@192.154.97.74:/home/q8s/q8s/server
#	echo scp prod-cfg.js pschlump@192.154.97.74:/home/q8s/q8s/server/www/js/cfg.js
#	rsync -r -a -v -e "ssh -l pschlump"    --delete --exclude www/qr 		./www      			pschlump@192.154.97.74:/home/q8s/q8s/server
#	echo rsync -r -a -v -e "ssh -l pschlump"    --delete 					./Fonts    			pschlump@192.154.97.74:/home/q8s/q8s/server
#	scp *.png pschlump@192.154.97.74:/home/q8s/q8s/server
#	echo "deploy-to-prod: " ${GIT_COMMIT} `date` >>build-log.txt 
#.PHONY: deploy_74


FR=./tmpl/md/

mk_doc: ./www/api-doc.html
.PHONY: mk_doc

./www/api-doc.html: ./tmpl/api-doc.md
	blackfriday-tool ./tmpl/api-doc.md ,api-doc.html
	cat ${FR}/css/pre ${FR}/css/markdown.css ${FR}/css/post ./tmpl/md/md.css ${FR}/css/hpre ,api-doc.html ${FR}/css/hpost >/tmp/api-doc.html
	mv /tmp/api-doc.html ./www/api-doc.html



l_test: l_test_pre l_test1 l_test2 l_test3 l_test4 l_test5 l_test6 l_test7

l_test_pre:
	mkdir -p out
.PHONY: l_test_pre

# URIPath:   "/api/v1/ct_homework_seen",
l_test1:
	wget -o out/ct_homework_seen.err -O out/ct_homework_seen.out "http://127.0.0.1:7001/api/v1/ct_homework_seen"

# URIPath:   "/api/v1/ct_file_list",
l_test2:
	wget -o out/ct_file_list.err -O out/ct_file_list.out "http://127.0.0.1:7001/api/v1/ct_file_list"

# URIPath:   "/api/v1/ct_val_homework",
l_test3:
	wget -o out/ct_val_homework.err -O out/ct_val_homework.out "http://127.0.0.1:7001/api/v1/ct_val_homework"

# URIPath:   "/api/v1/ct_file_list",
l_test4:
	wget -o out/ct_file_list.err -O out/ct_file_list.out "http://127.0.0.1:7001/api/v1/ct_file_list"

# Test where capability
l_test5:
	wget -o out/ct_file_list.2.err -O out/ct_file_list.2.out "http://127.0.0.1:7001/api/v1/ct_file_list?homework_no=26"

# Test where capability
l_test6:
	wget -o out/ct_file_list.3.err -O out/ct_file_list.3.out "http://127.0.0.1:7001/api/v1/ct_file_list?homework_no=13"

# Test the API for getting the list of lessons
l_test7:
	wget -o out/ct_homework.err -O out/ct_homework.out "http://127.0.0.1:7001/api/v1/ct_homework"




#			URIPath:   "/api/v1/ct_file_list",
l_test8:
	wget -o out/ct_file_list.err -O out/ct_file_list.out \
		"http://127.0.0.1:7001/api/v1/ct_file_list?homework_no=26"




xx="http://127.0.0.1:7001/api/v1/run-sql-in-db?homework_id=6e6fa42d-a556-4bc8-61ee-addcae6a6622&stmt=%0A--%20--------------------------------------------------------%20--%20--------------------------------------------------------%0A--%20Note%20the%20%22auth_token%22%20is%20the%20%22ID%22%20for%20this%20row.%20(Primnary%20Key)%0A--%20--------------------------------------------------------%20--%20--------------------------------------------------------%0A%0Adrop%20TABLE%20if%20exists%20%22t_ymux_auth_token%22%20%3B%0ACREATE%20TABLE%20%22t_ymux_auth_token%22%20(%0A%09%20%20%22id%22%09%09%09%09%09uuid%20DEFAULT%20uuid_generate_v4()%20not%20null%20primary%20key%0A%09%2C%20%22user_id%22%09%09%09%09uuid%20not%20null%0A%09%2C%20%22updated%22%20%09%09%09timestamp%0A%09%2C%20%22created%22%20%09%09%09timestamp%20default%20current_timestamp%20not%20null%0A)%3B%0ACOMMENT%20ON%20TABLE%20%22t_ymux_auth_token%22%20IS%20%27version%3A%20b14871d878aa2f9fbaa35f4ad2c6c2baf9584c06%20tag%3A%20v0.1.9%20build_date%3A%20Fri%20Dec%2011%2014%3A30%3A43%20MST%202020%27%3B%0A%0Acreate%20index%20%22t_ymux_auth_token_p1%22%20on%20%22t_ymux_auth_token%22%20(%20%22user_id%22%20)%3B%0Acreate%20index%20%22t_ymux_auth_token_p2%22%20on%20%22t_ymux_auth_token%22%20(%20%22created%22%20)%3B%0A%0A%0AALTER%20TABLE%20%22t_ymux_auth_token%22%0A%09ADD%20CONSTRAINT%20%22t_ymux_auth_token_user_id_fk1%22%0A%09FOREIGN%20KEY%20(%22user_id%22)%0A%09REFERENCES%20%22t_ymux_user%22%20(%22id%22)%0A%3B%0A%0ACREATE%20OR%20REPLACE%20function%20t_ymux_auth_token_upd()%0ARETURNS%20trigger%20AS%20%24%24%0ABEGIN%0A%09--%20version%3A%20b14871d878aa2f9fbaa35f4ad2c6c2baf9584c06%20tag%3A%20v0.1.9%20build_date%3A%20Fri%20Dec%2011%2014%3A30%3A43%20MST%202020%0A%09NEW.updated%20%3A%3D%20current_timestamp%3B%0A%09RETURN%20NEW%3B%0AEND%0A%24%24%20LANGUAGE%20%27plpgsql%27%3B%0A%0A%0ACREATE%20TRIGGER%20t_ymux_auth_token_trig%0ABEFORE%20update%20ON%20%22t_ymux_auth_token%22%0AFOR%20EACH%20ROW%0AEXECUTE%20PROCEDURE%20t_ymux_auth_token_upd()%3B%0A&_ran_=4827085.297543914"

jwt_token=eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRoX3Rva2VuIjoiMzUzMjkxNDUtZWE4MS00MWRhLWFiMGItM2Y5Yzg5NDcxZDNmIiwidXNlcl9pZCI6IjdhOTU1ODIwLTA1MGEtNDA1Yy03ZTMwLTMxMGRhODE1MmI2ZCJ9.vlDkCWvjCFJxoBb2A04_hI4OLfpGP2RH2Mp9qmrbBEr0XIwJrUFRK6O_uMXSXDV_3sQ_Eu0bu0KKeJ9ohtfFqQ

l_test9:
	wget -o out/run-it.err -O out/run-it.out \
		--header="Authorization: bearer ${jwt_token}"  \
		--header='Accept: application/json' \
		${xx}

