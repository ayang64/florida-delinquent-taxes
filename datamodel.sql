
drop table if exists taxes cascade;
create table taxes (
	id				serial,
	business	text,
	owner			text,
	location	text,
	county		text,
	warrant		text,
	amount		money
);
