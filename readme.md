# YYC Community Hub

A project that links out to every other developer/tech community that's hosted in
Calgary and area. The goal isn't to replace every tech community's website with one
MEGA website, but rather link out to existing resources.

# Goals

## Organizer

1. See other events and locations to lessen event conflict
1. Be able to link out to Meetup
1. (stretch) Be able to link out to Event Brite
1. Be able to give a description of their group, their website, and who to contact easily without manual intervention
1. (stretch) Be able to create and manage events natively in the application

## Attendees

1. Have a calendar that shows local events in Calgary
1. Be able to filter down to topics of interest
1. Be able to go to the event page to RSVP

## Admin/Moderator

1. Be able to remove communities that aren't in Calgary

# MVP

1. Scrape meetup events and put them on a calendar
1. Users are able to see and go out to the meetup event

# MVP+

1. Basic event CRUD
1. Basic event attendance
1. github OAUTH

# Tech Stack

The goal is to keep it simple and accessible. Folks shouldn't struggle with setting stuff up just to contribute

1. golang for the backend
2. HTMX for the front-end
3. TBD for specific libraries
4. sqlite database to start with. May expand to postgres later

# TODO

- [x] Hello web app
- [ ] Query the meetup api
- [ ] Sign in with Github
- [ ] Form to create a new community (don't bother with update or delete yet)
- [ ] Event feed (list format first)
- [ ] Event calendar
- [ ] Link to event

# Ideas

- Have the community be able to link to a discord server so that a bot could post about events automatically
- Try the coworking idea except with each community
    - A user can post a scheduled coworking session
    - That session has a limited amount
    - User can flag (if they choose) anyone who doesn't show up
    - The user can post it to a community if they choose to do so
