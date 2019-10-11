# zmxp v.02a - 10/1/2019

zmxp is an open source Zimbra experience testing utility.


# Features

 As of 9/26/2019, there is only one feature of this app: It will login and screenshot the web client of a user automatically. The purpose is to be sure that the web client is working and that folders appear properly.

## Screenshot of webclient

This app uses Google Chrome headless to screenshot the inbox of an account. This allows rapid automated testing.
The screenshot is saved to the directory under the email address.

## Config

zmxp.ini will be generated automatically upon first run. You will be prompted to for values and given the chance to save those values to the config file.
The app does not save passwords, so you will be prompted each run for a password.

## Input

You can specify a list of users to test using -f=filename.txt or --file=filename.txt
The file should contain an email address of a user to be tested, one per line.

If no file is provided, you will be prompted to test a single email account. Simply enter the address and it will save it.
