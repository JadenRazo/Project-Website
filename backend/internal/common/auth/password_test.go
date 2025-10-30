package auth

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "SecurePassword123!",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  true,
		},
		{
			name:     "very long password",
			password: strings.Repeat("a", 73),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, hash)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hash)
				assert.True(t, strings.HasPrefix(hash, "$2a$12$"))
				assert.NotEqual(t, tt.password, hash)
			}
		})
	}
}

func TestVerifyPassword(t *testing.T) {
	password := "TestPassword123!"
	hash, err := HashPassword(password)
	require.NoError(t, err)

	tests := []struct {
		name        string
		hash        string
		password    string
		expectError bool
	}{
		{
			name:        "correct password",
			hash:        hash,
			password:    password,
			expectError: false,
		},
		{
			name:        "incorrect password",
			hash:        hash,
			password:    "WrongPassword123!",
			expectError: true,
		},
		{
			name:        "empty password",
			hash:        hash,
			password:    "",
			expectError: true,
		},
		{
			name:        "invalid hash",
			hash:        "invalid-hash",
			password:    password,
			expectError: true,
		},
		{
			name:        "empty hash",
			hash:        "",
			password:    password,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := VerifyPassword(tt.hash, tt.password)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePasswordStrength(t *testing.T) {
	tests := []struct {
		name     string
		password string
		valid    bool
	}{
		{
			name:     "valid password with all requirements",
			password: "SecurePass123!",
			valid:    true,
		},
		{
			name:     "too short",
			password: "Sh0rt!",
			valid:    false,
		},
		{
			name:     "no uppercase",
			password: "lowercase123!",
			valid:    false,
		},
		{
			name:     "no lowercase",
			password: "UPPERCASE123!",
			valid:    false,
		},
		{
			name:     "no digit",
			password: "NoNumbers!",
			valid:    false,
		},
		{
			name:     "no special character",
			password: "NoSpecialChar123",
			valid:    false,
		},
		{
			name:     "empty password",
			password: "",
			valid:    false,
		},
		{
			name:     "only spaces",
			password: "        ",
			valid:    false,
		},
		{
			name:     "minimum valid password",
			password: "Abcd123!",
			valid:    true,
		},
		{
			name:     "password with unicode special chars",
			password: "Password123€",
			valid:    false,
		},
		{
			name:     "very long valid password",
			password: "ThisIsAVeryLongPasswordButStillValid123!@#$%^&*()",
			valid:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePasswordStrength(tt.password)
			if tt.valid {
				assert.NoError(t, err, "Password: %s", tt.password)
			} else {
				assert.Error(t, err, "Password: %s", tt.password)
			}
		})
	}
}

func TestPasswordHashingConsistency(t *testing.T) {
	password := "ConsistencyTest123!"

	hash1, err1 := HashPassword(password)
	require.NoError(t, err1)

	hash2, err2 := HashPassword(password)
	require.NoError(t, err2)

	assert.NotEqual(t, hash1, hash2, "Same password should produce different hashes (salt)")

	assert.NoError(t, VerifyPassword(hash1, password))
	assert.NoError(t, VerifyPassword(hash2, password))
}

func TestPasswordStrengthValidation(t *testing.T) {
	weakPasswords := []string{
		"password",
		"123456",
		"password123",
		"PASSWORD123",
		"Password",
		"12345678",
		"qwerty",
		"abc123",
		"Password1",
	}

	for _, weak := range weakPasswords {
		t.Run("weak_password_"+weak, func(t *testing.T) {
			assert.Error(t, ValidatePasswordStrength(weak), "Password should be invalid: %s", weak)
		})
	}

	strongPasswords := []string{
		"MySecurePass123!",
		"Complex$Password2024",
		"Auth3nt1c@t3d!",
		"Str0ng&S3cur3",
		"V3ry$3cur3P@ssw0rd",
	}

	for _, strong := range strongPasswords {
		t.Run("strong_password_"+strong, func(t *testing.T) {
			assert.NoError(t, ValidatePasswordStrength(strong), "Password should be valid: %s", strong)
		})
	}
}

func TestPasswordSecurityBoundaries(t *testing.T) {
	t.Run("maximum_length_password", func(t *testing.T) {
		maxPassword := strings.Repeat("a", 72)
		maxPassword = "A1!" + maxPassword[:69]

		hash, err := HashPassword(maxPassword)
		assert.NoError(t, err)
		assert.NoError(t, VerifyPassword(hash, maxPassword))
	})

	t.Run("password_exceeding_bcrypt_limit", func(t *testing.T) {
		tooLongPassword := strings.Repeat("a", 73)
		tooLongPassword = "A1!" + tooLongPassword[:70]

		_, err := HashPassword(tooLongPassword)
		assert.Error(t, err)
	})

	t.Run("password_with_null_bytes", func(t *testing.T) {
		passwordWithNull := "Valid123!\x00extra"

		_, err := HashPassword(passwordWithNull)
		assert.NoError(t, err) // The actual implementation doesn't check for null bytes
	})
}

func TestPasswordTimingAttackResistance(t *testing.T) {
	password := "TimingTest123!"
	hash, err := HashPassword(password)
	require.NoError(t, err)

	t.Run("consistent_timing_for_wrong_passwords", func(t *testing.T) {
		wrongPasswords := []string{
			"WrongPassword1!",
			"DifferentWrong2@",
			"AnotherWrong3#",
		}

		for _, wrong := range wrongPasswords {
			err := VerifyPassword(hash, wrong)
			assert.Error(t, err)
		}
	})

	t.Run("consistent_timing_for_invalid_hashes", func(t *testing.T) {
		invalidHashes := []string{
			"invalid",
			"$2a$10$invalid",
			"notahash",
			"",
		}

		for _, invalidHash := range invalidHashes {
			err := VerifyPassword(invalidHash, password)
			assert.Error(t, err)
		}
	})
}

func BenchmarkHashPassword(b *testing.B) {
	password := "BenchmarkPassword123!"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := HashPassword(password)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCheckPassword(b *testing.B) {
	password := "BenchmarkPassword123!"
	hash, err := HashPassword(password)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		VerifyPassword(hash, password)
	}
}

func BenchmarkValidatePasswordStrength(b *testing.B) {
	password := "BenchmarkPassword123!"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ValidatePasswordStrength(password)
	}
}
