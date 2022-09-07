# Kissu-Feedback-and-Forms
Kissu's application for handling, creating and linking forms. Usefull for feedback or moderator applications.<br/>
Mod UI is done through ReactJS while the user UI is done through simple HTML which is augmented by JS.<br/>
Javascript is handled through a package containing functions which are used by the HTML pages inside of SCRIPT tags

## Notes
- Routing will be manually tested.
- The UI will be manually tested.
- Templating will be manually tested.

## Tech
 - Gin Server
 - Twig Templates

## Support
Help support Kissu and my development processes: https://ko-fi.com/kissu

## To-Do List
The following is a public list of tests and tasks to be created for this project.
### Server Side
- After each TDD, do a refactor which will (hopefully) fix up any poorly designed code 
- Tests should be using fuzzing to verify multiple different possibilities of inputs(Fuzzing tests can be written after refactor)
- After each test case a cleanup will be done. This means constant reinitialization.
1. Setup: <ins>golang, golang test suite, gin, MySQL</ins> & React
2. Routing into presumed locations and templating
3. TDD for server initialization
    - start and test for settings directory, create test settings directory, build tables in test-DB, password settings, UI settings and finally verify startup.
    - Future changes to initialization must be done here
4. TDD to handle mod login
    - Initialize and test for login function passing on correct pass
    - Initialize and test for login function failing on wrong pass
    - Future changes to passwords must be done here
5. TDD for mod creating a form.
    - A form which will contain all the valid inputs and all realistic combinations of inputs
    - Handle rejection on empty inputs
    - Verify that the form JSON is correct, the form directory is created and it is inserted into the database
    - Future changes to form creation must be done here
6. TDD for displaying a form to a user.
    - Send out the propper JSON for an initialized form
    - Server sends out HTML that is not rendered by JS
    - Form needs a verifiable hash to prevent form spam
7. TDD to handle a user response.
    - Handle inputs on all varieties of potential form inputs
    - Handle error case
    - handle success case, inserting new data into SQL and creating response file
8. TDD to display all forms.
    - As a mod, get a list of all the forms as JSON
9. TDD to display all responses to a given form.
    - As a mod, get a list of all the responses to a form as JSON, these are just methods of acces no extra data
10. TDD to display a singular response.
    - As a mod, get the JSON for a respons
11. TDD to download all responses to a given form.
    - Initialize, create a form and create some test responses
    - Test to determin that the function will properly zip all the given directory properly
    - Will have to unzip to determine correct contents.
12. TDD to delete a form   
    - Delete data from DB for the form and all sql entries on it.
    - Verify the directory and contents within it are unchanged
13. TDD to delete a response
    - Delete given post data from the DB
    - Remove given data file
14. TDD to edit a form
    - Test that new data being sent into the application results in propper changes to data
### Client Side
Create a React library which will initialize multiple types of displays depending on the page's request
1. Mod Login UI
2. Mod Form displays
     - Display of all forms in a select list
     - Selection of list leads to a display of all responses
     - Allows for deleting a form or getting all responses
     - Allows for edits of forms
     - A form can be viewed on it's own
     - Each response is a dropdown and expanding will render the given one
     - Allows for deleting a response
     - A response can be viewed on it's own
3. Mod form edits
     - Form behaves similar to form creation but with a pre-initialized input
     - Old data will remain on already submitted forms, but viewing missing new data should should a NULL
4. Mod form creation
     - Verify no empty data is being sent
5. User form display and submission
     - Past history with given form
     - Shouldn't be done by JS 
     - rather JS adds to the form

## Closing Remarks
  - In the future these test cases will either be augmented or new ones will be created
  - Test cases are to handle the largest possible use-case and as of such changes to the function of use cases should be straight forward
  - My mindset for TDD comes from experiences working with Laravel's PHP testing on https://github.com/ECHibiki/Community-Banners and various issues I had with what I did previously. I'm no pro with automated testing, but having experienced multiple bugs in Vichan whenever I make changes, I see value in it... though I am no diehard who is ideologically commited to it.
