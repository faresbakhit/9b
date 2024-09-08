package store

import "testing"

func TestIsCorrectUsername(t *testing.T) {
	tests := map[string][]struct {
		string
		bool
	}{
		"dash": {
			{"-username", false},
			{"username-", false},
			{"use--rname", false},
			{"use-.rname", false},
			{"use.-rname", false},
			{"use@-rname", false},
			{"use-@rname", false},
			{"~-username", false},
			{"us-ername", true},
			{"us-ern-ame", true},
			{"us-Ern-Ame", false},
			{"us-ern-am-e", true},
			{"u.s-e.rn-am-e", true},
			{"u.s-e.r@n-am-e", true},
			{"~u.s-e.r@n-am-e", true},
			{"~u.S-e.R@n-am-e", false},
		},
		"dots": {
			{".username", false},
			{"username.", false},
			{"use..rname", false},
			{"use.-rname", false},
			{"use-.rname", false},
			{"use@.rname", false},
			{"use.@rname", false},
			{"~.username", false},
			{"us.ername", true},
			{"us.ern-ame", true},
			{"us.Ern.Ame", false},
			{"us.ern.am.e", true},
			{"u-s.e-rn.am.e", true},
			{"u-s.e-r@n.am.e", true},
			{"~u-s.e-r@n.am.e", true},
			{"~u-S.e-R@n.am.e", false},
		},
		"at sign": {
			{"@username", false},
			{"us@ername", true},
			{"uS@ername", false},
			{"us@Ername", false},
			{"username@", false},
			{"~@username", false},
		},
		"tilde": {
			{"~username", true},
			{"~Username", false},
			{"us~ername", false},
			{"username~", false},
			{"~username~", false},
			{"~user~name", false},
			{"us~erna~me", false},
		},
	}

	for name, tests := range tests {
		t.Run(name, func(t *testing.T) {
			for _, test := range tests {
				isValid := isUsernameValid(test.string)
				if test.bool != isValid {
					t.Errorf("expected %v, got %v", test.bool, isValid)
				}
			}
		})
	}
}
