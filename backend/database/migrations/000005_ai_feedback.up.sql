-- Ensure ai_feedback tracking exists for continuous AI improvement
CREATE TABLE ai_feedback (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    edital_id UUID NOT NULL, -- Logical Reference to actual edital
    
    ai_decision VARCHAR(50) NOT NULL, -- ENTER, DO_NOT_ENTER
    real_result VARCHAR(50) NOT NULL, -- WON, LOST, NOT_PARTICIPATED
    
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    UNIQUE (tenant_id, edital_id)
);
