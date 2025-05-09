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
package util

import (
	"os"
	"path/filepath"
	"runtime"
)

func GetCurrentDir() string {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		panic("Não foi possível obter diretório atual")
	}
	return filepath.Dir(filename)
}

func GetAbsolutePath(path string) string {
	// Sobe dois níveis a partir do diretório atual
	currentDir := GetCurrentDir()
	rootDir := filepath.Join(currentDir, "..", "..")
	absRootDir, err := filepath.Abs(rootDir)
	if err != nil {
		panic(err)
	}
	return filepath.Join(absRootDir, path)
}

// GetExecutableDir tenta usar o diretório do executável e, caso seja um diretório temporário,
// usa o diretório atual.
func GetExecutableDir() string {
	exePath, err := os.Executable()
	if err != nil {
		panic("não foi possível obter o caminho do executável: " + err.Error())
	}

	// Se for um diretório temporário, retorna o diretório atual.
	if filepath.HasPrefix(exePath, os.TempDir()) {
		dir, err := os.Getwd()
		if err != nil {
			panic("não foi possível obter o diretório atual: " + err.Error())
		}
		return dir
	}

	return filepath.Dir(exePath)
}

func GetExecutableAbsolutePath(path string) string {
	return filepath.Join(GetExecutableDir(), path)
}
