# Kissu-Feedback-and-Forms
Kissu's application for handling, creating and linking forms. Usefull for feedback or moderator applications.<br/>
Mod UI is done through ReactJS while the user UI is done through simple HTML which is augmented by JS.<br/>
Javascript is handled through a package containing functions which are used by the HTML pages inside of SCRIPT tags

## Features
Create, drag, drop, respond, view, edit or download(CSV and JSON) for forms with a single user login. Allows for file submissions. text submissions and a variety of input types. Do not need to know any HTML to use this service. Forms allow for anonymous submission, hashing the users IP, or for it to be IP based. User display of forms does not use any javascript, but mod displays allow for a variety of vanilla javascript functionality.

## Tech
 - Gin Server
 - Twig Templates
 - Fuzzing tests and 'Test Driven Development'

## Support
Help support Kissu and my development processes: https://ko-fi.com/kissu



## Closing Remarks
  - In the future these test cases will either be augmented or new ones will be created
  - Test cases are to handle the largest possible use-case and as of such changes to the function of use cases should be straight forward
  - My mindset for TDD comes from experiences working with Laravel's PHP testing on https://github.com/ECHibiki/Community-Banners and various issues I had with what I did previously. I'm no pro with automated testing, but having experienced multiple bugs in Vichan whenever I make changes, I see value in it... though I am no diehard who is ideologically commited to it.
