CREATE TABLE invoices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_id UUID NOT NULL,
    user_id UUID NOT NULL,
    invoice_number VARCHAR(255) UNIQUE NOT NULL,
    status VARCHAR(50) NOT NULL, -- Enum-like field (e.g., Paid, Overdue, Draft)
    total_amount DECIMAL(15, 2) NOT NULL,
    due_date TIMESTAMP NOT NULL,
    issue_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT fk_customer FOREIGN KEY (customer_id) REFERENCES customers (id) ON DELETE CASCADE,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE SET NULL
);
