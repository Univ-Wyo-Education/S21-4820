

m4_include(../lect-setup.m4)

# Lecture 11 - Security

## News

[https://www.infoq.com/news/2021/02/amazon-aurora-postgresql-12/](https://www.infoq.com/news/2021/02/amazon-aurora-postgresql-12/)

Amazon AWS has a modified version of PostgreSQL 12 - that runs much faster on AWS than a regular install.  The Regular install is available via
Amazon RDS - Relational Data Services.  The Aurora version is 2x to 5x faster than regular and 3x faster than mySQL.  That is for a single
instance running.   You also have a auto-scale capability that allows up to 15 read-replicas.   

Also there is an auto-upgrade capability.  The primary reason that the NyTimes claims that they used mySQL instead of PostgreSQL is because 
the upgrades in PostreSQL require down time.  I tried the Upgrade from a 11x PostgreSQL to the new 12x vision it was - as far as I can tell -
invisible to the users.   I did it on a running system - with active users on it - without interruption of the system.

[https://www.depesz.com/2021/02/04/waiting-for-postgresql-14-search-and-cycle-clauses/](https://www.depesz.com/2021/02/04/waiting-for-postgresql-14-search-and-cycle-clauses/)

This new recursive select capability is way in advance of all other databases giving you depth-first, or breath-first or circuit searches.  This allows for
non-trees - like graph search to take place.

[https://www.jdsupra.com/legalnews/why-it-matters-whether-hashed-passwords-3305065/](https://www.jdsupra.com/legalnews/why-it-matters-whether-hashed-passwords-3305065/)

A former client (well... now a current client) of mine brought up this problem.   How to address the legal security needs and store passwords.
The client is switching to using Secure Remote Password version 6a (SRP6a) at a future point - I developed the code for them to do the
conversion - but they are not ready to take the plunge.  SO they wanted an interim solution.  This involved doing encryption of the data
and protects other data in the database.

## How to implement encryption based security

https://blog.andreiavram.ro/encrypt-postgresql-column/

https://medium.com/@jianshi_94445/encrypt-decrypt-your-data-using-native-postgres-functions-with-sequelize-js-c04948d96833


