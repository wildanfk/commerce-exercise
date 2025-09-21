CREATE TABLE IF NOT EXISTS warehouses (
    id              BIGINT PRIMARY KEY AUTO_INCREMENT,
    shop_id         BIGINT NOT NULL,
    name            VARCHAR(255) NOT NULL,
    active          BOOLEAN,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB;

CREATE INDEX idx_warehouses_shop_id ON warehouses (shop_id);
CREATE INDEX idx_warehouses_active ON warehouses (active);
