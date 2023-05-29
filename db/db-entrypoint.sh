#!/bin/bash

sqlite3 db.sqlite3 ".read init.sql"
if [ $? -ne 0 ]
  then
    echo "Error while creating database!"; exit 1
  else
    echo "Success creating database!"; sqlite3 db.sqlite3
fi