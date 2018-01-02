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
	INSERT INTO user VALUES(5,'Mikke','Mus','mikke@andeby.com','1 Musveien','1 Andeby',1432,'null',0);
	INSERT INTO user VALUES(7,'Kit','Walker','kit@fantomet.com','Hodeskallegrotten','De dype skoger','Apepost','null',0);

CREATE TABLE country (
	country_id string PRIMARY KEY,
	country_name string,
	vat_percent int
	);


CREATE TABLE items (
	item_id int PRIMARY KEY,
	item_description string,
	item_name string,
	price_ex_vat real,
	storage_amount int
);

CREATE TABLE bill_lines (
	indx int PRIMARY KEY,
	bill_id int,
	line_id int,
	item_id int,
	description string,
	quantity int,
	discount_percentage int,
	vat_used int,
	price_ex_vat real
);

CREATE TABLE bills (
	bill_id int PRIMARY KEY,
	user_id int,
	created_date text,
	due_date text,
	comment string,
	totalt_ex_vat real,
	total_inc_vat real
);

COMMIT;
