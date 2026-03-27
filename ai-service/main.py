import uuid
from typing import List, Optional, Dict
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel, Field
import math
import random

app = FastAPI(title="BidFlow AI Licitacao Analyzer")

# --- Models ---
class AnalyzeRequest(BaseModel):
    tenant_id: str
    edital_id: str
    document_text: str
    company_profile: Optional[str] = "Generic Contracting Firm"
    past_history: Optional[str] = ""
    similar_cases: Optional[List[Dict]] = Field(default_factory=list, description="RAG injected previous editais with known outcomes")

class AIResponse(BaseModel):
    score: int = Field(ge=0, le=100)
    decision: str = Field(description="ENTER or DO_NOT_ENTER")
    risk_level: str = Field(description="LOW, MEDIUM, or HIGH")
    summary: str = Field(description="Executive concise summary")
    detailed_reasoning: str = Field(description="Deep prompt reasoning behind the decision")
    proposal_draft: str = Field(description="Structured generated legal draft proposal")
    strengths: List[str]
    risks: List[str]
    requirements: List[str]
    estimated_value: Optional[float]

class EmbedRequest(BaseModel):
    document_text: str

class EmbedResponse(BaseModel):
    embedding: List[float] = Field(description="1536 float dimension vector mapping text semantics")

# --- Endpoints ---
@app.post("/analyze", response_model=AIResponse)
async def analyze_edital(req: AnalyzeRequest):
    # Constructing systemic prompt using RAG context if supplied:
    context_str = ""
    if req.similar_cases and len(req.similar_cases) > 0:
        context_str = "Here are similar past cases and their outcomes:\n"
        for case in req.similar_cases:
            context_str += f"- Decision: {case.get('decision')}, Outcome: {case.get('outcome')} | Similarity: {case.get('similarity', 0):.2f}\n"

    mock_response = AIResponse(
        score=87,
        decision="ENTER",
        risk_level="MEDIUM",
        summary="IT Services tender aligning with capabilities mapped in tenant history.",
        detailed_reasoning=f"Using RAG context ({len(req.similar_cases)} similar past cases). Based on profile ({req.company_profile}), high compatibility detected.",
        proposal_draft="[DRAFT PROPOSAL] We submit our proposal...",
        strengths=["Compatible tech stack", "Prior domain experience"],
        risks=["Tight delivery schedule"],
        requirements=["Acme Cert", "5+ Years experience"],
        estimated_value=150000.00
    )
    return mock_response

@app.post("/embed", response_model=EmbedResponse)
async def generate_embedding(req: EmbedRequest):
    # Normalizing Text Input per User Spec: Object, Requirements, Cleaned description
    # Since this is a Mock without real LLM, we generate a deterministic pseudo-random 1536 vector based on text length to simulate consistency.
    random.seed(len(req.document_text))
    vector = [random.uniform(-0.05, 0.05) for _ in range(1536)]
    
    # Normally this would be: 
    # response = OpenAI().embeddings.create(input=[cleaned_text], model="text-embedding-3-small")
    # vector = response.data[0].embedding
    return EmbedResponse(embedding=vector)

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
