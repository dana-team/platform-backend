# Hack Scripts

This directory includes several scripts that are helpful for setting up the OpenShift environment for end-to-end (e2e) testing.

## Script Explanation

| **Script Name**      | **Explanation**                                                                                                                                                                                                                                              | **Where is it used?**                                                                                      | **How to use it?**          | **Makefile target** |
|----------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------|-----------------------------|---------------------|
| `create-e2e-user.sh` | This script generates a `bcrypt` hash for a given password and outputs it along with a username in the format `username:hashed_password`. It includes a function to create the bcrypt hash using `OpenSSL` and then prints the username and hashed password. | The output is later used as a constant in the `e2e-tests` setup to create an `HTPasswd` user in OpenShift. | `bash create-e2e-user.sh`   | -                   |
