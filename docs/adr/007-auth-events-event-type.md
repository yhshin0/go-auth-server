# ADR-007: auth_events event_type — DB enum 대신 varchar + 코드 상수

- **Status**: Accepted
- **Date**: 2026-03-15

---

## Context

`auth_events` 테이블의 `event_type` 컬럼 타입을 결정해야 했다.
허용 값은 `SIGNUP`, `SIGNIN_SUCCESS`, `SIGNIN_FAIL`, `REFRESH_SUCCESS`, `SIGNOUT` 네 가지다.

주요 선택지는 두 가지였다.

- **PostgreSQL enum 타입**: DB 레벨에서 허용 값을 강제.
- **varchar + 코드 상수**: 컬럼은 문자열로 두고, 허용 값은 애플리케이션 코드에서 강제.

---

## Decision

**`varchar(16) NOT NULL`로 선언하고, 허용 값은 코드 상수로 강제한다.**

---

## Rationale

### PostgreSQL enum을 선택하지 않은 이유

PostgreSQL enum은 새로운 값을 추가할 때 `ALTER TYPE ... ADD VALUE`가 필요하다. 이 DDL은 트랜잭션 안에서 실행할 수 없어 마이그레이션 롤백이 어렵다.
또한 값 추가 시 DB 마이그레이션이 애플리케이션 배포보다 먼저 실행되어야 하므로 배포 순서를 엄격하게 관리해야 한다.

`auth_events`는 로그성 테이블로, 향후 이벤트 타입이 추가될 가능성이 있다. 매번 마이그레이션 부담을 감수하는 것보다 코드 레벨에서 관리하는 편이 변경 비용이 낮다.

### varchar + 코드 상수를 선택한 이유

코드 상수(Go의 `const` 또는 `iota` 기반 타입)로 허용 값을 강제하면, 새로운 이벤트 타입 추가 시 코드 변경만으로 충분하다. DB 마이그레이션이 필요 없고 배포 순서 제약도 없다.

타입 안정성은 애플리케이션 레벨에서 충분히 보장할 수 있다. `auth_events`는 외부에서 직접 쓰는 테이블이 아니라 애플리케이션을 통해서만 기록되므로, DB 레벨 강제가 없어도 무결성 리스크가 낮다.

---

## Consequences

- `event_type`의 허용 값은 코드에서 상수로 정의하며, 해당 상수 외의 값은 삽입하지 않는다.
- DB 레벨 강제가 없으므로 직접 SQL로 잘못된 값을 삽입하는 것을 막지 못한다. 운영 환경에서 직접 쓰기는 지양한다.
- 이벤트 타입 추가 시 코드 상수 정의와 삽입 로직만 수정하면 된다. DB 마이그레이션은 불필요하다.
- 필요 시 `CHECK constraint`를 추가해 DB 레벨 검증을 보완할 수 있으나, 현재는 적용하지 않는다.
