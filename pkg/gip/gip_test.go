package gip

import (
	"context"
	"errors"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var client *gipClient

func TestMain(m *testing.M) {
	err := godotenv.Load("../../service/main/.env")
	if err != nil {
		log.Fatalln("Error loading .env file in main service ", err)
	}
	provider, err := NewGIPClient()
	client = provider.(*gipClient)
	if err != nil {
		log.Fatalln(err)
	}

	os.Exit(m.Run())
}

func Test_gipClient_CreateUser(t *testing.T) {
	uid, err := client.CreateUser(context.Background(), "test@test.com", "test", "+14155552671", false)
	if err != nil {
		t.Errorf("CreateUser(); err = %s; want err = <nil>", err)
	}

	if uid == "" {
		t.Errorf("CreateUser(); uid = %s; want non-empty string", uid)
	}

	_ = client.DeleteUser(context.Background(), uid)
}

func Test_gipClient_IsUserExists(t *testing.T) {
	email := "test@test.com"
	exists, err := client.IsUserExists(context.Background(), email)
	if err != nil {
		t.Errorf("IsUserExists(); err = %s; want err = <nil>", err)
	}

	if exists {
		t.Errorf("IsUserExists(); exists = %t; want exists = false", exists)
	}

	uid, err := client.CreateUser(context.Background(), email, "test", "+14155552671", false)
	if err != nil {
		t.Fatal(err)
	}

	if uid == "" {
		t.Fatal("uid is empty")
	}

	defer client.DeleteUser(context.Background(), uid)

	exists, err = client.IsUserExists(context.Background(), email)
	if err != nil {
		t.Errorf("IsUserExists(); err = %s; want err = <nil>", err)
	}

	if !exists {
		t.Errorf("IsUserExists(); exists = %t; want exists = true", exists)
	}

}

func Test_gipClient_DeleteUser(t *testing.T) {
	uid, err := client.CreateUser(context.Background(), "test@test.com", "test", "+14155552671", false)
	if err != nil {
		t.Fatal(err)
	}

	err = client.DeleteUser(context.Background(), uid)
	if err != nil {
		t.Errorf("DeleteUser(); err = %s; want err = <nil>", err)
	}
}

func Test_gipClient_DeleteUserByEmail(t *testing.T) {
	email := "test@test.com"
	_, err := client.CreateUser(context.Background(), email, "test", "+14155552671", false)
	if err != nil {
		t.Fatal(err)
	}

	err = client.DeleteUserByEmail(context.Background(), email)
	if err != nil {
		t.Errorf("DeleteUser(); err = %s; want err = <nil>", err)
	}
}

func Test_gipClient_UpdateUser(t *testing.T) {
	uid, err := client.CreateUser(context.Background(), "test@test.com", "test", "+14155552671", false)
	if err != nil {
		t.Fatal(err)
	}
	defer client.DeleteUser(context.Background(), uid)

	updateParams := map[string]interface{}{
		"email":       "lucy@test.com",
		"displayName": "Lucy",
		"phoneNumber": "+14155552672",
		"disableUser": false,
		"customClaims": map[string]interface{}{
			"organization_account": true,
			"organization_user_id": 100000,
			"organization_status":  false,
			"tenant_uuid":          "162765857916276585",
			"organization_uuid":    "162765857916276585",
			"flag":                 "1",
		},
	}
	err = client.UpdateUser(context.Background(), uid, updateParams)
	if err != nil {
		t.Errorf("UpdateUser(); err = %s; want err = <nil>", err)
	}

	authClient, err := client.app.Auth(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	user, err := authClient.GetUser(context.Background(), uid)
	if err != nil {
		t.Fatal(err)
	}

	if user.Email != updateParams["email"] {
		t.Errorf("UpdateUser(); user.Email = %s; want user.Email = %s", user.Email, updateParams["email"])
	}
	if user.DisplayName != updateParams["displayName"] {
		t.Errorf("UpdateUser(); user.DisplayName = %s; want user.DisplayName = %s", user.DisplayName, updateParams["displayName"])
	}
	if user.PhoneNumber != updateParams["phoneNumber"] {
		t.Errorf("UpdateUser(); user.PhoneNumber = %s; want user.PhoneNumber = %s", user.PhoneNumber, updateParams["phoneNumber"])
	}
	if user.Disabled != updateParams["disableUser"] {
		t.Errorf("UpdateUser(); user.Disabled = %t; want user.Disabled = %t", user.Disabled, updateParams["disableUser"])
	}
	if !assert.NotEqual(t, user.CustomClaims, updateParams["customClaims"].(map[string]interface{})) {
		t.Errorf("UpdateUser(); user.CustomClaims = %v; want user.CustomClaims = %v", user.CustomClaims, updateParams["customClaims"])
	}
}

func Test_gipClient_SignIn(t *testing.T) {
	uid, err := client.CreateUser(context.Background(), "test@test.com", "test", "+14155552671", false)
	if err != nil {
		t.Fatal(err)
	}
	defer client.DeleteUser(context.Background(), uid)

	idToken, err := client.SignIn(context.Background(), uid, map[string]interface{}{
		"organization_account": true,
		"organization_user_id": 100000,
		"organization_status":  false,
		"tenant_uuid":          "162765857916276585",
		"organization_uuid":    "162765857916276585",
	})
	if err != nil {
		t.Errorf("SignIn(); err = %v; want nil", err)
	}

	if idToken == "" {
		t.Errorf("SignIn(); idToken = %s; want non-empty string", idToken)
	}
}

func Test_gipClient_VerifyIDToken(t *testing.T) {
	uid, err := client.CreateUser(context.Background(), "test@test.com", "test", "+14155552671", false)
	if err != nil {
		t.Fatal(err)
	}
	defer client.DeleteUser(context.Background(), uid)

	claims := map[string]interface{}{
		"uid":                  "1627658579",
		"organization_account": true,
		"organization_user_id": 100000,
		"organization_status":  false,
		"tenant_uuid":          "162765857916276585",
		"organization_uuid":    "162765857916276585",
	}
	idToken, err := client.SignIn(context.Background(), uid, claims)
	if err != nil {
		t.Fatal(err)
	}

	idTokenClaims, err := client.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		t.Errorf("VerifyIDToken(); err = %v; want nil", err)
	}

	if !assert.NotEqual(t, idTokenClaims, claims) {
		t.Errorf("VerifyIDToken(); token.Claims = %v; want not equal to claims", idTokenClaims)
	}

	err = client.LogOut(context.Background(), uid)
	if err != nil {
		t.Fatal(err)
	}

	idTokenClaims, err = client.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		t.Errorf("VerifyIDToken(); err = %v; want nil", err)
	}

	if !assert.NotEqual(t, idTokenClaims, claims) {
		t.Errorf("VerifyIDToken(); token.Claims = %v; want not equal to claims", idTokenClaims)
	}

}

