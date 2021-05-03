
m4_include(../../Lect/lect-setup.m4)

# Use pg_dump to backup a PostgreSQL database

Pts: 50
Due May 6th.


Using you PostgreSQL system dump out a database
with at least 1 table in it.


Usually you will need to be the "postgre" user to do this.
Log in as the "postgres" user:

```
$ su - bash
Password: .......
# su - postgres
$
```

Dump the contents of a database to a file by running the following command. Replace yourDbName with the name of the database to be backed up.
You will probably want to dump the database that you used in Assignment 04.


```
$ pg_dump --format=t yourDbName > yourDbName.tar

```

Do a list of the contents of the .tar file and extract the
recover.sql file.

```
$ tar -tf yourDbName.tar >list-of-tar.txt
$ tar -xf yourDbName.tar recover.sql
```

Submit the list-of-tar.txt and the restore.sql file

(If your database is on Windows then you won't have access to the "tar" command.
Please use the Git install of bash to access tar or use 7-zip)





