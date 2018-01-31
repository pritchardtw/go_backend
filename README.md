# go_backend

Hi! Thanks for reviewing my code challenge submission.

When you have the go files you can run go build and build the project. Then run the executable.

The go server will launch on localhost:5000. 

When the server is launched you can run 'go test' and test the server with the few tests I have implemented. 

These tests are not exhaustive, but are just to show that I can write some tests.

I believe the server implements all the requested requirements.

I do not use a real database, I just use a slice and simulate it. Everytime the server resets all passwords
will be lost and it will start from 0 passwords.

# Test hash password post. Will return password index.
curl -X POST --data "password=angryMonkey" localhost:5000/hash

# After 5 seconds, verify password hashed. If before 5 seconds returns invalid request.
curl localhost:5000/hash/0

# Shutdown. There is a 6 second shutdown delay in case it is issued right after a hash password.
curl localhost:5000/shutdown

# Statistics. Will return json formatted statistics. The average is a float because 'int' ms was so small it was rounding to 0.
curl localhost:5000/stats

This is my very first exposure to go and I'm personally excited about how far I've gotten.
