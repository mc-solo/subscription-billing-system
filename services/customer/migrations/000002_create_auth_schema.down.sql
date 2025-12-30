-- rollback: removes authentication chema
-- this removes the users table and related changes

-- drops the audit first
drop table if exists user_auth_audit_logs;

-- remove the fk constriants 
alter table customers
drop foreign key fk_customers_user,
drop index idx_customers_user_id,
drop column user_id;

-- and then drop the user table
drop table if exists users;