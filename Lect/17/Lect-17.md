
<style>
.pagebreak { page-break-before: always; }
.half { height: 200px; }
</style>


# Lecture 17 - Hash Based Databases v.s. SqLite

[SqLite3 v.s. Hash Store - https://youtu.be/0GG2s9uz8BY](https://youtu.be/0GG2s9uz8BY)<br>
[get-issue-detail  walk through - https://youtu.be/lfCNVEJEOxQ](https://youtu.be/lfCNVEJEOxQ)<br>

From Amazon S3 - for download (same as youtube videos)

[SqLite3 v.s. Hash Store](http://uw-s20-2015.s3.amazonaws.com/4820-L17-pt1-sqlite3-very-brief-overivew.mp4)<br>
[get-issue-detail  walk through](http://uw-s20-2015.s3.amazonaws.com/4820-L17-pt2-get-issue-detail-walkthrough.mp4)<br>

The stated reason for using a "hash" database like
"bitcoin" uses is that you don't have an external
database that requires setup and maintenance.  
Bitcoin uses LevelDB a hash-file store.

SqLite is a SQL database that stores its data in a single file.
All of you are already users of SqLite - it is builtin to 
Google Chrome, Firefox, Microsoft Edge and Apple Safari.
This is how your bookmarks are stored (And a lot of other
stuff).

SqLite is designed to be compiled right into an application.
So if you need to ship a "demo" of a SQL" tool with the data
and you don't want to force the demo-user to setup an entire
database - this is the tool for you.

It is reasonably fast - I have run websites with 100's of requests
a minute on it.

It works at scale for data.   You can lode gigabyte size data.
Inside a single application it is ACID compliant - even
with multiple threads.


Let's take a look at performance.

| Database/Rows | 10		| 100		| 1000		| 10000		|
|--------------:|----------:|----------:|----------:|----------:|
| SqLite3		| 8			| 47		| 391		| 3141		|
| LevelDb		| 5			| 52		| 431		| 5210		|
| PostgreSQL	| 19		| 61		| 532		| 4328		|

Known "issue" or problems that I have encountered.

1. Inserting a NULL row will not do an insert.
2. Porting data files between big-ndian and little-ndian systems is "iffy".
3. Performance is not the same as 50,000 TPS from PostgreSQL - probably more on the 500-1000 TPS range.





Copyright (C) University of Wyoming, 2021
