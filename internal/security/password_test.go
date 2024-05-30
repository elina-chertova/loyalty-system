package security

import "testing"

func TestPasswordHash(t *testing.T) {
	type args struct {
		firstPassword  string
		secondPassword string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Passwords are equal",
			args: args{
				firstPassword:  "123qwerty",
				secondPassword: "123qwerty",
			},
			want: true,
		},
		{
			name: "Passwords aren't equal",
			args: args{
				firstPassword:  "123qwerty",
				secondPassword: "1234qwerty",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				hashPassword, _ := HashPassword(tt.args.firstPassword)
				if got := CheckPasswordHash(tt.args.secondPassword, hashPassword); got != tt.want {
					t.Errorf("HashPassword() + CheckPasswordHash() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
