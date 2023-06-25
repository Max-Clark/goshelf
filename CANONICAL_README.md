# Assessment Notes

- An effort was made to overlap as much of the CLI and API commands as possible
- An effort was made to make the database layer agnostic so future databases could be easily added
- Many tests were left behind due to time
- The following functions would have been added if time was available:
  - CollectionAddBook
  - CollectionRemoveBook
  - BookEdit
- Documentation is a little lacking due to time
- Migrations are configured to be performed with the database schema with the general idea being the program would know current/next version and run the SQL commands in the appropriate folder. Time didn't allow for the implementation.
- The CLI was written for human interaction, but needs to be updated for automation interaction (e.g., shell scripting)