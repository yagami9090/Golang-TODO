package greeting

import "testing"

// TDD: Test Driven Development
func TestGreetingYourName(t *testing.T) {
	t.Skip()
	// AAA Pattern

	// Arrange
	given := "Bob"
	want := "Hello, Bob."

	// Act
	get := Greet(given)

	// Assert
	if want != get {
		t.Errorf("given a name %q want greeting %q, but got %q", given, want, get)
	}
}

func TestGreetingMyFriend(t *testing.T) {
	t.Skip()
	// AAA Pattern

	// Arrange
	given := ""
	want := "Hello, my friend."

	// Act
	get := Greet(given)

	// Assert
	if want != get {
		t.Errorf("given a name %q want greeting %q, but got %q", given, want, get)
	}
}

func TestGreetingCaptital(t *testing.T) {
	t.Skip()
	// AAA Pattern

	// Arrange
	given := "BOB"
	want := "HELLO, BOB."

	// Act
	get := Greet(given)

	// Assert
	if want != get {
		t.Errorf("given a name %q want greeting %q, but got %q", given, want, get)
	}
}
