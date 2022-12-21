# API's

#### Login/registration user
| Endpoint | Method | Task | Body/Query | Authorization | Implemented |
|--|--|--|--|--|--|
| /api/v1/auth/register | POST | Register user | username, password, invite | | [x] |
| /api/v1/auth/login | POST | Login | username, password | | [x] |
| /api/v1/auth/user | GET | Get information about current user | | + | [x] |
| /api/v1/auth/logout | POST | Logout active session | | + | [x] |

### Invites
| Endpoint | Method | Task | Body/Query | Authorization | Implemented |
|--|--|--|--|--|--|
| /api/v1/invites/generate | GET | Generate invite token | | + | [x] |
| /api/v1/invites/revoke | DELETE | Revoke generated invite token | | + | [x] |

### Post
| Endpoint | Method | Task | Body/Query | Authorization | Implemented |
|--|--|--|--|--|--|
| /api/v1/posts | POST | Create post | title, short, body, tags, is_publised | + | [x] |
| /api/v1/posts | GET | Get list of posts for authorized user | limit, page, query | + | [x] |
| /api/v1/posts/:id | PUT | Edit post | title, short, body, tags, is_published | + | [x] |
| /api/v1/posts/feed | GET | Get list of all posts from all users (main page) | page, limit, query | | [x] |

#### Authorized user posts control (required cookie)
| Endpoint | Method | Task | Body/Query | Implemented |
|--|--|--|--|--|
| /api/v1/users/:id | PATCH | Update user info | password, displayed_name, email, image |

#### Common RO user api's
| Endpoint | Method | Task | Body/Query | Implemented |
|--|--|--|--|--|
| /api/v1/posts/:id | GET | Get full post |
| /api/v1/users | GET | Get list of all users with short information | query |
| /api/v1/users/:id | GET | Get full information about selected user |
| /api/v1/users/:id/posts | GET | Get list of all posts selected user | page, limit |
| /api/v1/posts/:id/comments | GET | Get list of all comments selected post | page, limit, order |
| /api/v1/posts/:id/comments | POST | Add comment to selected post | text |
| /api/v1/posts/:id/comments/:id | PATCH | Edit selected comment | text |

# Data types

Post:
```go
struct Post {
	ID int32
	UserID int32
	Title string
	Short string
	Body string
	Tags []string
	IsPublished bool
}
```

Comment:
```go
struct Comment {
	ID int32
	UserID int32
	PostID int32
	Text string
}
```

User:
```go
struct User {
	ID int32
	Username string
	DisplayedName string
	Email string
	InvitedByUser int
}
```

Invite:
```go
struct InviteUser {
	ID int32
	UserID int32
	InviteHash string
	IsActivated bool
}
```

Session:
```go
struct Session {
	ID int
	UserID int
	TokenHash string
}
```
