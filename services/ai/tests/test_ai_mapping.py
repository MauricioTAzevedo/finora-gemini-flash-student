import pytest
from app.providers import MockAIProvider
from app.models import SpreadsheetMappingRequest, MerchantNormalizationRequest, TransactionCategorizationRequest

@pytest.fixture
def provider():
    return MockAIProvider()

def test_spreadsheet_column_mapping_brazilian_headers(provider):
    # Tests standard Brazilian Excel headers common to household spreadsheets
    headers = ["Data", "Descrição", "Valor Pg.", "Categoria", "Cartão M.", "Observações"]
    req = SpreadsheetMappingRequest(columns=headers)
    res = provider.map_columns(req)

    assert len(res.mappings) == 6

    # Mapping checks
    mapping_dict = {m.source_column: m for m in res.mappings}

    assert mapping_dict["Data"].target_field == "date"
    assert mapping_dict["Data"].confidence >= 0.90
    assert mapping_dict["Data"].requires_confirmation is False

    assert mapping_dict["Descrição"].target_field == "description"
    assert mapping_dict["Descrição"].confidence >= 0.95
    assert mapping_dict["Descrição"].requires_confirmation is False

    assert mapping_dict["Valor Pg."].target_field == "amount"
    assert mapping_dict["Valor Pg."].confidence >= 0.90
    assert mapping_dict["Valor Pg."].requires_confirmation is False

    assert mapping_dict["Categoria"].target_field == "category"
    assert mapping_dict["Categoria"].confidence >= 0.90

    # "Cartão M." is ambiguous and should require user confirmation
    assert mapping_dict["Cartão M."].target_field == "account"
    assert mapping_dict["Cartão M."].requires_confirmation is True

def test_merchant_normalization(provider):
    req1 = MerchantNormalizationRequest(raw_description="AMZN Mktp BR*128490")
    res1 = provider.normalize_merchant(req1)
    assert res1.normalized_merchant == "Amazon"
    assert res1.category_hint == "Assinaturas & Lazer"

    req2 = MerchantNormalizationRequest(raw_description="DROGASIL LOJA 92 SAO PAULO")
    res2 = provider.normalize_merchant(req2)
    assert res2.normalized_merchant == "Drogasil"
    assert res2.category_hint == "Saúde & Farmácia"

def test_transaction_categorization(provider):
    categories = ["Alimentação & Supermercado", "Moradia & Energia", "Transporte & Combustível", "Saúde & Farmácia"]
    req = TransactionCategorizationRequest(
        description="POSTO SHELL AV PAULISTA",
        amount_cents=18000,
        available_categories=categories
    )
    res = provider.categorize_transaction(req)
    assert res.suggested_category == "Transporte & Combustível"
    assert res.confidence >= 0.90