func Test_gipClient_VerifyIDTokenAndCheckRevoked(t *testing.T) {
	uid, err := client.CreateUser(context.Background(), "test@test.com", "test", "+14155552671", false)
	if err != nil {
		t.Fatal(err)
	}
	defer client.DeleteUser(context.Background(), uid)

	claims := map[string]interface{}{
		"uid":                  "1627658579",
		"organization_account": true,
		"organization_user_id": 100000,
		"organization_status":  false,
		"tenant_uuid":          "162765857916276585",
		"organization_uuid":    "162765857916276585",
	}
	idToken, err := client.SignIn(context.Background(), uid, claims)
	if err != nil {
		t.Fatal(err)
	}

	idTokenClaims, err := client.VerifyIDTokenAndCheckRevoked(context.Background(), idToken)
	if err != nil {
		t.Errorf("VerifyIDTokenAndCheckRevoked(); err = %v; want nil", err)
	}

	if !assert.NotEqual(t, idTokenClaims, claims) {
		t.Errorf("VerifyIDTokenAndCheckRevoked(); token.Claims = %v; want not equal to claims", idTokenClaims)
	}

	err = client.LogOut(context.Background(), uid)
	if err != nil {
		t.Fatal(err)
	}

	idTokenClaims, err = client.VerifyIDTokenAndCheckRevoked(context.Background(), idToken)
	if !errors.Is(err, ErrIDTokenInvalid) {
		t.Errorf("VerifyIDTokenAndCheckRevoked(); err = %v; want ErrIDTokenInvalid", err)
	}
}
