# Readme
All files under MPL2.0<br/>
Designed for kissu.moe by ECVerniy

## Most Urgent Refactors
 - Tools.go removal of functions into former subpackages and etc.
 - fuzzing
 - Gin tests https://gin-gonic.com/docs/testing/ shows how to do unit tests with Gin
   - Notable is "net/http/httptest"
 - Converting Twig to an alternative  https://medium.com/@kataras/whats-the-fastest-template-engine-in-go-fdf0cb95899b
 - Using structs for some initializers... Further condense files( runGin args for example ) into structs
 - Int64 into more generic appropriate types
 - Pass configs and etc. by reference into function
 - Fully remove panics
 - FailureObject should be a map with the FailPosition as a key
 - Responder's response deletes should go into destroyer. Might be other cases like this
 - Storage of versions
 - Status code chekc
 - Rollback inputs on fail

## Testing
Tests cover all functionality except for server startup, templating and routing.<br/>
This project follows, for the most-part, the TDD design philosophy. Additions to the program should be accompanied by a test suite for the given use case. That is, as long as it is not to do with the cases listed above.<br/>
### Reason for not doing 100% TDD
I simply don't want to figure out how to do it at this time. In the next project, or perhaps upon finishing this one.
### Reason for no Fuzzing
I plan on doing this after the project is finished in preparation for the next project: Kissu-Search-and-Archive
