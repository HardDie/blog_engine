# API's

#### Login/registration user
| Endpoint | Method | Task | Body/Query |
|--|--|--|--|
| /api/v1/auth/registration | POST | Register user | username, password, invite |
| /api/v1/auth/login | POST | Login | username, password |
| /api/v1/auth/generate_invite | GET | Generate invite token for authorized user |

#### Authorized user posts control (required cookie)
| Endpoint | Method | Task | Body/Query |
|--|--|--|--|
| /api/v1/users/:id | PATCH | Update user info | password, displayed_name, email, image |
| /api/v1/posts | POST | Create post | title, short, long |
| /api/v1/posts/:id | PATCH | Edit post | title, short, long |
| /api/v1/users/:id/posts | GET | Get list of all posts selected user | page, limit, published, deleted |

#### Common RO user api's
| Endpoint | Method | Task | Body/Query |
|--|--|--|--|
| /api/v1/posts | GET | Get list of all posts from all users (main page) | page, limit, query |
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
	ID int
	UserID int
	Short string
	Body string
	Tags string
	IsPublished bool
	IsDeleted bool
}
```

Comment:
```go
struct Comment {
	ID int
	UserID int
	PostID int
	Text string
}
```

User:
```go
struct User {
	ID int
	Name string
	DisplayedName string
	Email string
	Password string
	InvitedByUser int
}
```

Invite:
```go
struct InviteUser {
	ID int
	UserID int
	InviteCode string
	IsUsed bool
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
