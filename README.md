# Kissu-Feedback-and-Forms
Kissu's application for handling, creating and linking forms. Usefull for feedback or moderator applications.

## Notes
- Routing will be manually tested.
- The UI will be manually tested.
- After each test case a cleanup will be done. This means constant reinitialization.
## To-Do List
1. Setup: <ins>golang, golang test suite, gin, MySQL</ins> & React
2. TDD for server initialization 
  - start and test for settings directory, create test settings directory, build tables in test-DB, password settings, UI settings and finally verify startup.
  - Future changes to initialization must be done here
3. TDD to handle mod login 
  - Initialize and test for login function passing on correct pass
  - Initialize and test for login function failing on wrong pass
  - Future changes to passwords must be done here
4. TDD for mod creating a form.
  - A form which will contain all the valid inputs and all realistic combinations of inputs
  - Verify that the response is correct, the form directory is created and it is inserted into the database
  - Future changes to form creation must be done here
5. TDD for displaying a form to a user.
6. TDD to handle a user response.
7. TDD to display all forms.
8. TDD to display all responses to a given form.
9. TDD to display a singular response.
10. TDD to download all responses to a given form.
  - Initialize, create a form and create some test responses
  - Test to determin that the function will properly zip all the given directory properly
  - Will have to unzip to determine correct contents.
  
  ## Closing Remarks
  - In the future these test cases will either be augmented or new ones will be created
  - Test cases are to handle the largest possible use-case and as of such changes to the function of use cases should be straight forward
  - My mindset for TDD comes from experiences working with Laravel's PHP testing on https://github.com/ECHibiki/Community-Banners and various issues I had with what I did previously. I'm no pro with automated testing, but having experienced multiple bugs in Vichan whenever I make changes, I see value in it... though I am no diehard who is ideologically commited to it.
