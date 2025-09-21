CREATE TABLE IF NOT EXISTS warehouse_stocks (
    id              BIGINT PRIMARY KEY AUTO_INCREMENT,
    warehouse_id    BIGINT NOT NULL,
    product_id      BIGINT NOT NULL,
    stock           INT NOT NULL,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB;

CREATE UNIQUE INDEX idx_warehouse_stocks_wh_id_p_id ON warehouse_stocks (warehouse_id, product_id);
CREATE INDEX idx_warehouse_stocks_p_id_wh_id ON warehouse_stocks (product_id, warehouse_id);
