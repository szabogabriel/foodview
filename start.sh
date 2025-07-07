#!/bin/bash

readonly PROPERTY_FILE="foodview.properties"

function cleanup_files() {
    if [[ -f "foodview.db" ]]; then
        rm "foodview.db"
    fi
}

function check_properties() {
    if [[ ! -f "$PROPERTY_FILE" ]]; then
        echo "$PROPERTY_FILE file not found. Please create it before running the application."
        exit 1
    fi
}

function set_os_env_from_file() {
    for i in $(cat $PROPERTY_FILE); do
        key=$(echo $i | awk -F '=' '{print $1}' | tr -d ' ')
        value=$(echo $i | awk -F '=' '{print $2}' | tr -d ' ')
        if [[ -n "$key" && -n "$value" ]]; then
            export "$key"="$value"
        fi
    done
}

function print_env() {
    echo "Current environment variables:"
    env | grep -E '^(DB_HOST|DB_PORT|DB_USER|DB_PASSWORD|DB_NAME|DB_SQLITE_FILE)='
}

function run_application() {
    go run main.go
}

function main() {
    cleanup_files
    check_properties
    set_os_env_from_file
    print_env
    run_application
}

main $@