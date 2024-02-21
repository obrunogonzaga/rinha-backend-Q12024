CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    limite INT NOT NULL,
    saldo INT NOT NULL
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL,
    amount INT NOT NULL,
    transaction_type CHAR(1) NOT NULL,
    description VARCHAR(10) NOT NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_customer_transaction_id
        FOREIGN KEY (customer_id) REFERENCES customers(id)
);

INSERT INTO customers (id, limite, saldo)
VALUES
    (1, 100000, 0),
    (2, 80000, 0),
    (3, 1000000, 0),
    (4, 10000000, 0),
    (5, 500000, 0);
