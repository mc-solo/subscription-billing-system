-- this migration creates the users table for authentication
-- and links it to the customers table

create table if not exists users (
    id char(36) primary key default (uuid()),
    email varchar(255) unique not null,
    password_hash varchar(255) not null,
    full_name varchar(255),
    status enum('pending_verification', 'active', 'suspended') default 'pending_verification',


    -- email verification
    email_verified boolean default false,
    verification_token varchar(100),
    verification_token_expires_at timestamp null,

    -- password reset functionality
    password_reset_token varchar(100),
    password_reset_expires_at timestamp null,

    -- security and login tracking
    last_login_at timestamp null,
    failed_login_attempts int default 0,
    locked_until timestamp null,

    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp on update current_timestamp,
    deleted_at timestamp null default null,

    -- stripe integration
    stripe_customer_id varchar(255),

    -- indexes for performance
    index idx_users_email (email),
    index idx_users_status (status),
    index idx_users_verification_token (verification_token),
    index idx_users_password_reset_token (password_reset_token),
    index idx_users_stripe_customer_id (stripe_customer_id),
    index idx_users_deleted_at (deleted_at)


) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- link users to the existing customers table
alter table customers add column user_id char(36) null after id,
add constraint fk_customers_user
    foreign key (user_id) references users(id) on delete set null,
add index idx_customers_user_id (user_id);

-- create an audit log table for user authentication events 
create table if not exists user_auth_audit_logs (
    id bigint primary key auto_increment,
    user_id char(36) not null,
    action enum('registered', 'logged_in', 'logged_out', 'password_changed', 'email_verified', 'account_locked', 'account_unlocked') not null,
    ip_address varchar(45), 
    user_agent text,
    metadata json,
    performed_at timestamp default current_timestamp,

    foreign key (user_id) references users(id) on delete cascade,
    index idx_auth_user_id (user_id),
    index idx_auth_performed_at (performed_at),
    index idx_auth_action (action)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;