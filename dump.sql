create table if not exists storage
(
	product_id int not null
		constraint storage_pk
			primary key
	quantity int default 0 not null
);

insert into storage (product_id, quantity) values
(1, 10),
(2, 55),
(3, 300),
(4, 40),
(5, 75),
(6, 15),
(7, 155);