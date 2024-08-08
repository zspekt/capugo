package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"testing"
)

type testCasesGetSecKey struct {
	Description        string
	Want               key
	WantError          error
	EnvVarKeyParameter string // the parameter we'll pass to GetSecKey
	EnvVarKeyToSet     string // the key of the environment variable we'll set for the test
	EnvVarVal          string // the value of the environment variable we'll set for the test
}

func TestGetSecKey(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	base64EncodedPemBlockSecKey, err := marshalAndEncode(t, key)
	if err != nil {
		t.Fatal(err)
	}

	base64EncodedPemBlockPubKey, err := marshalAndEncode(t, &key.PublicKey)
	if err != nil {
		t.Fatal(err)
	}

	trashString := "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEFUHAEIUFHAIUE1203984612380768DUF-GDSAKLJFGHADHFGAKLJGH213UI560129836120835KLJGH"
	length := base64.StdEncoding.EncodedLen(len(trashString))
	base64EncodedTrash := make([]byte, length)
	base64.StdEncoding.Encode(base64EncodedTrash, []byte(trashString))

	failCases := []testCasesGetSecKey{
		{
			Description:        "calling GetSecKey with an empty EnvVarKeyParameter",
			Want:               nil,
			WantError:          errors.New("caller passed an empty string as envVar parameter"),
			EnvVarKeyParameter: "",
			EnvVarKeyToSet:     "",
			EnvVarVal:          "",
		},
		{
			Description: "calling GetSecKey with a parameter that does not correspond to any env vars",
			Want:        nil,
			// we hardcode the EnvVarKey string into WantedError
			// because otherwise we would have to set it at runtime
			WantError:          errors.New("testingInvalid var missing from environment..."),
			EnvVarKeyParameter: "testingInvalid",
			EnvVarKeyToSet:     "",
			EnvVarVal:          "",
		},
		{
			Description: "calling GetSecKey with a parameter that corresponds to an empty env var",
			Want:        nil,
			WantError: errors.New(
				"testingEmptyEnvVarKey var set but empty (how did this happen?)...",
			),
			EnvVarKeyParameter: "testingEmptyEnvVarKey",
			EnvVarKeyToSet:     "testingEmptyEnvVarKey",
			EnvVarVal:          "",
		},
		{
			Description:        "calling GetSecKey with a parameter that corresponds to an env var but does not hold a valid RSA sec key",
			Want:               nil,
			WantError:          errors.New("pem.Decode couldn't find any pem data"),
			EnvVarKeyParameter: "testingInvalidContentsEnvVarKey",
			EnvVarKeyToSet:     "testingInvalidContentsEnvVarKey",
			EnvVarVal:          string(base64EncodedTrash),
		},
		{
			Description: "calling GetSecKey with a parameter that corresponds to an env var that holds an RSA pub key",
			Want:        nil,
			WantError: errors.New(
				"key block type <RSA PUBLIC KEY> does not match isPrivate argument <true>",
			),
			EnvVarKeyParameter: "testingPubKey",
			EnvVarKeyToSet:     "testingPubKey",
			EnvVarVal:          string(base64EncodedPemBlockPubKey),
		},
	}

	for _, test := range failCases {
		t.Run(test.Description, func(t *testing.T) {
			slog.Info("running test", "test", test.Description)
			// if we are supposed to set an environment variable... AKA if this field is NOT empty
			if test.EnvVarKeyToSet != "" {
				t.Setenv(test.EnvVarKeyToSet, test.EnvVarVal)
			}
			key, err := GetSecKey(test.EnvVarKeyParameter)

			if err == nil { // only non-nil errors expected here
				t.Fatal("nil error returned while ranging over the fail cases")
			}

			gotVal := key == nil && test.Want == nil
			gotErr := err.Error() == test.WantError.Error()
			// gotErr := errors.Is(err, test.WantError)

			switch {
			case gotVal && gotErr: // WIN
				break
			case gotVal && !gotErr: // GOOD VAL, BAD ERROR
				t.Fatalf("good value. got error -> %v || wanted -> %v", err, test.WantError)
			case !gotVal && gotErr: // BAD VAL, GOOD ERROR
				t.Fatalf("good error. got value -> %v || wanted -> %v", key, test.Want)
			case !gotVal && !gotErr: // BAD VAL, BAD ERROR
				t.Fatalf(
					"got error -> <%v> wanted error -> <%v> ##### got value -> <%v> wanted value <%v>",
					err,
					test.WantError,
					key,
					test.Want,
				)
			}
		})
	}
	passCase := testCasesGetSecKey{
		Description:        "calling GetSecKey with a parameter that corresponds to an env var which holds an RSA sec key",
		Want:               nil,
		WantError:          nil,
		EnvVarKeyParameter: "testingSecKey",
		EnvVarKeyToSet:     "testingSecKey",
		EnvVarVal:          string(base64EncodedPemBlockSecKey),
	}

	t.Run(passCase.Description, func(t *testing.T) {
		t.Setenv(passCase.EnvVarKeyToSet, passCase.EnvVarVal)

		seckey, err := GetSecKey(passCase.EnvVarKeyParameter)
		if err != nil {
			t.Fatalf(err.Error())
		}

		if !reflect.DeepEqual(key, seckey) {
			t.Fatalf("Deep equal is false: keys aren't identical") // oh no :(
		}
	})
}

func marshalAndEncode(t testing.TB, key key) ([]byte, error) {
	t.Helper()
	keyType := reflect.TypeOf(key).String()

	var (
		marshalledKey []byte
		pemType       string
		err           error
	)
	switch keyType {
	case "*rsa.PrivateKey":
		marshalledKey, err = x509.MarshalPKCS8PrivateKey(key)
		if err != nil {
			t.Fatal(err)
		}
		pemType = "RSA PRIVATE KEY"
	case "*rsa.PublicKey":
		marshalledKey, err = x509.MarshalPKIXPublicKey(key)
		if err != nil {
			t.Fatal(err)
		}
		pemType = "RSA PUBLIC KEY"
	default:
		return nil, fmt.Errorf("unknown key type: %v", keyType)
	}

	pemBlockMarshalledKey := pem.EncodeToMemory(&pem.Block{
		Type:  pemType,
		Bytes: []byte(marshalledKey),
	})

	// we determine how long the bytes buffer needs to be and make it
	length := base64.StdEncoding.EncodedLen(len(pemBlockMarshalledKey))
	base64EncodedPemBlockKey := make([]byte, length)

	base64.StdEncoding.Encode(base64EncodedPemBlockKey, pemBlockMarshalledKey)

	return base64EncodedPemBlockKey, nil
}
