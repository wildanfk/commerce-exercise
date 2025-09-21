CREATE TABLE IF NOT EXISTS orders (
    id              BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id         BIGINT NOT NULL,
    shop_id         BIGINT NOT NULL,
    state           TINYINT NOT NULL,
    total_stock     INT NOT NULL,
    total_price     DECIMAL(15,3) NOT NULL,
    expired_at      DATETIME NOT NULL,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB;

CREATE INDEX idx_orders_user_id ON orders (user_id);
CREATE INDEX idx_orders_shop_id ON orders (shop_id);
CREATE INDEX idx_orders_state_expired_at ON orders (state, expired_at);
