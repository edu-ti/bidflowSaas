DROP TABLE IF EXISTS usage_tracking;
DROP TABLE IF EXISTS usage_limits;
DROP TABLE IF EXISTS plan_modules;
DROP TABLE IF EXISTS modules;

ALTER TABLE subscriptions DROP COLUMN external_subscription_id;
ALTER TABLE subscriptions DROP COLUMN current_period_end;
ALTER TABLE subscriptions DROP COLUMN trial_ends_at;
