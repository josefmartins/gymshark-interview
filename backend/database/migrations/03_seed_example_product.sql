-- +migrate Up

INSERT INTO products (id,name)
VALUES ("0196b5d3-c52c-7e50-ac45-f83b35ee9e3d","Product ABC");

INSERT INTO package_sizes (id, product_id, size)
VALUES ("0196b5d1-9010-74de-8f3e-f11149df2319","0196b5d3-c52c-7e50-ac45-f83b35ee9e3d",250),
("0196b5d1-f64b-737c-8cf6-42c3ea81cf1b","0196b5d3-c52c-7e50-ac45-f83b35ee9e3d",500), 
("0196b5d1-f64b-710e-bb44-49fb981fd7a7","0196b5d3-c52c-7e50-ac45-f83b35ee9e3d",1000),
("0196b5d1-f64b-7632-b865-0c46cb2dfe44","0196b5d3-c52c-7e50-ac45-f83b35ee9e3d",2000),
("0196b5d1-f64b-784d-bd68-23f35eeaf73a","0196b5d3-c52c-7e50-ac45-f83b35ee9e3d",5000);

-- +migrate Down

DELETE FROM package_sizes WHERE id IN ("0196b5d1-9010-74de-8f3e-f11149df2319","0196b5d1-f64b-737c-8cf6-42c3ea81cf1b","0196b5d1-f64b-710e-bb44-49fb981fd7a7","0196b5d1-f64b-7632-b865-0c46cb2dfe44","0196b5d1-f64b-784d-bd68-23f35eeaf73a");

DELETE FROM products WHERE id = "0196b5d3-c52c-7e50-ac45-f83b35ee9e3d";

