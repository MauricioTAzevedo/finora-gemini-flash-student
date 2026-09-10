from typing import List, Optional
from pydantic import BaseModel, Field

class ColumnMappingProposal(BaseModel):
    source_column: str = Field(description="Name of the header in user's spreadsheet")
    target_field: Optional[str] = Field(description="Mapped Finora canonical field: date, description, amount, category, account")
    confidence: float = Field(ge=0.0, le=1.0, description="Confidence score from 0.0 to 1.0")
    reason: str = Field(description="Explainable rationale for the mapping")
    requires_confirmation: bool = Field(description="True if confidence is below safety threshold (<0.85)")

class SpreadsheetMappingRequest(BaseModel):
    columns: List[str] = Field(description="List of header strings extracted from spreadsheet")
    sample_rows: Optional[List[List[str]]] = Field(default=None, description="Optional sample row values for type inference")

class SpreadsheetMappingResponse(BaseModel):
    mappings: List[ColumnMappingProposal]

class MerchantNormalizationRequest(BaseModel):
    raw_description: str

class MerchantNormalizationResponse(BaseModel):
    normalized_merchant: str
    confidence: float
    category_hint: Optional[str] = None

class TransactionCategorizationRequest(BaseModel):
    description: str
    amount_cents: int
    available_categories: List[str]

class TransactionCategorizationResponse(BaseModel):
    suggested_category: str
    confidence: float
    reason: str
