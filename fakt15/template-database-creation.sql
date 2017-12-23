CREATE TABLE user (
user_id integer PRIMARY KEY,
first_name string not null,
last_name string,
mail string, 
address string,
post_nr_place string,
phone_nr string,
org_nr string,
country_id string);
	INSERT INTO user VALUES(1,'Donald','Duck','donald@andeby.com','Ducksvei 1','1 Andeby',333,'333.333.333',0);
	INSERT INTO user VALUES(2,'Dolly','Duck','dolly@andeby.com','Ducksvei 2','1 Andeby',222,'null',0);
	INSERT INTO user VALUES(3,'Doffen','Duck','doffen@andeby.com','Ducksvei 1','1 Andeby',333,'null',0);
	INSERT INTO user VALUES(4,'Skrue','McDuck','skrue@andeby.com','Pengebingen','1 Andeby',99999999,'999.999.999',0);
	INSERT INTO user VALUES(5,'Mikke','Mus','mikke@andeby.com','1 Musveien','1 Andeby',1432,'null',0
	);

create TABLE country (
	country_id string PRIMARY KEY,
	country_name string,
	vat_value
	);


create TABLE items (
	item_id int PRIMARY KEY,
	item_description string,
	item_name,
	price_ex_vat int,
	storage_amount int
);

create TABLE bill_lines (
	bill_id int PRIMARY KEY,
	item_id int,

);

COMMIT;
