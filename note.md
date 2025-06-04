# Proposition
- **Event Management System**

## Spec/Requirement
- **3 Person/Team**
  - 2 API services
  - 1 service bus
  - 2 hours

## Prerequisite
- No AI
- No Banpu lib
- No chat with Copilot
- FreeStyle framework
- FreeStyle ORM
- Typescript/Golang (Golang optional)

## Score Threshold
- File and folder structure
- Request/response object design
- API design
- DB transaction controller (commit/rollback)
- Handle error
- Code quality (function separation)
- Logging (optional)
- Unit test (optional)

## 1st API (Sender)
- Create event by organizer
- Enroll event by participant (enroll many participants)
- Service bus connect
- As an organizer, can create/update many events and each event can have many participants
- As a participant, can enroll in many events (optional)
- **API**
  - Create event (only event)
  - Update event (only event)
  - Enroll event (enroll many participants)

## Service Bus
- Send message to 2nd API
- Receive message from 1st API
- Rollback transaction (optional)

## 2nd API (Receiver)
- Receive message from service bus and save to database (create/update)
- Get event by ID
- Get all events
- **API**
  - Create event (only event)
  - Update event (only event)
  - Enroll event (enroll many participants)
  - Find event by ID
  - Find all events
- Distinct event response expected (transform) (Example)
  ```json
  [
      {
          "eventId": 1,
          "name": "Event 1",
          "description": "Description 1",
          "startDate": "2021-01-01",
          "endDate": "2021-01-02",
          "location": "Location 1",
          "status": "active",
          "createdAt": "2021-01-01",
          "participant": [
              {
                  "participantId": 1,
                  "name": "John Doe",
                  "email": "john.doe@example.com"
              },
              {
                  "participantId": 2,
                  "name": "Jane Doe",
                  "email": "jane.doe@example.com"
              }
              // ...
          ]
      }
  ]
  ```

## DB
- DB transaction (begin, commit, rollback)
- Update 1 to many (temp table/delete and insert)
- Give ER Diagram to each team
