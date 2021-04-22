
<style>
.pagebreak { page-break-before: always; }
.half { height: 200px; }
</style>


# Lecture 01 - Intro to Class

[Syllabus part 1 - https://youtu.be/HGtT48bdxII](https://youtu.be/HGtT48bdxII)<br>
[Lecture 2 - Why Database - https://youtu.be/z6dy9bi0piU](https://youtu.be/z6dy9bi0piU)<br>
[Lecture 2 - What is a Database - https://youtu.be/zoS8YWuXzss](https://youtu.be/zoS8YWuXzss)<br>

From Amazon S3 - for download (same as youtube videos)

[Syllabus part 1](http://uw-s20-2015.s3.amazonaws.com/4280-L01-Syllabus.mp4)<br>
[Lecture 2 - Why Database](http://uw-s20-2015.s3.amazonaws.com/4280-L02-pt1-Why-Database.mp4)<br>
[Lecture 2 - What is a Database](http://uw-s20-2015.s3.amazonaws.com/4280-L02-what-is-a-database.mp4)<br>


## Why Database

1. Most jobs require it
2. Lots of businesses can be built with it

## What are databases

There are a number of different kinds of databases.

1. SQL - this is what we will spend a large portion of the semester on.
2. NoSQL - The most popular is MongoDB - we will look at it a little.
3. Redis - Shared Data Structure
4. File System, Amazon S3, and IPFS
5. Object Store
6. Blockchain

Most applications need some sort of persistent data.  This is most often
stored in a database.

Usually a database class starts out with the database and that's it - you
don't get any context of where it all fits in.   That is not how we are going
to do it.   We are going to start out with an application and burro down into
where the database fits into the bigger picture.  By doing that we are 
going to see that just a SQL database is usually not a complete solution.

## Tables are 

A set of columns with types and a set of rows.

