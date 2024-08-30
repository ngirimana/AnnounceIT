# AnnounceIT

## Project Overview

Since the rise of the internet, more people have switched their attention from spending hours
watching TV and listening to the radio or even reading newspapers to scrolling and tapping on
their phones and laptops. Businesses can now reach more eyeballs online than the more
traditional approaches.

AnnounceIT comes in as a solution to broadcasting agencies which will will allow them to be able
to receive and manage announcements.

Required Features

1. User can sign up.
2. User can sign in.
3. User (advertiser) can create an announcement.
4. User (advertiser) can update the details of the announcement.
5. User (advertiser) can view all his/her announcements.
6. User (advertiser) can view all announcements of a specific state - accepted, declined,
active, deactivated.
7. User (advertiser) can view a specific announcement.
8. User (admin) can delete an announcement.
9. User (admin) can change the status of an announcement.
10. User (admin) can view announcements from all users.
11. User (admin) can view all announcements of a specific state - accepted, declined
12. User can reset the password.
13. flag/report an announcement as inappropriate.
14. Add a user on a blacklist of people who can’t create announcements.

## Entity Specification

- Users

```{
“id” : Integer,
“email” : String,
“first_name” : String,
“last_name” : String,
“password” : String,
“phoneNumber” : String,
“address” : String,
“is_admin” : Boolean,
}```


- Announcement

```

{
“id” : Integer,
“owner” : Integer, // user id
“status” : String, // pending,accepted, declined,
active,deactivated - default is pending
“text” : String, // announcement copy
“start_date” : DateTime,
“end_date” : DateTime,
“created_on” : DateTime,
}```

- Flags

```{
“id” : Integer,
“announcement_id” : Integer,
“created_on” : DateTime,
“reason” : String, // [sexist, racist,bad language, etc]
“description” : String,
}

```

### Run tests

```go test -v ./...
```

Run the application

```go run main.go
```
Run Swagger

```swag init
```
Swagger URL

```http://localhost:8080/swagger/index.html
```
