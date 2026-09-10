# Synthetic Financial Fixtures

These sample files contain purely synthetic financial transactions mimicking standard Brazilian bank statements and Excel spreadsheets for testing the Finora Ingestion & Reconciliation engine.

- `extrato_nubank_agosto.ofx`: Standard OFX SGML format from Brazilian banking apps with FITID, DTPOSTED, TRNAMT, and MEMO.
- `planilha_gastos_familia.csv`: Semicolon-delimited CSV with Brazilian decimal formatting (`R$ 1.250,50`) and pt-BR headers.

## Testing via API
```bash
# Test OFX reconciliation
curl -X POST "http://localhost:8080/api/v1/imports/reconcile?format=ofx" \
  -H "X-Household-ID: b0000000-0000-0000-0000-000000000001" \
  --data-binary "@database/samples/extrato_nubank_agosto.ofx"

# Test CSV reconciliation
curl -X POST "http://localhost:8080/api/v1/imports/reconcile?format=csv" \
  -H "X-Household-ID: b0000000-0000-0000-0000-000000000001" \
  --data-binary "@database/samples/planilha_gastos_familia.csv"
```
