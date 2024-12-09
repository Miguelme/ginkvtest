CREATE TABLE key_value_store (
     key_field VARCHAR(255) PRIMARY KEY,
     value_field VARCHAR(255)
);

INSERT INTO key_value_store (key_field, value_field) VALUES ('test_key', 'test_value');
