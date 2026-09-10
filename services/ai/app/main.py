from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from .models import (
    SpreadsheetMappingRequest,
    SpreadsheetMappingResponse,
    MerchantNormalizationRequest,
    MerchantNormalizationResponse,
    TransactionCategorizationRequest,
    TransactionCategorizationResponse,
)
from .providers import MockAIProvider

app = FastAPI(
    title="Finora Gemini Flash Student — AI & Intelligence Service",
    version="1.0.0",
    description="Intelligent spreadsheet mapping, merchant normalization, and financial document extraction."
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

provider = MockAIProvider()

@app.get("/health")
def health_check():
    return {"status": "healthy", "service": "finora-ai-service", "provider": "mock"}

@app.post("/api/v1/ai/map-spreadsheet-columns", response_model=SpreadsheetMappingResponse)
def map_spreadsheet_columns(request: SpreadsheetMappingRequest):
    return provider.map_columns(request)

@app.post("/api/v1/ai/normalize-merchant", response_model=MerchantNormalizationResponse)
def normalize_merchant(request: MerchantNormalizationRequest):
    return provider.normalize_merchant(request)

@app.post("/api/v1/ai/categorize-transaction", response_model=TransactionCategorizationResponse)
def categorize_transaction(request: TransactionCategorizationRequest):
    return provider.categorize_transaction(request)
