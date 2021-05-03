create table exp_1 (
	i int
);

EXPLAIN 
	SELECT * 
	FROM exp_1 
	WHERE i = 0;
