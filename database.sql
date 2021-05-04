-- Database: Rational_Room_Reservations
-- DROP DATABASE "Rational_Room_Reservations";

CREATE DATABASE "Rational_Room_Reservations"
    WITH 
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.UTF-8'
    LC_CTYPE = 'en_US.UTF-8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;
	
CREATE OR REPLACE FUNCTION pseudo_encrypt(VALUE int) returns int AS $$
DECLARE
l1 int;
l2 int;
r1 int;
r2 int;
i int:=0;
BEGIN
 l1:= (VALUE >> 16) & 65535;
 r1:= VALUE & 65535;
 WHILE i < 3 LOOP
   l2 := r1;
   r2 := l1 # ((((1366 * r1 + 150889) % 714025) / 714025.0) * 32767)::int;
   l1 := l2;
   r1 := r2;
   i := i + 1;
 END LOOP;
 RETURN ((r1 << 16) + l1);
END;
$$ LANGUAGE plpgsql strict immutable;

CREATE SEQUENCE seq maxvalue 2147483647;

CREATE TABLE departments (
	department_id INT DEFAULT pseudo_encrypt(nextval('seq')::INT) UNIQUE PRIMARY KEY NOT NULL,
	"name" VARCHAR(100) NOT NULL
);

CREATE TABLE buildings	(
	building_id INT DEFAULT pseudo_encrypt(nextval('seq')::INT) UNIQUE PRIMARY KEY NOT NULL,
	department_id INT REFERENCES departments (department_id),
	"name" VARCHAR(100) NOT NULL,
	address VARCHAR(100) NOT NULL
);

CREATE TABLE offices	(
	office_id INT DEFAULT pseudo_encrypt(nextval('seq')::INT) UNIQUE PRIMARY KEY NOT NULL,
	department_id INT REFERENCES departments (department_id),
	"name" VARCHAR(100) NOT NULL,
	priority_level INT NOT NULL
);

CREATE TABLE users	(
	user_id INT DEFAULT pseudo_encrypt(nextval('seq')::INT) UNIQUE PRIMARY KEY NOT NULL,
	office_id INT REFERENCES offices (office_id),
	email VARCHAR(50) NOT NULL,
	password_hash VARCHAR(50) NOT NULL,
	phone VARCHAR(20) NOT NULL,
	first_name VARCHAR(50) NOT NULL,
	last_name VARCHAR(50) NOT NULL,
	access_level INT NOT NULL,
	is_deleted BOOL NOT NULL
);

CREATE TABLE user_proxies	(
	user_proxy_id INT DEFAULT pseudo_encrypt(nextval('seq')::INT) UNIQUE PRIMARY KEY NOT NULL,
	parent_user INT REFERENCES users (user_id),
	access_level INT NOT NULL,
	password_hash VARCHAR(50) NOT NULL
);

CREATE TABLE sectional_rooms (
	sectional_room_id INT DEFAULT pseudo_encrypt(nextval('seq')::INT) UNIQUE PRIMARY KEY NOT NULL,
	"name" VARCHAR(50) NOT NULL,
	quantity_of_rooms INT NOT NULL
);

CREATE TABLE rooms (
	room_id INT DEFAULT pseudo_encrypt(nextval('seq')::INT) UNIQUE PRIMARY KEY NOT NULL,
	building_id INT REFERENCES buildings (building_id),
	office_id INT REFERENCES offices (office_id),
	sectional_room_id INT REFERENCES sectional_rooms (sectional_room_id),
	price FLOAT8 NOT NULL,
	room_number INT,
	"floor" INT NOT NULL,
	height INT,
	"size" INT NOT NULL,
	address VARCHAR(50) NOT NULL,
	description TEXT,
	gis_coordinate VARCHAR(50) NOT NULL,
	reservation_lock BOOL NOT NULL
);

CREATE TABLE equipment	(
	equipment_id INT DEFAULT pseudo_encrypt(nextval('seq')::INT) UNIQUE PRIMARY KEY NOT NULL,
	room_id INT REFERENCES rooms (room_id),
	price FLOAT8 NOT NULL,
	description VARCHAR(100) NOT NULL,
	"type" VARCHAR(50) NOT NULL,
	unavailable BOOL NOT NULL,
	fixed BOOL NOT NULL
);

CREATE TABLE reservations	(
	reservation_id INT DEFAULT pseudo_encrypt(nextval('seq')::INT) UNIQUE PRIMARY KEY NOT NULL,
	room_id INT REFERENCES rooms (room_id),
	user_id INT REFERENCES users (user_id),
	priority_level INT NOT NULL,
	account_number INT NOT NULL,
	start_time TIMESTAMP without time zone NOT NULL,
	end_time TIMESTAMP without time zone NOT NULL,
	purpose VARCHAR(100) NOT NULL,
	number_of_people INT NOT NULL,
	created_at TIMESTAMP without time zone DEFAULT now() NOT NULL,
	cancelled_pending_reopen TIMESTAMP without time zone DEFAULT now() NOT NULL,
	is_deleted BOOL NOT NULL
);

CREATE TABLE reservation_attendees	(
	reservation_id INT PRIMARY KEY REFERENCES reservations (reservation_id),
	user_id INT REFERENCES users (user_id)
);

CREATE TABLE reserved_equipment	(
	equipment_id INT PRIMARY KEY REFERENCES equipment (equipment_id),
	reservation_id INT REFERENCES reservations (reservation_id)
);


-- SELECT * FROM departments;
-- SELECT * FROM buildings;
-- SELECT * FROM offices;
-- SELECT * FROM users;
-- SELECT * FROM user_proxies;
-- SELECT * FROM sectional_rooms;
-- SELECT * FROM rooms;
-- SELECT * FROM equipment;
-- SELECT * FROM reservations;
-- SELECT * FROM reservation_attendees;
-- SELECT * FROM reserved_equipment;

-- DROP TABLE reserved_equipment;
-- DROP TABLE reservation_attendees;
-- DROP TABLE reservations;
-- DROP TABLE equipment;
-- DROP TABLE rooms;
-- DROP TABLE sectional_rooms;
-- DROP TABLE user_proxies;
-- DROP TABLE users;
-- DROP TABLE offices;
-- DROP TABLE buildings;
-- DROP TABLE departments;