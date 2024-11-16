### Logger
    I have implemented a custom logger that records detailed log information. The logger captures the following key elements:

        Timestamp: The exact time the log event occurred.
        File Path: The path to the file where the log event was triggered.
        Log Message: A clear description of the event or error being logged.
        Log Type: The severity or level of the log (e.g., INFO, DEBUG, ERROR, etc.).
    This custom logger helps in tracking events more effectively by providing precise context, making it easier to debug and   
    monitor application behavior.

### Struct instead of parameters
    I refactored the code by passing a struct as a parameter, where its fields represent the parameters needed by various methods 
    and functions. Previously, parameters were passed one by one in every method. By using a struct to group these parameters, it 
    is now possible to add or modify them by updating the struct without having to change each method, function, or interface 
    individually. This approach has made the code more organized, easier to maintain, and has simplified parameter management.

### Error Response refactor
    Instead of creating an error response and writing the HTTP code in every method, I implemented a method that generates a      
    structured error response along with the appropriate HTTP code. This refactor centralizes error handling, ensuring 
    consistency across the application. It reduces code duplication, as I no longer need to repeat the error response logic in 
    each method. 

### In memory cache
    I am introducing Redis as an in-memory cache, which significantly improves the speed of my APIs. By caching frequently 
    accessed data, Redis reduces the need for repetitive database queries, resulting in faster response times. This enhancement 
    optimizes performance, especially for read-heavy operations, and helps to alleviate the load on the database, leading to a 
    more efficient and scalable application.

### Query optimization
    I efficiently fetch all tasks from the database by optimizing the query. This involves using indexing, selecting only the 
    necessary columns, and applying filtering criteria to reduce the data size. These optimizations ensure that the database 
    queries run efficiently, resulting in quicker responses when retrieving task data.

### Add filters for task and status
    I implemented filters for task and status in both the frontend and backend. This allows users to easily search for specific 
    tasks while optimizing data retrieval on the backend to improve performance and enhance the user experience.

### Add sorting feature for task
    I added a sorting feature for tasks by priority in the UI. First, I included a priority field in the Todo model, then 
    implemented functionality to ensure higher priority tasks are displayed at the top. This enhancement improves task 
    organization and visibility.

### Add ability to modify tasks and change their status
    I added the ability to modify tasks and change their status in the UI. This enhancement allows users to easily update task 
    details and manage their progress, improving overall task management.

### Added test cases for FindAll method
