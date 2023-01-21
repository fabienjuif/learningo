## Communicaton interface

I choose to be low level to understand Go better.
I am sure there is a lot of cool lib out there, I might want to try a 0mq implementation.

Here is the communication interface chosen for this exercice:

1. The client send the number of `line` it will be sending
2. Then the client send the command name
3. Then the command value for each `line`
4. One line can not exceed 1024 bytes
5. Only strings can be exchanged
6. On each step, the server send an acknowledgment packet of 1 byte

Example 1:

```
1
SET_NAME
Fabien
```

Example 2:

```
3
MESSAGE
Amrit
how are you?
1674303374741
```

Not that `1674303374741` if the client JS Timestamp

## Resources

- socket: https://www.developer.com/languages/intro-socket-programming-go/
- env var: https://www.geeksforgeeks.org/golang-environment-variables/
- dotenv: https://github.com/joho/godotenv
- read the stdin: https://stackoverflow.com/questions/20895552/how-can-i-read-from-standard-input-in-the-console

## Questions

1. I wanted to have a utils.go next to the main.go but I don't know how to do this and if I should do this
