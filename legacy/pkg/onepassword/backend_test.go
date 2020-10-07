package onepassword

import (
	"fmt"
	"testing"
)

type MockOnePassword struct {
	value    string
	signInOk bool
}

func (m *MockOnePassword) Authenticate(domain string, email string, secretKey string, masterPassword string) error {
	if m.signInOk {
		return nil
	}
	return fmt.Errorf("mock op sign in failed")
}

func (m *MockOnePassword) GetItem(vault string, item string) (string, error) {
	if m.value != "" {
		return m.value, nil
	} else {
		return "", fmt.Errorf("mock op get item failed")
	}
}

func TestGet(t *testing.T) {
	item := "item"
	value := "value"

	backend := &Backend{}
	backend.OnePassword = &MockOnePassword{value: value}

	actual, err := backend.Get(item)

	if err != nil {
		t.Fail()
		fmt.Printf("expected nil but got '%s'", err)
	}
	if actual != value {
		t.Fail()
		fmt.Printf("expected '%s' got %s'", value, actual)
	}
}

func TestGet_ErrGet(t *testing.T) {
	backend := &Backend{}
	backend.OnePassword = &MockOnePassword{}

	_, err := backend.Get("nonExistentItem")

	switch err.(type) {
	case *ErrGet:
		actual := err.Error()
		expected := "1password backend get 'nonExistentItem' failed: mock op get item failed"
		if actual != expected {
			t.Fail()
			fmt.Printf("expected '%s' got '%s'", expected, actual)
		}
	default:
		t.Fail()
	}
}

func TestInit(t *testing.T) {
	domain := "https://externalsecretoperator.1password.com"
	email := "externalsecretoperator@example.com"
	secretKey := "AA-BB-CC-DD-EE-FF-GG-HH-II-JJ"
	masterPassword := "MasterPassword12346!"
	vault := "production"

	backend := &Backend{}
	backend.OnePassword = &MockOnePassword{signInOk: true}

	params := map[string]string{
		"domain":         domain,
		"email":          email,
		"secretKey":      secretKey,
		"masterPassword": masterPassword,
		"vault":          vault,
	}

	err := backend.Init(params)
	if err != nil {
		t.Fail()
		fmt.Println("expected signin to succeed")
	}
}

func TestInit_ErrInitFailed_SignInFailed(t *testing.T) {
	domain := "https://externalsecretoperator.1password.com"
	email := "externalsecretoperator@example.com"
	secretKey := "AA-BB-CC-DD-EE-FF-GG-HH-II-JJ"
	masterPassword := "MasterPassword12346!"
	vault := "production"

	backend := &Backend{}
	backend.OnePassword = &MockOnePassword{signInOk: false}

	params := map[string]string{
		"domain":         domain,
		"email":          email,
		"secretKey":      secretKey,
		"masterPassword": masterPassword,
		"vault":          vault,
	}

	err := backend.Init(params)
	switch err.(type) {
	case *ErrInitFailed:
		actual := err.Error()
		expected := "1password backend init failed: mock op sign in failed"
		if actual != expected {
			t.Fail()
			fmt.Printf("expected '%s' got '%s'", expected, actual)
		}
	default:
		t.Fail()
		fmt.Println("expected init failed error")
	}
}

func TestInit_ErrInitFailed_ParameterMissing(t *testing.T) {
	domain := "https://externalsecretoperator.1password.com"
	secretKey := "AA-BB-CC-DD-EE-FF-GG-HH-II-JJ"
	masterPassword := "MasterPassword12346!"

	backend := NewBackend()

	params := map[string]string{
		"domain":         domain,
		"secretKey":      secretKey,
		"masterPassword": masterPassword,
	}

	err := backend.Init(params)
	switch err.(type) {
	case *ErrInitFailed:
		actual := err.Error()
		expected := "1password backend init failed: expected parameter 'email'"
		if actual != expected {
			t.Fail()
			fmt.Printf("expected '%s' got '%s'", expected, actual)
		}
	default:
		t.Fail()
		fmt.Println("expected init failed error")
	}
}

func TestNewBackend(t *testing.T) {
	backend := NewBackend()

	switch backend.(*Backend).OnePassword.(type) {
	case *Op:
		switch backend.(*Backend).OnePassword.(*Op).GetterBuilder.(type) {
		case *OpGetterBuilder:
		default:
			t.Fail()
			fmt.Println("expected OnePassword GetterBuilder to be OpGetterBuilder")
		}
	default:
		t.Fail()
		fmt.Println("expected OnePassword implementation to be Op")
	}

	expectedVault := "Personal"
	if backend.(*Backend).Vault != expectedVault {
		t.Fail()
		fmt.Printf("expected vault to be equal to '%s'", expectedVault)
	}
}
