-- ==============================================================================
-- Finora Gemini Flash Student — Demo Dataset ("Família Silva")
-- ==============================================================================

-- 1. Demo User: Maurício (mauricio@example.com)
INSERT INTO users (id, email, password_hash, name)
VALUES (
    'a0000000-0000-0000-0000-000000000001',
    'mauricio@example.com',
    '$2a$10$7EqJtq98hPqEX7fNZaFWoO.8/b8GqLp5m7E0u3m2qQ6k1o5wZ2K9K', -- demo123456
    'Maurício Silva'
) ON CONFLICT (email) DO NOTHING;

-- 2. Household: Família Silva Demo
INSERT INTO households (id, name, currency, locale)
VALUES (
    'b0000000-0000-0000-0000-000000000001',
    'Família Silva Demo',
    'BRL',
    'pt-BR'
) ON CONFLICT (id) DO NOTHING;

-- 3. Household Membership (Owner)
INSERT INTO household_members (household_id, user_id, role)
VALUES (
    'b0000000-0000-0000-0000-000000000001',
    'a0000000-0000-0000-0000-000000000001',
    'owner'
) ON CONFLICT (household_id, user_id) DO NOTHING;

-- 4. Accounts
-- 4.1 Checking Account: Nubank Conta (R$ 8.420,00 = 842000 cents)
INSERT INTO financial_accounts (id, household_id, name, type, currency, initial_balance_minor, current_balance_minor)
VALUES (
    'c0000000-0000-0000-0000-000000000001',
    'b0000000-0000-0000-0000-000000000001',
    'Nubank Conta Principal',
    'checking',
    'BRL',
    842000,
    842000
) ON CONFLICT (id) DO NOTHING;

-- 4.2 Savings: Reserva de Emergência (R$ 15.800,00 = 1580000 cents)
INSERT INTO financial_accounts (id, household_id, name, type, currency, initial_balance_minor, current_balance_minor)
VALUES (
    'c0000000-0000-0000-0000-000000000002',
    'b0000000-0000-0000-0000-000000000001',
    'Reserva de Emergência',
    'savings',
    'BRL',
    1580000,
    1580000
) ON CONFLICT (id) DO NOTHING;

-- 4.3 Credit Card: Nubank Platinum (Current bill: R$ 3.120,00 = 312000 cents liability)
INSERT INTO financial_accounts (id, household_id, name, type, currency, initial_balance_minor, current_balance_minor, credit_limit_minor, statement_closing_day, statement_due_day)
VALUES (
    'c0000000-0000-0000-0000-000000000003',
    'b0000000-0000-0000-0000-000000000001',
    'Nubank Platinum',
    'credit_card',
    'BRL',
    0,
    -312000,
    1200000,
    25,
    5
) ON CONFLICT (id) DO NOTHING;

-- 5. Categories
INSERT INTO categories (id, household_id, name, icon, color) VALUES
('d0000000-0000-0000-0000-000000000001', 'b0000000-0000-0000-0000-000000000001', 'Alimentação', 'shopping-cart', '#10B981'),
('d0000000-0000-0000-0000-000000000002', 'b0000000-0000-0000-0000-000000000001', 'Moradia & Contas', 'home', '#3B82F6'),
('d0000000-0000-0000-0000-000000000003', 'b0000000-0000-0000-0000-000000000001', 'Transporte', 'car', '#F59E0B'),
('d0000000-0000-0000-0000-000000000004', 'b0000000-0000-0000-0000-000000000001', 'Saúde & Farmácia', 'heart', '#EF4444'),
('d0000000-0000-0000-0000-000000000005', 'b0000000-0000-0000-0000-000000000001', 'Salário & Renda', 'briefcase', '#059669'),
('d0000000-0000-0000-0000-000000000006', 'b0000000-0000-0000-0000-000000000001', 'Assinaturas & Lazer', 'tv', '#8B5CF6')
ON CONFLICT (id) DO NOTHING;
