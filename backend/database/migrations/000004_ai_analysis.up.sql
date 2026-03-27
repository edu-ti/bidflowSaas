-- Comprehensive AI Analytics and Decision Storage
CREATE TABLE ai_analysis_results (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    edital_id UUID NOT NULL, -- Logical Reference to actual edital
    
    status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, processing, completed, failed
    score INT,
    decision VARCHAR(50), -- ENTER, DO_NOT_ENTER
    risk_level VARCHAR(20), -- LOW, MEDIUM, HIGH
    summary TEXT,
    detailed_reasoning TEXT,
    proposal_draft TEXT,
    
    structured_data JSONB, -- Requirements, Strengths, Risks Arrays
    raw_response TEXT, -- Log raw LLM response natively for audit
    prompt_version VARCHAR(50) NOT NULL DEFAULT 'v1',
    error TEXT,
    
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    UNIQUE (tenant_id, edital_id)
);

-- Note: We map actual usage metrics inside usage_tracking independently directly on job success.
