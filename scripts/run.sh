#!/bin/bash

DATABASE_NAME="BUYMEAGIFT"
DB_FILE_NAME="buymeagift.db"
USERNAME=""

# Function to check if the database exists
function database_exists {
    if [ -f $DB_FILE_NAME ]; then
        return 0
    else
        return 1
    fi
}

# Function to drop db, use it when you want to recreate the db
function drop_database {
    echo  "Dropping database: $DATABASE_NAME";
    rm $DB_FILE_NAME
}

# Ask wheter to drop
read -p "Do you want to drop the database? [y/n] " -n 1 -r
if [[ $REPLY =~ ^[Yy]$ ]]
then
    drop_database
fi

# Create database if it doesn't exist
if ! database_exists; then
    echo "Database does not exist. Creating database: $DATABASE_NAME";
    touch $DB_FILE_NAME
    sqlite3 $DB_FILE_NAME < ./scripts/init.sql
else
    echo "Database already exists: $DATABASE_NAME";
fi

# Ask if the user wants to insert some data
read -p "Do you want to insert some data? [y/n] " -n 1 -r
if [[ $REPLY =~ ^[Yy]$ ]]
then
    echo "Inserting data into the database: $DATABASE_NAME";
    sqlite3 $DB_FILE_NAME < ./scripts/seed.sql
fi

# wait for user input
read -p "Press [Enter] key to continue..."
