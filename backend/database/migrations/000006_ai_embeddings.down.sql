DROP INDEX IF EXISTS ai_embeddings_ivfflat_idx;
DROP TABLE IF EXISTS ai_embeddings;
-- Cannot safely drop EXTENSION vector automatically as other tables could implicitly inherit it in the future, safe to leave.
