Simple CLI to create md-files with a template.

To build the project run:
    go build -o md_journal

If you run:
    ./md_journal 
it will create by default a task.md file in the `docs/journal` folder.

You can change the default path by running:
    ./md_journal [-w|--wizard]
to run the wizard and run all optional flags.