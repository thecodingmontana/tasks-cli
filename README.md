# tasks-cli

tasks-cli is a command-line interface tool for managing tasks directly in your terminal. With tasks-cli, you can add, view, and manage tasks stored in a CSV file, and display them in a neat tabular format. The project is built using Go and the Cobra library.

## Features

- **Add Tasks**: Easily add tasks with a title, description, and status.
- **View Tasks**: Display tasks in a clear tabular format.
- **Edit Tasks**: Update task details directly from the CLI.
- **Delete Tasks**: Remove tasks from the CSV file.
- **Persist Data**: All tasks are stored in a CSV file for easy management and persistence.

## Requirements

- Go 1.23.4 or later
- Cobra library installed

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/thecodingmontana/tasks-cli.git
   cd tasks-cli
   ```

2. Install dependencies and build the binary:

   ```bash
   go build -o tasks-cli
   ```

3. Move the binary to your PATH:

   ```bash
   mv tasks-cli /usr/local/bin/
   ```

4. Verify installation:

   ```bash
   tasks-cli --help
   ```

## Usage

### Initialize the CSV File

Before adding tasks, initialize the CSV file:

```bash
tasks-cli init
```

### Add a Task

Add a new task by providing a title, description, and optional status:

```bash
tasks-cli add --title "Buy groceries" --description "Milk, eggs, bread" --status "Pending"
```

### View Tasks

View all tasks in a tabular format:

```bash
tasks-cli list
```

### Edit a Task

Update the details of an existing task:

```bash
tasks-cli edit --id 1 --title "Buy groceries" --status "Completed"
```

### Delete a Task

Remove a task by its ID:

```bash
tasks-cli delete --id 1
```

## Example Output

```plaintext
ID   Title            Description          Status
1    Buy groceries    Milk, eggs, bread    Pending
2    Complete report  Due by Monday        In Progress
```

## Configuration

- **CSV File Path**: By default, tasks-cli creates a `tasks.csv` file in the current directory. You can specify a custom file path using the `--file` flag.

```bash
tasks-cli list --file /path/to/custom.csv
```

## Contributing

1. Fork the repository
2. Create a new branch (`git checkout -b feature/your-feature-name`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin feature/your-feature-name`)
5. Open a Pull Request

## License

This project is licensed under the MIT License. See the LICENSE file for details.

## Acknowledgments

- [Cobra Library](https://github.com/spf13/cobra) for the CLI framework.
