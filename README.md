an Interactive CLI program created using two libraries, [Cobra](https://github.com/spf13/cobra) and [PromptUI](https://github.com/manifoldco/promptui). This program allows users to enter certain topics and store them in a database in a table called wikis. This table has columns such as id, topic, description, created_at, and updated_at. When the user enters a topic via CLI, the data stored includes id, topic, created_at, updated_at, while the description column will remain empty.

Next, the program will use the [GoCron](https://github.com/go-co-op/gocron) library to create a worker that runs every one minute. Every time the worker runs, the program will retrieve all data from the wikis table whose description is still empty. Then, the program will create a concurrent http client to access Wikipedia and retrieve the first paragraph of each topic using the [goquery](https://github.com/PuerkitoBio/goquery) library.

The first paragraph will be inserted into the description column for each topic and will update the updated_at column. In this program, we will input 10 topics to be retrieved from Wikipedia.

