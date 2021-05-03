create index exp_3_p1 on exp_3 ( user_data );

explain analyze
	select *
	from exp_3
	where user_data = 55
;
