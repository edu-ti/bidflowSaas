-- pgvector requirement for storing similarity embeddings
CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE ai_embeddings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    edital_id UUID NOT NULL,
    
    embedding VECTOR(1536), -- Assuming standard OpenAI embedding model 3-small boundaries
    text_hash TEXT NOT NULL,
    model_version VARCHAR(50) NOT NULL DEFAULT 'text-embedding-3-small',
    
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    
    UNIQUE (tenant_id, edital_id)
);

-- Optimize high performance neighbor search scaling
CREATE INDEX ai_embeddings_ivfflat_idx ON ai_embeddings USING ivfflat (embedding vector_cosine_ops);
