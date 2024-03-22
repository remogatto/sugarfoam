# Star Wars characters browser

This Go program is a terminal-based application that browses Star Wars
characters using [this](https://starwars-databank.vercel.app/) Star
Wars Database API. It utilizes the Bubble Tea framework for building
terminal user interfaces (TUIs) and the Lip Gloss library for styling
and layout. The application's main components include:

- **Key Bindings**: Defined for navigating through the application,
  including quitting the application.
- **Model**: Represents the application state, including the current
  state of the application (e.g., checking connection, downloading
  data, browsing characters), the characters fetched from the API, and
  various UI components like tables, viewports, and a status bar.
- **View**: Renders the application UI based on the current state of
  the model. It uses Lip Gloss for styling and layout, creating a
  visually appealing interface within the terminal.
- **Update**: Handles incoming messages (events) and updates the
  application state accordingly. This includes handling window size
  changes, key presses, and responses from the API.
- **Init**: Initializes the application, setting up the initial state
  and commands to be executed.

The application starts by checking the connection to the Star Wars
Database API. If the connection is successful, it proceeds to download
character data. The downloaded characters are displayed in a table,
and selecting a character from the table shows detailed information
about the character in a viewport. The application also includes a
spinner to indicate when data is being downloaded.

The program uses the Bubble Tea framework's model-view-update (MVU)
architecture, where the model represents the application state, the
view renders the UI based on the model, and the update function
handles events and updates the model. This architecture allows for a
clean separation of concerns and makes the application easier to
understand and maintain.

The application leverages the Lip Gloss library for styling and
layout, which provides a declarative approach to terminal rendering,
similar to CSS for web development. This allows for the creation of
visually appealing and well-structured terminal interfaces with ease.

To execute the program, please use the following command:

```bash
go run .
```

Ensure you have an active internet connection.
