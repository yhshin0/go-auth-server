# ADR-006: SameSite=Lax + OriginGuard 조합 CSRF 대응

- **Status**: Accepted
- **Date**: 2026-03-15

---

## Context

Refresh token을 HttpOnly 쿠키로 전달하는 구조에서 CSRF 공격에 대한 대응이 필요했다.
CSRF는 다른 오리진의 페이지가 쿠키를 자동 포함한 요청을 유도하는 공격이다.

주요 선택지는 세 가지였다.

- **SameSite=Strict**: 크로스 사이트 요청에서 쿠키를 완전히 차단.
- **SameSite=Lax + OriginGuard**: 일반 탐색은 허용하되, 민감한 엔드포인트에서 Origin 검증 추가.
- **CSRF token**: 서버가 발급한 토큰을 요청마다 헤더에 포함.

---

## Decision

**SameSite=Lax + OriginGuard 조합을 사용한다.**
OriginGuard는 `/auth/refresh`, `/auth/signout` 에만 적용한다.

---

## Rationale

### SameSite=Strict를 선택하지 않은 이유

Strict는 외부 사이트에서 링크를 클릭해 이동하는 경우에도 쿠키를 전송하지 않는다.
예를 들어 이메일이나 다른 서비스에서 링크를 타고 들어올 때 로그인 상태가 유지되지 않아 UX가 불필요하게 제한된다.

Refresh token 쿠키는 auth 서버에서만 사용되며, 설정되어 있어 일반 페이지 탐색에서는 전송되지 않는다. Strict의 강한 제한이 필요한 상황이 아니다.

### SameSite=Lax만으로 부족한 이유

Lax는 GET 요청의 크로스 사이트 전송은 허용한다. `/auth/refresh`와 `/auth/signout`은 POST이므로 Lax만으로도 자동 전송은 막을 수 있다.
그러나 Lax는 브라우저 구현에 따라 동작이 다를 수 있고, 동일 사이트 내 공격 시나리오를 완전히 차단하지 못한다.

OriginGuard를 추가해 `Origin` 헤더를 allowlist와 대조함으로써 요청 출처를 서버에서 직접 검증한다. 브라우저 정책에만 의존하지 않는 이중 방어다.

### CSRF token을 선택하지 않은 이유

CSRF token은 서버가 토큰을 발급하고 클라이언트가 요청마다 헤더에 포함해야 한다.
이 서버는 동일 오리진 환경을 전제로 하며, SameSite=Lax + OriginGuard 조합으로 동일한 수준의 보호를 더 단순하게 달성할 수 있다.
CSRF token은 구현과 관리 비용 대비 추가 이점이 없다고 판단했다.

### OriginGuard 적용 범위를 refresh/signout으로 제한한 이유

`/auth/signin`은 쿠키가 아직 없는 상태에서 호출된다. CSRF의 전제 조건인 "쿠키 자동 포함"이 성립하지 않으므로 OriginGuard가 필요 없다.

---

## Consequences

- 프론트엔드와 API가 서로 다른 서브도메인 또는 스킴을 사용하는 경우 SameSite=Lax가 의도대로 동작하지 않을 수 있다. 이 경우 SameSite, CORS 정책, OriginGuard allowlist를 함께 재검토해야 한다.
- OriginGuard allowlist는 환경(local/dev/prod)별로 다르게 관리한다.
- `Vary: Origin` 헤더를 응답에 포함해 CDN/프록시의 캐시 오염을 방지한다.
