# Transform

1. Create Event Participant POST: /participants => transform to insertable to DB
```json
{
    "event_id":1,
    "participants":[2,3,4,5,6]
}
```

2. Find All Event Participant GET: /participants => transform to easy to use
```json
{
    "data": [
        {
            "event_id": 1,
            "event_name": "Food Fair",
            "event_date": "2025-09-01T16:00:00+07:00",
            "available": false,
            "cancelled": false,
            "organizer": {
                "user_id": 1,
                "role": "organizer",
                "name": "Thanachot_Arm",
                "email": "test@mail.com"
            },
            "participants": [
                {
                    "user_id": 2,
                    "role": "participant",
                    "name": "John",
                    "email": "john@mail.com"
                },
                {
                    "user_id": 3,
                    "role": "participant",
                    "name": "Jane",
                    "email": "Jane@mail.com"
                },
                {
                    "user_id": 4,
                    "role": "participant",
                    "name": "Jack",
                    "email": "Jack@mail.com"
                },
                {
                    "user_id": 5,
                    "role": "participant",
                    "name": "Jim",
                    "email": "Jim@mail.com"
                },
                {
                    "user_id": 6,
                    "role": "participant",
                    "name": "A",
                    "email": "A@mail.com"
                }
            ]
        },
        {
            "event_id": 2,
            "event_name": "Jod Fair",
            "event_date": "2025-09-01T16:00:00+07:00",
            "available": false,
            "cancelled": false,
            "organizer": {
                "user_id": 1,
                "role": "organizer",
                "name": "Thanachot_Arm",
                "email": "test@mail.com"
            },
            "participants": [
                {
                    "user_id": 6,
                    "role": "participant",
                    "name": "A",
                    "email": "A@mail.com"
                }
            ]
        }
    ],
    "limit": 0,
    "offset": 0,
    "total": 2,
    "message": "Success",
    "success": true
}
```