package util

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"
	"sync"

	"github.com/yhshin0/go-auth-server/internal/config"
	"golang.org/x/crypto/argon2"
	"golang.org/x/sync/semaphore"
)

const (
	argon2IDHashFormat = "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
)

type argon2IDBlock struct {
	version   uint8  // 버전
	memoryKiB uint32 // 메모리(KiB)
	timeCost  uint32 // 반복 횟수
	threads   uint8  // 사용할 스레드 수
	keyLen    uint32 // 해시값 길이

	sem *semaphore.Weighted // hash 사용을 위한 pool
}

var (
	Argon2ID     *argon2IDBlock
	onceArgon2ID sync.Once
)

func newArgon2IDBlock() {
	onceArgon2ID.Do(func() {
		cfg := config.GetInstance()
		Argon2ID = &argon2IDBlock{
			version:   cfg.Argon2ID.Version,
			memoryKiB: cfg.Argon2ID.MemoryKiB,
			timeCost:  cfg.Argon2ID.TimeCost,
			threads:   cfg.Argon2ID.Threads,
			keyLen:    cfg.Argon2ID.KeyLen,
			sem:       semaphore.NewWeighted(int64(cfg.Argon2ID.MaxConcurrentHashes)),
		}
	})
}

func (a *argon2IDBlock) Hash(ctx context.Context, password string) (string, error) {
	// semaphore 획득
	// 자리가 없으면 대기
	err := a.sem.Acquire(ctx, 1)
	if err != nil {
		return "", err
	}

	// 반드시 반환
	defer a.sem.Release(1)

	// make salt
	salt := make([]byte, 16)
	_, err = rand.Read(salt)
	if err != nil {
		return "", err
	}

	// 실제 Argon2 실행
	hashPw := argon2.IDKey(
		[]byte(password),
		salt,
		a.timeCost,
		a.memoryKiB,
		a.threads,
		a.keyLen,
	)

	// encode base64
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64HashPw := base64.RawStdEncoding.EncodeToString(hashPw)

	phcStr := fmt.Sprintf(
		argon2IDHashFormat,
		a.version,
		a.memoryKiB,
		a.timeCost,
		a.threads,
		b64Salt,
		b64HashPw,
	)

	return phcStr, nil
}

func (a *argon2IDBlock) Verify(password, encodedHash string) (bool, error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return false, fmt.Errorf("invalid encoded hash. encoded: %s", encodedHash)
	}

	decodedSalt, _ := base64.RawStdEncoding.DecodeString(vals[4])
	decodedPw, _ := base64.RawStdEncoding.DecodeString(vals[5])

	hashPw := argon2.IDKey(
		[]byte(password),
		decodedSalt,
		a.timeCost,
		a.memoryKiB,
		a.threads,
		a.keyLen,
	)

	// 타이밍 공격(Timing Attack)을 방지를 위해 ConstantTimeCompare() 사용
	// 일반적인 비교 연산자는 다른 문자를 발견하는 즉시 비교를 멈춤(early return)
	return subtle.ConstantTimeCompare(hashPw, decodedPw) == 1, nil
}
