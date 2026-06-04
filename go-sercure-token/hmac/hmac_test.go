package hmac

import (
	"strings"
	"testing"
)

type MockTimerInterface struct {
	TimeValue int64
}

func (mt MockTimerInterface) TimeNow() int64 {
	return mt.TimeValue
}

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name      string
		username  string
		timestamp int64
		signature string
	}{
		{
			name:      "Valid Username",
			username:  "arun",
			timestamp: 1743983631,
			signature: "DDo+qxzsLcLwJTl0Yxw+F8lhGGLAVsy8oxmr+B2bATc=",
		},
		{
			name:      "Empty Username",
			username:  "",
			timestamp: 1743983631,
			signature: "z/Elu322Rjy1y99El9uLS7yFbe4VujzAdGhcnZ8j9j0=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := MockTimerInterface{TimeValue: tt.timestamp}
			hmacToken := TokenHandler{_TimeProvider: mock}
			got := hmacToken.GenerateToken(tt.username)
			parts := strings.Split(got, ".")
			if len(parts) != 3 {
				t.Fatalf("invalid token format: %s", got)
			}
			if parts[2] != tt.signature {
				t.Errorf("GenerateToken() prefix = %q, want %q", parts[2], tt.signature)
			}
		})
	}
}

func TestVerifyToken(t *testing.T) {
	tests := []struct {
		name        string
		token       string
		mockTime    int64
		want        string
		expectError string
	}{
		{
			name:        "Valid Token",
			token:       "YXJ1bg==.1743983631.DDo+qxzsLcLwJTl0Yxw+F8lhGGLAVsy8oxmr+B2bATc=",
			mockTime:    1743983631 + 100, // Within 5 minutes
			want:        "arun",
			expectError: "",
		},
		{
			name:        "Expired Token",
			token:       "YXJ1bg==.1743983631.DDo+qxzsLcLwJTl0Yxw+F8lhGGLAVsy8oxmr+B2bATc=",
			mockTime:    1743983631 + 400, // Beyond 5 minutes
			want:        "",
			expectError: "token expired",
		},
		{
			name:        "Invalid Signature",
			token:       "YXJ1bg==.1743983631.invalidsignature",
			mockTime:    1743983631,
			want:        "",
			expectError: "invalid signature",
		},
		{
			name:        "Invalid Token Format",
			token:       "YXJ1bg==.1743983631",
			mockTime:    1743983631,
			want:        "",
			expectError: "invalid token format",
		},
		{
			name:        "Invalid Timestamp",
			token:       "YXJ1bg==.17439836x1.X/6hsJkKspxN5VjxW9jVgqrCfzGzAyCBD08jwt6ifKY=",
			mockTime:    1743983631,
			want:        "",
			expectError: "invalid timestamp",
		},
		{
			name:        "Invalid Base64 Encoding",
			token:       "YXJ1bg==!.1743983631.lhCwhSiJe3JCDia74rH5Te7sZkhQZETiLDqTlmYwjw4=",
			mockTime:    1743983631,
			expectError: "invalid username encoding",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := MockTimerInterface{TimeValue: tt.mockTime}
			hmacToken := TokenHandler{_TimeProvider: mock}

			got, err := hmacToken.VerifyToken(tt.token)
			if err != nil {
				if tt.expectError == "" || err.Error() != tt.expectError {
					t.Errorf("VerifyToken() error = %q, want %q", err.Error(), tt.expectError)
				}
			} else if got != tt.want {
				t.Errorf("VerifyToken() = %q, want %q", got, tt.want)
			}
		})
	}
}
