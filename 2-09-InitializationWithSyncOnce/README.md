# sync.Once

How to test:

- Run the file once (`go run main.go`) and observe how often the file is loaded.
- In `main()`, change the call to `BadGetFile` into `GoodGetFile` and test again. You should then see a single file load only.