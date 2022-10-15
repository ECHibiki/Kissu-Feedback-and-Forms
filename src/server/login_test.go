package main

import (
	prebuilder "github.com/ECHibiki/Kissu-Feedback-and-Forms/testing"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/tools"
	"github.com/ECHibiki/Kissu-Feedback-and-Forms/types"
	"strconv"
	"testing"
	"time"
)

func TestPassCreation(t *testing.T) {
	var initialization_folder string = "../../test"
	var err error

	db, init_fields, _ := prebuilder.DoTestingIntializations(initialization_folder)
	defer prebuilder.CleanupTestingInitializations(initialization_folder)

	var stored_pass types.PasswordsDBFields = hashPassword(init_fields.ApplicationPassword, "bcrypt", "10")
	if stored_pass.HashedPassword == init_fields.ApplicationPassword {
		t.Fatal("HashedPassword is same as init_fields.ApplicationPassword")
	}
	if stored_pass.HashSystem != "bcrypt" {
		t.Fatal("HashSystem is not bcrpyt")
	}
	if stored_pass.HashScrambler != "10" {
		t.Fatal("HashScrambler is not using a tested value")
	}
	err = CheckPasswordValid(init_fields.ApplicationPassword, stored_pass.HashedPassword)
	if err != nil {
		t.Fatal("Assigned password does not register as correct before storage")
	}

	err = storePassword(db, stored_pass)
	if err != nil {
		t.Fatal("Error on password storage", err)
	}
	var a_second_stored_pass types.PasswordsDBFields = hashPassword("second-"+init_fields.ApplicationPassword, "bcrypt", "10")
	err = storePassword(db, a_second_stored_pass)
	if err == nil {
		t.Fatal(err)
	}
	// write to it anyways, immitating the effect of a manual DB insertion. It shouldn't effect the outcome
	tools.WritePassToDB(db, a_second_stored_pass)
	potentially_the_second_stored_pass, err := getStoredPassword(db)
	if err != nil {
		t.Fatal(err)
	}
	if potentially_the_second_stored_pass.HashedPassword == a_second_stored_pass.HashedPassword {
		t.Fatal("Passwords being read incorrectly from DB, password table with two collumns should always be reading the top-most, first inserted value")
	}

	retrieved_pass, err := getStoredPassword(db)
	if err != nil {
		t.Fatal(err)
	}
	if stored_pass.HashedPassword != retrieved_pass.HashedPassword {
		t.Fatal("HashedPassword was not stored correctly", retrieved_pass, stored_pass)
	}
	if stored_pass.HashSystem != retrieved_pass.HashSystem {
		t.Fatal("HashSystem was not stored correctly", retrieved_pass, stored_pass)
	}
	if stored_pass.HashScrambler != retrieved_pass.HashScrambler {
		t.Fatal("HashScrambler was not stored correctly", retrieved_pass, stored_pass)
	}
	err = CheckPasswordValid(init_fields.ApplicationPassword, retrieved_pass.HashedPassword)
	if err != nil {
		t.Fatal("Assigned password does not register as correct after storage", init_fields.ApplicationPassword, retrieved_pass)
	}

	err = CheckPasswordValid("Not-"+init_fields.ApplicationPassword, retrieved_pass.HashedPassword)
	if err == nil {
		t.Fatal("The incorrect password registers as valid", init_fields.ApplicationPassword, retrieved_pass)
	}
}

func TestCookieCreation(t *testing.T) {
	// init /////////////

	var initialization_folder string = "../../test"
	var err error

	db, init_fields, _ := prebuilder.DoTestingIntializations(initialization_folder)
	defer prebuilder.CleanupTestingInitializations(initialization_folder)

	var stored_pass types.PasswordsDBFields = hashPassword(init_fields.ApplicationPassword, "bcrypt", "10")

	err = storePassword(db, stored_pass)
	if err != nil {
		t.Fatal("Error on password storage", err)
	}

	// login ///////////

	input_password := init_fields.ApplicationPassword
	retrieved_pass, err := getStoredPassword(db)
	if err != nil {
		t.Fatal(err)
	}
	err = CheckPasswordValid(input_password, retrieved_pass.HashedPassword)
	if err != nil {
		t.Fatal("Assigned password does not register as correct after storage", input_password, retrieved_pass)
	}
	// create cookie
	// SHA256 of Name+Password+Time + random-characters
	session_key_unencrypted := "ADMIN" + input_password + strconv.Itoa(int(time.Now().Unix()))
	session_key_safe := CreateAuthenticationHash(session_key_unencrypted)
	// Store cookie
	var login_fields types.LoginDBFields
	login_fields = CreateLoginFields(session_key_safe, "192.168.1.1")
	err = StoreLogin(db, login_fields)
	if err != nil {
		t.Fatal("Login not stored", err)
	}

	// Verify retrieval for given IP with cookie
	// Test success
	err = CheckCookieValid(db, session_key_safe, "192.168.1.1")
	if err != nil {
		t.Error("true cookie invalid", err, session_key_safe, login_fields)
	}
	// Test failure (Cookie not associated with IP)
	err = CheckCookieValid(db, session_key_safe, "192.168.1.2")
	if err == nil {
		t.Error("false IP is valid")
	}

	// Test failure (Cookie incorrect)
	err = CheckCookieValid(db, session_key_safe+"^", "192.168.1.1")
	if err == nil {
		t.Error("false cookie is valid")
	}
}
