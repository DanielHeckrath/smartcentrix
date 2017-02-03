package main

import (
	"os"
	"testing"
)

func TestNewDatabase(t *testing.T) {
	t.Run("missing host", testMissingHost)
	t.Run("missing database name", testMissingDatabase)
	t.Run("missing user", testMissingUsername)
	t.Run("missing password", testMissingPassword)
}

func testMissingHost(t *testing.T) {
	os.Setenv(envDatabaseHost, "")
	os.Setenv(envDatabaseUser, "username")
	os.Setenv(envDatabasePassword, "password")
	os.Setenv(envDatabaseName, "database")

	db, err := newDatabase()

	if err != errHostEmpty {
		t.Error("newDatabase() should return errHostEmpty if host is not set in environment variables")
	}

	if db != nil {
		t.Error("database struct should be empty if we cannot initialize database")
	}
}

func testMissingDatabase(t *testing.T) {
	os.Setenv(envDatabaseHost, "host")
	os.Setenv(envDatabaseUser, "username")
	os.Setenv(envDatabasePassword, "password")
	os.Setenv(envDatabaseName, "")

	db, err := newDatabase()

	if err != errNameEmpty {
		t.Error("newDatabase() should return errNameEmpty if database name is not set in environment variables")
	}

	if db != nil {
		t.Error("database struct should be empty if we cannot initialize database")
	}
}

func testMissingUsername(t *testing.T) {
	os.Setenv(envDatabaseHost, "host")
	os.Setenv(envDatabaseUser, "")
	os.Setenv(envDatabasePassword, "password")
	os.Setenv(envDatabaseName, "database")

	db, err := newDatabase()

	if err != errUserEmpty {
		t.Error("newDatabase() should return errUserEmpty if user name is not is set in environment variables")
	}

	if db != nil {
		t.Error("database struct should be empty if we cannot initialize database")
	}
}

func testMissingPassword(t *testing.T) {
	os.Setenv(envDatabaseHost, "host")
	os.Setenv(envDatabaseUser, "username")
	os.Setenv(envDatabasePassword, "")
	os.Setenv(envDatabaseName, "database")

	db, err := newDatabase()

	if err != errPasswordEmpty {
		t.Error("newDatabase() should return errPasswordEmpty if password is not set in environment variables")
	}

	if db != nil {
		t.Error("database struct should be empty if we cannot initialize database")
	}
}
