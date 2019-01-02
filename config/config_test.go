package config_test

import (
	"github.com/stretchr/testify/assert"
	"gocleancode/config"
	"os"
	"strconv"
	"testing"
)

func TestFillFromEnvironmentVariables(t *testing.T) {
	expectedPort := 8888
	err := os.Setenv("APP_PORT", strconv.Itoa(expectedPort))
	if err != nil {
		t.Errorf("Expected no error, but got %s instead", err)
		return
	}
	expectedDbName := "test"
	err = os.Setenv("DB_NAME", expectedDbName)
	if err != nil {
		t.Errorf("Expected no error, but got %s instead", err)
		return
	}
	appConfig := config.New()
	if err != nil {
		t.Errorf("Expected no error, but got %s instead", err)
		return
	}
	assert.Equal(t, expectedPort, appConfig.Port)
	assert.Equal(t, expectedDbName, appConfig.DbName)
	assert.Equal(t, 3306, appConfig.DbPort) // Default
}
