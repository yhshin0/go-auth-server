-- server version: v0.0.1

BEGIN TRANSACTION ISOLATION LEVEL READ COMMITTED;

-- 로그인/로그아웃 이벤트 기록 테이블
CREATE TABLE IF NOT EXISTS auth_events
(
    _id   bigserial   PRIMARY KEY,
    _uuid uuid        NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    _cts  timestamptz NOT NULL        DEFAULT now(),
---
    event_type varchar(16) NOT NULL,
    user__id   bigint,
    session_id varchar(64),
    ip         inet,
    user_agent text
);

COMMENT ON TABLE auth_events IS '로그인/로그아웃 이벤트 기록 테이블';
COMMENT ON COLUMN auth_events.event_type IS 'auth event 타입(e.g. SIGNUP, SIGNIN_SUCCESS, ...)';
COMMENT ON COLUMN auth_events.user__id IS 'users._id 참조';

COMMIT;
