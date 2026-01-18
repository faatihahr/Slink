-- URL Shortener Database Schema for Supabase

-- URLs table to store shortened links
CREATE TABLE urls (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    original_url TEXT NOT NULL,
    short_code VARCHAR(10) UNIQUE NOT NULL,
    custom_alias VARCHAR(20) UNIQUE,
    hit_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_short_code ON urls(short_code);
CREATE INDEX idx_custom_alias ON urls(custom_alias);
CREATE INDEX idx_created_at ON urls(created_at);

-- Function to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger to automatically update updated_at
CREATE TRIGGER update_urls_updated_at 
    BEFORE UPDATE ON urls 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Function to generate unique short code
CREATE OR REPLACE FUNCTION generate_short_code()
RETURNS TEXT AS $$
DECLARE
    new_code TEXT;
    code_exists BOOLEAN;
BEGIN
    LOOP
        -- Generate 6-character alphanumeric code
        new_code := substring(
            encode(gen_random_bytes(6), 'base64'), 
            1, 
            6
        );
        
        -- Remove problematic characters and ensure URL safety
        new_code := regexp_replace(new_code, '[+/=]', '', 'g');
        new_code := lower(new_code);
        
        -- Check if code already exists
        SELECT EXISTS(SELECT 1 FROM urls WHERE short_code = new_code) INTO code_exists;
        
        IF NOT code_exists THEN
            EXIT;
        END IF;
    END LOOP;
    
    RETURN new_code;
END;
$$ LANGUAGE plpgsql;

-- Row Level Security (RLS) for public access
ALTER TABLE urls ENABLE ROW LEVEL SECURITY;

-- Policy to allow anyone to read URLs (needed for redirects)
CREATE POLICY "Anyone can view URLs" ON urls
    FOR SELECT USING (true);

-- Policy to allow anyone to insert URLs (needed for creating short links)
CREATE POLICY "Anyone can create URLs" ON urls
    FOR INSERT WITH CHECK (true);

-- Policy to allow anyone to update hit_count (needed for tracking clicks)
CREATE POLICY "Anyone can update hit count" ON urls
    FOR UPDATE USING (true);
