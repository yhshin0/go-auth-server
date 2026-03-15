# ADR-005: Refresh Token 재발급

- **Status**: Accepted
- **Date**: 2026-03-13

---

## Context

Refresh token 재발급 시 두 가지를 결정해야 했다.

1. Refresh token을 1회용(rotating)으로 만들지, 정적으로 유지할지
2. 동시 요청으로 인한 이중 소비(double-spend) 문제를 어떻게 처리할지

---

## Decision

**Refresh token은 1회용으로 운영한다. Rotate 전체 과정을 Lua 스크립트 하나로 묶어 원자적으로 실행한다.**

---

## Rationale

### Rotating을 선택한 이유

정적 refresh token은 탈취되어도 서버가 감지할 방법이 없다. 유효 기간(7일) 동안 탈취자가 조용히 사용해도 정상 유저와 구분이 불가능하다.

Rotating은 토큰을 1회용으로 만들어 다음과 같은 탈취 감지 메커니즘을 제공한다.

- 탈취된 토큰이 먼저 사용되면 → 정상 유저의 다음 refresh가 401로 실패한다
- 정상 유저가 먼저 사용하면 → 탈취자의 토큰은 이미 폐기되어 무효다

어느 쪽이든 탈취 상황이 자동으로 드러난다. 완전한 방어는 아니지만, 정적 토큰 대비 공격 창을 크게 줄인다.

### Lua 스크립트로 원자성을 보장한 이유

Rotate 과정은 다음 Redis 명령들로 구성된다.

1. `SET lock:refresh:{oldHash} NX PX 3000` — 락 획득
2. `GET rt:{oldHash}` — sid 조회
3. `HGETALL sess:{sid}` — 세션 검증
4. `DEL rt:{oldHash}`, `DEL sess:{oldSid}` — 기존 데이터 삭제
5. `SET rt:{newHash}`, `HSET sess:{newSid}`, `ZADD user_sessions` — 새 데이터 저장
6. `DEL lock:refresh:{oldHash}` — 락 해제

이 과정을 개별 명령으로 순차 실행하면, 중간 단계에서 장애가 발생했을 때 구 토큰은 삭제되었지만 새 토큰은 아직 저장되지 않은 상태가 생길 수 있다.
이 경우 유저는 유효한 토큰 없이 재로그인을 강제당한다.

Lua 스크립트는 Redis에서 원자적으로 실행된다. 스크립트 전체가 하나의 단위로 처리되므로 중간 상태가 발생하지 않는다.

### SET NX PX를 선택한 이유 (Redlock 등 대안 대비)

Redlock은 다중 Redis 노드에 걸친 분산 락이다. 이 서버는 단일 Redis 노드를 사용하므로 Redlock의 복잡성이 필요 없다.
`SET NX PX`는 단일 노드에서 원자적 락 획득을 보장하며, Lua 스크립트 안에서 실행되므로 충분하다.

---

## Consequences

- 락 TTL은 3000ms(PX 3000)로 설정한다. 스크립트 실행이 이 시간을 초과하는 경우는 사실상 없으나, 락이 영구적으로 남는 상황을 방지하기 위한 안전장치다.
- 동시 요청으로 락 획득에 실패하면 409를 반환한다. 클라이언트는 200~500ms 랜덤 지터 후 1~2회 재시도를 권장한다.
- Redis 단일 노드 장애 시 인증 전체가 중단된다. 고가용성이 필요한 시점에 Redis Sentinel 또는 Cluster 도입을 검토한다.
