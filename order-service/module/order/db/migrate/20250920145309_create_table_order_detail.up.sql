CREATE TABLE IF NOT EXISTS order_details (
    id              BIGINT PRIMARY KEY AUTO_INCREMENT,
    order_id        BIGINT NOT NULL,
    product_id      BIGINT NOT NULL,
    warehouse_id    BIGINT NOT NULL,
    stock           INT NOT NULL,
    price           DECIMAL(15,3) NOT NULL,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB;

CREATE INDEX idx_order_details_order_id ON order_details (order_id);
