POST {{addr}}/signup
[FormParams]
username: upvote.poster
password: 123
HTTP 200

POST {{addr}}/posts
[FormParams]
title: Hello World
url: https://example.com
HTTP 200
[Captures]
post-id: body

GET {{addr}}/posts/{{post-id}}/upvotes
HTTP 200
[Asserts]
body == "1"

POST {{addr}}/signup
[FormParams]
username: upvote.upvoter
password: 123
HTTP 200

POST {{addr}}/posts/{{post-id}}/upvotes
HTTP 200
[Asserts]
body == "2"

POST {{addr}}/posts/{{post-id}}/upvotes
HTTP 200
[Asserts]
body == "2"

POST {{addr}}/posts/{{post-id}}/downvotes
HTTP 200
[Asserts]
body == "1"

POST {{addr}}/posts/{{post-id}}/downvotes
HTTP 200
[Asserts]
body == "1"
