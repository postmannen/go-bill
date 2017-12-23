CREATE TABLE user (
user_id integer PRIMARY KEY,			--pid integer PRIMARY KEY,
first_name string not null,				--firstname string not null,
last_name string,						--lastname string,
mail string, 
address string,
post_nr_place string,					--postnrandplace string,
phone_nr string,						--phonenr string,
org_nr string,							--orgnr string);
country_id integer);							--new
INSERT INTO user VALUES(1,'Donald','Duck','donald@andeby.com','Ducksvei 1','1 Andeby',333,'333.333.333',0);
INSERT INTO user VALUES(2,'Dolly','Duck','dolly@andeby.com','Ducksvei 2','1 Andeby',222,'null',0);
INSERT INTO user VALUES(3,'Doffen','Duck','doffen@andeby.com','Ducksvei 1','1 Andeby',333,'null',0);
INSERT INTO user VALUES(4,'Skrue','McDuck','skrue@andeby.com','Pengebingen','1 Andeby',99999999,'999.999.999',0);
INSERT INTO user VALUES(5,'Mikke','Mus','mikke@andeby.com','1 Musveien','1 Andeby',1432,'null',0);
COMMIT;
/*CREATE TABLE bills (
						pid integer PRIMARY KEY,
						billNR string not null,
						service string not null,
						amount string not null,
						priceExVat integer not null,
						discount integer,
						priceIncVat integer,
						totalSumExVat integer,
						totalSumIncVat integer);
COMMIT;*/