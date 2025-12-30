-- Миграция: Создание базовой схемы БД
-- Версия: 001
-- Дата: 2024

-- Расширения PostgreSQL
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Таблица пользователей
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    telegram_user_id BIGINT UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    age INTEGER,
    gender VARCHAR(10) CHECK (gender IN ('male', 'female', 'other')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_users_telegram_user_id ON users(telegram_user_id);
CREATE INDEX idx_users_deleted_at ON users(deleted_at) WHERE deleted_at IS NULL;

-- Таблица записей симптомов
CREATE TABLE symptom_entries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date_time TIMESTAMP NOT NULL,
    description TEXT NOT NULL,
    wellbeing_scale INTEGER NOT NULL CHECK (wellbeing_scale >= 1 AND wellbeing_scale <= 10),
    temperature DECIMAL(4,1),
    blood_pressure_systolic INTEGER,
    blood_pressure_diastolic INTEGER,
    pulse INTEGER,
    photo_url VARCHAR(500),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_symptom_entries_user_id ON symptom_entries(user_id);
CREATE INDEX idx_symptom_entries_date_time ON symptom_entries(date_time);
CREATE INDEX idx_symptom_entries_user_date ON symptom_entries(user_id, date_time);

-- Таблица анализов
CREATE TABLE analyses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL CHECK (type IN ('blood', 'urine', 'ultrasound', 'xray', 'other')),
    name VARCHAR(255) NOT NULL,
    date_taken DATE NOT NULL,
    file_url VARCHAR(500) NOT NULL,
    file_type VARCHAR(10) NOT NULL CHECK (file_type IN ('image', 'pdf')),
    next_reminder_date DATE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_analyses_user_id ON analyses(user_id);
CREATE INDEX idx_analyses_date_taken ON analyses(date_taken);
CREATE INDEX idx_analyses_type ON analyses(type);
CREATE INDEX idx_analyses_next_reminder ON analyses(next_reminder_date) WHERE next_reminder_date IS NOT NULL;

-- Таблица лекарств
CREATE TABLE medications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    dosage VARCHAR(100) NOT NULL,
    schedule_type VARCHAR(20) NOT NULL CHECK (schedule_type IN ('daily', 'weekly', 'as_needed')),
    schedule_details JSONB NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_medications_user_id ON medications(user_id);
CREATE INDEX idx_medications_is_active ON medications(is_active);
CREATE INDEX idx_medications_user_active ON medications(user_id, is_active);

-- Таблица приёмов лекарств
CREATE TABLE medication_intakes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    medication_id UUID NOT NULL REFERENCES medications(id) ON DELETE CASCADE,
    scheduled_time TIMESTAMP NOT NULL,
    taken_at TIMESTAMP,
    is_taken BOOLEAN NOT NULL DEFAULT FALSE,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_medication_intakes_medication_id ON medication_intakes(medication_id);
CREATE INDEX idx_medication_intakes_scheduled_time ON medication_intakes(scheduled_time);
CREATE INDEX idx_medication_intakes_is_taken ON medication_intakes(is_taken);

-- Таблица визитов к врачу
CREATE TABLE doctor_visits (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    visit_date DATE NOT NULL,
    doctor_name VARCHAR(255),
    specialty VARCHAR(100),
    questions TEXT,
    report_generated_at TIMESTAMP,
    report_data JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_doctor_visits_user_id ON doctor_visits(user_id);
CREATE INDEX idx_doctor_visits_visit_date ON doctor_visits(visit_date);
CREATE INDEX idx_doctor_visits_user_date ON doctor_visits(user_id, visit_date);

-- Таблица напоминаний
CREATE TABLE reminders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL CHECK (type IN ('medication', 'analysis', 'symptom_check')),
    related_id UUID,
    scheduled_time TIMESTAMP NOT NULL,
    message TEXT NOT NULL,
    is_sent BOOLEAN NOT NULL DEFAULT FALSE,
    sent_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_reminders_user_id ON reminders(user_id);
CREATE INDEX idx_reminders_scheduled_time ON reminders(scheduled_time);
CREATE INDEX idx_reminders_is_sent ON reminders(is_sent);
CREATE INDEX idx_reminders_type ON reminders(type);
CREATE INDEX idx_reminders_upcoming ON reminders(scheduled_time, is_sent) WHERE is_sent = FALSE;

-- Функция для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Триггеры для автоматического обновления updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_symptom_entries_updated_at BEFORE UPDATE ON symptom_entries
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_analyses_updated_at BEFORE UPDATE ON analyses
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_medications_updated_at BEFORE UPDATE ON medications
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_doctor_visits_updated_at BEFORE UPDATE ON doctor_visits
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

