#!/usr/bin/env bash

function separator {
  echo
  echo "===================================================================================="
  echo "==========================[ CLI HELP CHECK SEPARATOR]==============================="
  echo "===================================================================================="
  echo
}

go build -o chast ./main.go

clear
separator

# Test help command
echo ./chast
./chast
separator

# -- Run
echo ./chast run
./chast run
separator

echo ./chast run refactoring
./chast run refactoring
separator

# -- Test
echo ./chast test
./chast test
separator

echo ./chast test refactoring
./chast test refactoring
separator

rm ./chast
