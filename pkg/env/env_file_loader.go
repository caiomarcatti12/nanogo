/*
 * Copyright 2023 Caio Matheus Marcatti Calimério
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package env

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// fileLoader checks if the .env file exists in the expected location and loads its variables
func fileLoader() error {
	// Get the directory where the binary is located
	executable, err := os.Executable()
	if err != nil {
		log.Fatal("Error obtaining the binary path:", err)
	}
	execDir := filepath.Dir(executable)

	// List of possible locations where the .env file might be
	possiblePaths := []string{
		filepath.Join(execDir, ".env"), // Same directory as the binary
		"../../configs/.env",           // Expected path for development
		".env",                         // Project root directory
	}

	// Check and load the first available .env file
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			// Load the .env file
			if err := godotenv.Load(path); err != nil {
				log.Fatalf("Error loading .env file at path: %s, error: %v", path, err)
			}
			log.Println("Successfully loaded .env file:", path)
			return nil
		}
	}

	// If no .env file is found, print an error and exit
	log.Fatal("No .env file found in the expected locations")
	return nil
}
