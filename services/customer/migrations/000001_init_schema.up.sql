-- enable UUID functions if not exists
create function if not exists UUID_TO_BIN(_uuid char(36))
returns binary(16) deterministic
return unhex(replace(_uuid, '-', ''))

create function if not exists BIN_TO_UUID(_bin binary(16))
returns char(36) deterministic
return lower(concat(
    hex(substr(_bin,1,4)), '-',
    hex(substr(_bin,5,2)), '-',
    hex(substr(_bin,7,2)), '-',
    hex(substr(_bin,9,2)), '-',
    hex(substr(_bin,11))
));

-- create customers table
CREATE TABLE IF NOT EXISTS customers (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255),
    status ENUM('active', 'inactive', 'suspended') DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    INDEX idx_customers_email (email),
    INDEX idx_customers_status (status),
    INDEX idx_customers_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- create subscription table
CREATE TABLE IF NOT EXISTS subscriptions (
    id CHAR(36) PRIMARY KEY DEFAULT (UUID()),
    customer_id CHAR(36) NOT NULL,
    plan_id VARCHAR(50) NOT NULL,
    status ENUM('active', 'canceled', 'past_due', 'unpaid') DEFAULT 'active',
    current_period_start TIMESTAMP NOT NULL,
    current_period_end TIMESTAMP NOT NULL,
    cancel_at_period_end BOOLEAN DEFAULT FALSE,
    stripe_subscription_id VARCHAR(255) UNIQUE NULL,
    metadata JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
    INDEX idx_subscriptions_customer_id (customer_id),
    INDEX idx_subscriptions_status (status),
    INDEX idx_subscriptions_period_end (current_period_end),
    INDEX idx_subscriptions_deleted_at (deleted_at),
    INDEX idx_subscriptions_stripe_id (stripe_subscription_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- Create an audit log table for customer changes
CREATE TABLE IF NOT EXISTS customer_audit_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    customer_id CHAR(36) NOT NULL,
    action ENUM('created', 'updated', 'deleted', 'status_changed') NOT NULL,
    changes JSON,
    performed_by VARCHAR(255) DEFAULT 'system',
    performed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
    INDEX idx_audit_customer_id (customer_id),
    INDEX idx_audit_performed_at (performed_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- create a view for active subscriptions
create or replace view active_subscriptions_view as 
select
    s.*,
    c.email as customer_email,
    c.name as customer_name
from subscriptions s
join customers c on s.customer_id = c.id 
where s.status = 'active'
    and s.deleted_at is NULL
    and c.deleted_at is null;