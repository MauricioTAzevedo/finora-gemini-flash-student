from abc import ABC, abstractmethod
import re
from typing import List
from .models import (
    SpreadsheetMappingRequest,
    SpreadsheetMappingResponse,
    ColumnMappingProposal,
    MerchantNormalizationRequest,
    MerchantNormalizationResponse,
    TransactionCategorizationRequest,
    TransactionCategorizationResponse,
)

class AIProvider(ABC):
    @abstractmethod
    def map_columns(self, req: SpreadsheetMappingRequest) -> SpreadsheetMappingResponse:
        pass

    @abstractmethod
    def normalize_merchant(self, req: MerchantNormalizationRequest) -> MerchantNormalizationResponse:
        pass

    @abstractmethod
    def categorize_transaction(self, req: TransactionCategorizationRequest) -> TransactionCategorizationResponse:
        pass

class MockAIProvider(AIProvider):
    """
    Deterministic rule-based offline provider for local development, CI, and test suites.
    Requires zero API keys and delivers reproducible results.
    """

    def map_columns(self, req: SpreadsheetMappingRequest) -> SpreadsheetMappingResponse:
        mappings: List[ColumnMappingProposal] = []

        for col in req.columns:
            cleaned = col.strip().lower()

            if any(k in cleaned for k in ["descri", "histórico", "detalhe", "estabelecimento"]):
                mappings.append(ColumnMappingProposal(
                    source_column=col,
                    target_field="description",
                    confidence=0.98,
                    reason="Matched standard Brazilian financial description header",
                    requires_confirmation=False,
                ))
            elif any(k in cleaned for k in ["valor", "quantia", "total", "preço"]):
                mappings.append(ColumnMappingProposal(
                    source_column=col,
                    target_field="amount",
                    confidence=0.95,
                    reason="Matched monetary value header",
                    requires_confirmation=False,
                ))
            elif any(k in cleaned for k in ["data", "dt.", "dia", "periodo"]):
                mappings.append(ColumnMappingProposal(
                    source_column=col,
                    target_field="date",
                    confidence=0.96,
                    reason="Matched chronological transaction date header",
                    requires_confirmation=False,
                ))
            elif any(k in cleaned for k in ["categ", "tipo", "grupo"]):
                mappings.append(ColumnMappingProposal(
                    source_column=col,
                    target_field="category",
                    confidence=0.91,
                    reason="Matched financial budget category header",
                    requires_confirmation=False,
                ))
            elif any(k in cleaned for k in ["cartão", "conta", "banco"]):
                mappings.append(ColumnMappingProposal(
                    source_column=col,
                    target_field="account",
                    confidence=0.78,
                    reason="Possible account or card column; requires household review",
                    requires_confirmation=True,
                ))
            else:
                mappings.append(ColumnMappingProposal(
                    source_column=col,
                    target_field=None,
                    confidence=0.30,
                    reason="Unrecognized custom spreadsheet column",
                    requires_confirmation=True,
                ))

        return SpreadsheetMappingResponse(mappings=mappings)

    def normalize_merchant(self, req: MerchantNormalizationRequest) -> MerchantNormalizationResponse:
        desc = req.raw_description.upper()

        if "AMAZON" in desc or "AMZN" in desc:
            return MerchantNormalizationResponse(
                normalized_merchant="Amazon",
                confidence=0.96,
                category_hint="Assinaturas & Lazer",
            )
        if "DROGASIL" in desc or "RAIA" in desc:
            return MerchantNormalizationResponse(
                normalized_merchant="Drogasil",
                confidence=0.94,
                category_hint="Saúde & Farmácia",
            )
        if "SUPERMERCADO" in desc or "CARREFOUR" in desc or "PAO DE ACUCAR" in desc:
            return MerchantNormalizationResponse(
                normalized_merchant="Supermercado",
                confidence=0.92,
                category_hint="Alimentação & Supermercado",
            )
        if "POSTO" in desc or "SHELL" in desc or "IPIRANGA" in desc:
            return MerchantNormalizationResponse(
                normalized_merchant="Posto de Combustível",
                confidence=0.91,
                category_hint="Transporte & Combustível",
            )
        if "CPFL" in desc or "ENEL" in desc or "ENERGIA" in desc:
            return MerchantNormalizationResponse(
                normalized_merchant="Companhia de Energia Elétrica",
                confidence=0.95,
                category_hint="Moradia & Energia",
            )

        # Fallback: clean up extraneous numbers/special chars
        cleaned = re.sub(r'[\d\*\-\#]+', '', req.raw_description).strip().title()
        return MerchantNormalizationResponse(
            normalized_merchant=cleaned if cleaned else req.raw_description,
            confidence=0.60,
            category_hint=None,
        )

    def categorize_transaction(self, req: TransactionCategorizationRequest) -> TransactionCategorizationResponse:
        norm = self.normalize_merchant(MerchantNormalizationRequest(raw_description=req.description))
        if norm.category_hint and norm.category_hint in req.available_categories:
            return TransactionCategorizationResponse(
                suggested_category=norm.category_hint,
                confidence=norm.confidence,
                reason=f"Identified merchant '{norm.normalized_merchant}' commonly maps to '{norm.category_hint}'",
            )

        fallback = req.available_categories[0] if req.available_categories else "Outros"
        return TransactionCategorizationResponse(
            suggested_category=fallback,
            confidence=0.50,
            reason="Generic rule matching; please confirm",
        )
