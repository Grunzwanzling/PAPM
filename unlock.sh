#!/bin/bash
unset password
prompt="Enter Password:"
while IFS= read -p "$prompt" -r -s -n 1 char 
do
    if [[ $char == $'\0' ]];     then
        break
    fi
    if [[ $char == $'\177' ]];  then
        prompt=$'\b \b'
        password="${password%?}"
    else
        prompt='*'
        password+="$char"
    fi
done
echo " "
echo "Done. Password=$password" 
