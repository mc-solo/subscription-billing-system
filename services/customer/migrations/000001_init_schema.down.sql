-- drop views first
drop view if exists active_subscriptions_view;


-- drop tables in reverse order
drop table if exists customer_audit_logs;
drop table if exists subscriptions;
drop table if exists customers;

-- drop uuid funcs
drop function if exists UUID_TO_BIN;
drop function if exists BIN_TO_UUID;