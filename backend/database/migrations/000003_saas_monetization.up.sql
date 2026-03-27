-- External Stripe Integration Identifiers
ALTER TABLE plans ADD COLUMN stripe_price_id VARCHAR(255);
ALTER TABLE tenants ADD COLUMN stripe_customer_id VARCHAR(255);

-- Alter Subscriptions for External Monetizations
ALTER TABLE subscriptions ADD COLUMN external_subscription_id VARCHAR(255);
ALTER TABLE subscriptions ADD COLUMN current_period_end TIMESTAMP WITH TIME ZONE;
ALTER TABLE subscriptions ADD COLUMN trial_ends_at TIMESTAMP WITH TIME ZONE;

-- Modularity mapping
CREATE TABLE modules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE plan_modules (
    plan_id UUID NOT NULL REFERENCES plans(id) ON DELETE CASCADE,
    module_id UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    PRIMARY KEY(plan_id, module_id)
);

-- Advanced Usage Constraints
CREATE TABLE usage_limits (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    plan_id UUID NOT NULL REFERENCES plans(id) ON DELETE CASCADE,
    resource VARCHAR(100) NOT NULL,
    limit_value INT NOT NULL,
    UNIQUE(plan_id, resource)
);

CREATE TABLE usage_tracking (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    resource VARCHAR(100) NOT NULL,
    current_usage INT NOT NULL DEFAULT 0,
    period_start TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    period_end TIMESTAMP WITH TIME ZONE NOT NULL,
    UNIQUE(tenant_id, resource)
);

