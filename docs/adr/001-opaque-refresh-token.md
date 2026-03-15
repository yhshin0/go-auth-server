# ADR-001: Opaque Refresh Token + SHA-256 Hash 저장

- **Status**: Accepted
- **Date**: 2026-03-08

---

## Context

인증 서버의 refresh token 구현 방식을 결정해야 했다.
주요 선택지는 두 가지였다.

- **JWT refresh token**: 서명된 토큰. 서버가 상태를 저장하지 않아도 검증 가능.
- **Opaque refresh token**: 랜덤 문자열. 서버에 저장된 데이터와 대조해 검증.

---

## Decision

**Opaque token을 사용하고, 서버에는 SHA-256 hash만 저장한다.**

---

## Rationale

### JWT refresh를 선택하지 않은 이유

JWT refresh token은 stateless가 장점이지만, 이 서버는 다음 요구사항을 가진다.

- 로그아웃 시 해당 refresh 토큰 및 세션 즉시 무효화 (revocation)
- 다중 세션 관리 (유저당 최대 5개)
- Rotating refresh token (1회용 폐기)

이 요구사항들을 충족하려면 어차피 서버에 토큰 상태를 저장해야 한다.
저장이 전제되는 순간 JWT의 stateless 장점은 사라진다.
게다가 JWT는 payload가 base64로 인코딩되어 있어 탈취 시 내용이 노출된다.
서버 저장이 필요한 상황에서 payload 노출 리스크까지 감수할 이유가 없다.

Opaque token은 서버 저장을 전제로 설계되어 있고, 토큰 자체에 의미 있는 정보가 없어 탈취되더라도 payload 노출이 없다.

### Raw token을 저장하지 않는 이유

Redis나 DB가 침해될 경우 raw token이 그대로 노출되면 즉시 재사용이 가능하다.
Hash를 저장하면 침해되더라도 원문 복원이 불가능하다.

### SHA-256을 선택한 이유

Refresh token은 충분한 엔트로피를 가진 랜덤 문자열이다.
패스워드와 달리 사전 공격(dictionary attack)이나 rainbow table의 대상이 아니다.
따라서 slow hash(argon2id, bcrypt)가 필요 없다.

SHA-256은 빠르고, Redis 조회 키로 사용하기에 적합하며, 충돌 저항성도 충분하다.

---

## Consequences

- Redis에 `rt:{sha256Hash}` 형태로 저장하며, 원문은 서버 어디에도 저장하지 않는다.
- 토큰 검증은 요청으로 받은 원문을 SHA-256으로 해싱한 뒤 Redis에서 조회하는 방식으로 처리한다.
- Rotating 방식이므로 토큰 사용 시(refresh) 즉시 폐기하고 새 토큰을 발급한다. 재사용 시도는 자동으로 401 처리된다.
