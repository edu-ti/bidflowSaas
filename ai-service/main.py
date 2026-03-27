import uuid
from typing import List, Optional
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel, Field
import os
import json

app = FastAPI(title="BidFlow AI Licitacao Analyzer")

# --- Models ---
class AnalyzeRequest(BaseModel):
    tenant_id: str
    edital_id: str
    document_text: str
    company_profile: Optional[str] = "Generic Contracting Firm"
    past_history: Optional[str] = ""

class AIResponse(BaseModel):
    score: int = Field(ge=0, le=100, description="1-100 Compatibility Score based on tenant profile")
    decision: str = Field(description="ENTER or DO_NOT_ENTER based strictly on matching specs")
    risk_level: str = Field(description="LOW, MEDIUM, or HIGH risk evaluation")
    summary: str = Field(description="Executive concise summary of the tender")
    detailed_reasoning: str = Field(description="Deep prompt reasoning behind the decision")
    proposal_draft: str = Field(description="A highly-structured generated legal draft proposal based on features")
    strengths: List[str]
    risks: List[str]
    requirements: List[str]
    estimated_value: Optional[float]

# --- Endpoints ---
@app.post("/analyze", response_model=AIResponse)
async def analyze_edital(req: AnalyzeRequest):
    # In a production layout, `client = OpenAI(api_key=...)` would be here
    # We would use `client.chat.completions.create(..., response_format={ "type": "json_object" })`
    # and heavily feed `req.company_profile` and `req.past_history` into the System Prompt.
    
    # Mocking execution to simulate standard response boundaries
    prompt_version = "v1.2-b2g"
    
    mock_response = AIResponse(
        score=87,
        decision="ENTER",
        risk_level="MEDIUM",
        summary="IT Services tender perfectly aligning with infrastructure capabilities mapped in tenant history.",
        detailed_reasoning=f"Based on the company profile ({req.company_profile}) and past wins in the public sector, the requirements heavily overlap with existing capacity. Only financial qualification proofs pose a slight delay risk.",
        proposal_draft="[DRAFT PROPOSAL] We hereby submit our proposal responding strictly to Edital {req.edital_id} fulfilling all compliance layers...",
        strengths=["Highly compatible tech stack", "Prior domain experience", "Favorable SLA boundaries"],
        risks=["Tight delivery schedule", "Heavy penalty clause on miss"],
        requirements=["Acme Cert", "5+ Years experience"],
        estimated_value=150000.00
    )
    
    # Normally we would log the full raw conversation loop to a secure audit DB directly or return it for the worker to log.
    return mock_response

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
