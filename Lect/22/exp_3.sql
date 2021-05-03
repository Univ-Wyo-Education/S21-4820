CREATE SEQUENCE exp_3_id_seq
  INCREMENT 1
  MINVALUE 1
  MAXVALUE 9223372036854775807
  START 1
  CACHE 1;

create table exp_3 (
	id			bigint DEFAULT nextval('exp_3_id_seq'::regclass) NOT NULL primary key,
	user_data 	bigint
);

explain analyze
	select *
	from exp_3
	where id = 55
;

explain analyze
	select *
	from exp_3
	where user_data = 55
;
