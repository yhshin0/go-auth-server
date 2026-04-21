-- server version: v0.0.1

BEGIN TRANSACTION ISOLATION LEVEL READ COMMITTED;

-- 수정 시간 업데이트를 위한 트리거 함수
CREATE OR REPLACE FUNCTION update_mts_column()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW._mts = now();
    RETURN NEW;
END;
$$ language 'plpgsql';


-- 사용자 계정 정보 테이블
CREATE TABLE IF NOT EXISTS users
(
-- 공통 컬럼
    _id            bigserial   PRIMARY KEY,
    _uuid          uuid        NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    _cts           timestamptz NOT NULL        DEFAULT now(),
    _mts           timestamptz NOT NULL        DEFAULT now(),
-- 비즈니스 컬럼
    email         varchar(128) NOT NULL UNIQUE,
    password_hash varchar(255) NOT NULL
);

CREATE TRIGGER tr_users_update_mts
    BEFORE UPDATE ON users
    FOR EACH ROW
EXECUTE PROCEDURE update_mts_column();

CREATE INDEX idx_users_email ON users (email);

COMMENT ON TABLE users IS '사용자 계정 정보 테이블';
COMMENT ON COLUMN users._uuid IS '사용자 외부 식별자. JWT sub에 사용';
COMMENT ON COLUMN users.password_hash IS 'argon2id hash';

COMMIT;
