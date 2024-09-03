package store

import (
	"strings"
	"testing"
)

func TestIsCorrectUsername(t *testing.T) {
	t.Run("length", func(t *testing.T) {
		min_a := strings.Repeat("a", 3)
		max_a := strings.Repeat("a", 21)
		assert(t, isUsernameValid(min_a), "minimum")
		assert(t, isUsernameValid(max_a), "maximum")
		assert(t, !isUsernameValid(min_a[:len(min_a)-1]), "minimum-1")
		assert(t, !isUsernameValid(max_a+"a"), "maximum+1")
	})

	t.Run("dash", func(t *testing.T) {
		assert(t, !isUsernameValid("-username"), "starts with")
		assert(t, !isUsernameValid("username-"), "ends with")
		assert(t, !isUsernameValid("use--rname"), "dash after")
		assert(t, !isUsernameValid("use-.rname"), "dot after")
		assert(t, !isUsernameValid("use.-rname"), "dot before")
		assert(t, !isUsernameValid("use@-rname"), "after at sign")
		assert(t, !isUsernameValid("use-@rname"), "before at sign")
		assert(t, !isUsernameValid("~-username"), "after tilde")
		assert(t, isUsernameValid("us-ername"), "one")
		assert(t, isUsernameValid("us-ern-ame"), "two")
		assert(t, !isUsernameValid("us-Ern-Ame"), "two with captial")
		assert(t, isUsernameValid("us-ern-am-e"), "three")
		assert(t, isUsernameValid("u.s-e.rn-am-e"), "mixed with dots")
		assert(t, isUsernameValid("u.s-e.r@n-am-e"), "mixed with dots and at sign")
		assert(t, isUsernameValid("~u.s-e.r@n-am-e"), "mixed with dots and at sign and tilde")
		assert(t, !isUsernameValid("~u.S-e.R@n-am-e"), "mixed with dots and at sign and tilde and capital lettters")
	})

	t.Run("dot", func(t *testing.T) {
		assert(t, !isUsernameValid(".username"), "starts with")
		assert(t, !isUsernameValid("username."), "ends with")
		assert(t, !isUsernameValid("use..rname"), "dash after")
		assert(t, !isUsernameValid("use.-rname"), "dash after")
		assert(t, !isUsernameValid("use-.rname"), "dash before")
		assert(t, !isUsernameValid("use@.rname"), "after at sign")
		assert(t, !isUsernameValid("use.@rname"), "before at sign")
		assert(t, !isUsernameValid("~.username"), "after tilde")
		assert(t, isUsernameValid("us.ername"), "one")
		assert(t, isUsernameValid("us.ern-ame"), "two")
		assert(t, !isUsernameValid("us.Ern.Ame"), "two with captial")
		assert(t, isUsernameValid("us.ern.am.e"), "three")
		assert(t, isUsernameValid("u-s.e-rn.am.e"), "mixed with dots")
		assert(t, isUsernameValid("u-s.e-r@n.am.e"), "mixed with dots and at sign")
		assert(t, isUsernameValid("~u-s.e-r@n.am.e"), "mixed with dots and at sign and tilde")
		assert(t, !isUsernameValid("~u-S.e-R@n.am.e"), "mixed with dots and at sign and tilde and capital lettters")
	})

	t.Run("at sign", func(t *testing.T) {
		assert(t, !isUsernameValid("@username"), "starts with")
		assert(t, isUsernameValid("us@ername"), "not at any end")
		assert(t, !isUsernameValid("uS@ername"), "not at any end with capital letter before")
		assert(t, !isUsernameValid("us@Ername"), "not at any end with capital letter after")
		assert(t, !isUsernameValid("username@"), "ends with")
		assert(t, !isUsernameValid("~@username"), "after a tilde")
	})

	t.Run("tilde", func(t *testing.T) {
		assert(t, isUsernameValid("~username"), "starts with")
		assert(t, !isUsernameValid("~Username"), "starts with capital letter after")
		assert(t, !isUsernameValid("us~ername"), "not at any end")
		assert(t, !isUsernameValid("username~"), "ends with")
		assert(t, !isUsernameValid("~username~"), "starts and ends with")
		assert(t, !isUsernameValid("~user~name"), "starts with and not at any end")
		assert(t, !isUsernameValid("us~erna~me"), "ends with, not at any end")
	})

	assert(t, isUsernameValid("~a@b"), "normal")
	assert(t, isUsernameValid("a-a"), "normal")
}

func assert(t *testing.T, cond bool, msg string) {
	t.Helper()
	assertEqual(t, true, cond, msg)
}

func assertEqual[T comparable](t *testing.T, expected, actual T, msg string) {
	t.Helper()
	if expected != actual {
		t.Errorf("%s: expected %v, got %v", msg, expected, actual)
	}
}
