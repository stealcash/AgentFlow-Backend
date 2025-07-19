CREATE TABLE IF NOT EXISTS users (
                                     id BIGSERIAL PRIMARY KEY,
                                     email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    company_name VARCHAR(255),
    user_type VARCHAR(15) NOT NULL DEFAULT 'admin' CHECK (user_type IN ('superadmin', 'admin', 'editor')),
    parent_id BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (parent_id) REFERENCES users(id) ON DELETE SET NULL
    );

CREATE TABLE IF NOT EXISTS chatbots (
                                        id BIGSERIAL PRIMARY KEY,
                                        user_id BIGINT NOT NULL,
                                        chatbot_name VARCHAR(255),
    logo_path VARCHAR(255),
    default_message TEXT,
    public_api_key VARCHAR(255) UNIQUE,
    type VARCHAR(20) NOT NULL DEFAULT 'chatbot' CHECK (type IN ('chatbot', 'slider', 'aichat')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );

CREATE TABLE IF NOT EXISTS categories (
                                          id BIGSERIAL PRIMARY KEY,
                                          chatbot_id BIGINT NOT NULL,
                                          parent_id BIGINT,
                                          name VARCHAR(255) NOT NULL,
    FOREIGN KEY (chatbot_id) REFERENCES chatbots(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES categories(id) ON DELETE SET NULL
    );

CREATE TABLE IF NOT EXISTS category_media (
                                              id BIGSERIAL PRIMARY KEY,
                                              category_id BIGINT NOT NULL UNIQUE,
                                              image_path VARCHAR(255),
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
    );

CREATE TABLE IF NOT EXISTS questions (
                                         id BIGSERIAL PRIMARY KEY,
                                         category_id BIGINT NOT NULL,
                                         question_text TEXT NOT NULL,
                                         answer_text TEXT NOT NULL,
                                         FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
    );

CREATE TABLE IF NOT EXISTS general_questions (
                                                 id BIGSERIAL PRIMARY KEY,
                                                 chatbot_id BIGINT NOT NULL,
                                                 question_text TEXT NOT NULL,
                                                 answer_text TEXT,
                                                 FOREIGN KEY (chatbot_id) REFERENCES chatbots(id) ON DELETE CASCADE
    );

CREATE TABLE IF NOT EXISTS plans (
                                     id BIGSERIAL PRIMARY KEY,
                                     name VARCHAR(100) NOT NULL,
    description TEXT,
    features JSONB,
    price DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE IF NOT EXISTS subscriptions (
                                             id BIGSERIAL PRIMARY KEY,
                                             user_id BIGINT NOT NULL,
                                             plan_id BIGINT NOT NULL,
                                             start_date DATE NOT NULL,
                                             end_date DATE NOT NULL,
                                             status VARCHAR(10) DEFAULT 'active' CHECK (status IN ('active', 'expired', 'cancelled')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (plan_id) REFERENCES plans(id) ON DELETE CASCADE
    );

CREATE TABLE IF NOT EXISTS allowed_domains (
                                               id BIGSERIAL PRIMARY KEY,
                                               chatbot_id BIGINT NOT NULL,
                                               domain VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (chatbot_id) REFERENCES chatbots(id) ON DELETE CASCADE
    );

CREATE TABLE IF NOT EXISTS analytics (
                                         id BIGSERIAL PRIMARY KEY,
                                         chatbot_id BIGINT NOT NULL,
                                         domain VARCHAR(255),
    category_id BIGINT,
    question_id BIGINT,
    input_query TEXT,
    response_source VARCHAR(20) CHECK (response_source IN ('exact_match', 'fuzzy_match', 'general_question', 'chatgpt')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (chatbot_id) REFERENCES chatbots(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL,
    FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE SET NULL
    );
